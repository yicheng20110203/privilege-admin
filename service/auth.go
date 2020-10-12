package service

import (
    "encoding/json"
    "errors"
    "github.com/gin-gonic/gin"
    "gitlab.ceibsmoment.com/c/mp/cache"
    "gitlab.ceibsmoment.com/c/mp/config"
    "gitlab.ceibsmoment.com/c/mp/global"
    "gitlab.ceibsmoment.com/c/mp/library"
    "gitlab.ceibsmoment.com/c/mp/logger"
    "gitlab.ceibsmoment.com/c/mp/model"
    pbs "gitlab.ceibsmoment.com/c/mp/mppbs"
    "strings"
    "sync"
)

type AuthorizationData struct {
    Token  string `json:"token"`
    UserId int32  `json:"user_id"`
    Os     string `json:"os"`
}

func CreateAuthorizationData(in *AuthorizationData) (resp string, err error) {
    resp, err = library.AesUtil.EncryptString(config.Cfg.Aes.Key, in)
    if err != nil {
        logger.Logger.Errorf("service.CreateAuthorizationData error: %#v", err)
        return "", err
    }
    return
}

func Auth(c *gin.Context, groupUri string, skipSign bool) (bool, error) {
    uri := c.Request.RequestURI
    for _, v := range global.WhiteList {
        if strings.Trim(v, "/") == strings.Trim(uri, "/") {
            return true, nil
        }
    }

    // 鉴权关闭
    if skipSign {
        return true, nil
    }

    // 以privilege为前缀的路由需要鉴权
    if strings.HasPrefix(strings.Trim(uri, "/"), strings.Trim(groupUri, "/")) {
        resp, err := GetLoginInfoByHeader(c)
        if err != nil {
            return false, err
        }

        xToken := resp.Token
        xUserId := resp.UserId
        xOs := resp.Os
        if xToken == "" || xUserId == 0 || xOs == "" {
            logger.Logger.Errorf("service.Auth error: xToken = %s, xUserId = %d, xOs = %s", xToken, xUserId, xOs)
            return false, errors.New("解密后登录参数为空")
        }

        // token有效性验证
        token := strings.Trim(xToken, " ")
        ck := cache.UserCache.CheckToken(xOs, xUserId, token)
        if !ck {
            logger.Logger.Errorf("service.Auth cache.UserCache.CheckToken(%v) false: %v", token)
            return false, errors.New(global.ErrMsg(global.AdminLoginTokenCheckErr))
        }

        users, err := model.PrivilegeAdmin.Find(&pbs.PrivilegeAdminFindReq{
            Page: 1,
            Size: 1,
            Where: &pbs.PrivilegeAdminFindWhere{
                Id: []int32{
                    xUserId,
                },
            },
        })

        // 查询错误
        if err != nil {
            logger.Logger.Errorf("service.Auth <auto login> model.PrivilegeAdmin.Find(%d) error: %v", xUserId, err)
            return false, errors.New(global.ErrMsg(global.AdminSystemErr))
        }

        // 用户不存在
        if users.GetTotalSize() == 0 {
            logger.Logger.Info("service.Auth <auto login> model.PrivilegeAdmin.Find( ", xUserId, ") login id not exist")
            return false, errors.New(global.ErrMsg(global.AdminLoginUserNotFoundErr))
        }
        userInfo := users.GetList()[0]
        roleKey := userInfo.GetRoleKey()
        // 无任何角色，即无任何后台权限
        if roleKey == "" {
            logger.Logger.Info("service.Auth 无任何权限，user_id = ", xUserId)
            return false, errors.New(global.ErrMsg(global.PrivilegeAuthValidErr))
        }

        // 超级管理员
        isAdmin := RoleService.CheckIsAdmin(roleKey)
        if isAdmin {
            logger.Logger.Info("service.Auth 超级管理员无需鉴权")
            return true, nil
        }

        // 只要包含超级管理员角色即超管
        roleKeys := strings.Split(roleKey, ",")

        // 普通管理鉴权
        // 并行获取
        menuKeys := make([]string, 0)
        var wg, wg1 sync.WaitGroup
        wg.Add(len(roleKeys))
        for _, v := range roleKeys {
            go func(v string) {
                // 普通管理员，非超管
                roleMenuMaps, err := model.RoleMenuModel.Find(&pbs.PrivilegeRoleMenuFindReq{
                    Page: 1,
                    Size: global.MaxUrlNum,
                    Where: &pbs.PrivilegeRoleMenuFindWhere{
                        RoleKey:      v,
                        FuzzyRoleKey: true,
                        IsDelete:     []int32{global.IsDeleteNormal},
                    },
                })
                if err != nil {
                    logger.Logger.Errorf("service.Auth model.RoleMenuModel.Find(role_key = %s) error: %v", v, err)
                    return
                }
                for _, l := range roleMenuMaps.GetList() {
                    menuKeys = append(menuKeys, l.GetMenuKey())
                }
                defer wg.Done()
            }(v)
        }
        wg.Wait()

        // 后台路由验证
        wg1.Add(len(menuKeys))
        var backUrls []string
        for _, v := range menuKeys {
            go func(v string) {
                menuBackMap, err := model.MenuBackUrlModel.Find(&pbs.PrivilegeMenuBackUrlFindReq{
                    Page: 1,
                    Size: global.MaxUrlNum,
                    Where: &pbs.PrivilegeMenuBackUrlFindWhere{
                        MenuKey:      v,
                        FuzzyMenuKey: false,
                        IsDelete:     []int32{global.IsDeleteNormal},
                    },
                })
                if err != nil {
                    logger.Logger.Errorf("service.Auth model.MenuBackUrlModel.Find(menu_key = %s) error: %v", v, err)
                    return
                }
                for _, back := range menuBackMap.GetList() {
                    backUrls = append(backUrls, back.GetBackUrl())
                }
                defer wg1.Done()
            }(v)
        }
        wg1.Wait()

        for _, backUrl := range backUrls {
            if strings.Trim(uri, "/") == strings.Trim(backUrl, "/") {
                return true, nil
            }
        }

        logger.Logger.Info("middleware auth: 普通管理员暂无权限【", uri, "】")
        return false, errors.New(global.ErrMsg(global.AdminCheckParamErr))
    }

    return true, nil
}

func GetLoginInfoByHeader(c *gin.Context) (authInfo *AuthorizationData, err error) {
    authorization := c.Request.Header.Get(global.KeyAuthorization)
    if authorization == "" {
        logger.Logger.Info("service.Auth info: authorization = %s", authorization)
        return &AuthorizationData{}, nil
    }

    decryptData, err := library.AesUtil.DecryptString(config.Cfg.Aes.Key, authorization)
    if err != nil {
        logger.Logger.Errorf("service.Auth error: library.AesUtil.DecryptString(%s) error: %#v", authorization, err)
        return nil, err
    }

    dataBytes, _ := json.Marshal(decryptData)
    var resp *AuthorizationData
    err = json.Unmarshal(dataBytes, &resp)
    if err != nil {
        logger.Logger.Errorf("service.Auth service json.Unmarshal(%s) error: %#v", string(dataBytes), err)
        err = errors.New("json反序列化错误")
        return nil, err
    }

    return resp, nil
}
