package controller

import (
    "github.com/gin-gonic/gin"
    api "gitlab.ceibsmoment.com/c/mp"
    "gitlab.ceibsmoment.com/c/mp/cache"
    def "gitlab.ceibsmoment.com/c/mp/global"
    "gitlab.ceibsmoment.com/c/mp/library"
    "gitlab.ceibsmoment.com/c/mp/logger"
    "gitlab.ceibsmoment.com/c/mp/model"
    pbs "gitlab.ceibsmoment.com/c/mp/mppbs"
    "gitlab.ceibsmoment.com/c/mp/service"
    "net/http"
    "strings"
)

type _privilegeUserController struct {
}

var (
    UserController *_privilegeUserController
)

// admin user login & auto login
func (*_privilegeUserController) Login(c *gin.Context) {
    resp, err := service.GetLoginInfoByHeader(c)
    if err != nil {
        logger.Logger.Errorf("_privilegeUserController.Login error: %#v", err)
        c.JSON(http.StatusOK, api.JsonData(def.AdminCheckParamErr, "", nil))
        return
    }

    xToken := resp.Token
    xUserId := resp.UserId
    xOs := resp.Os
    if strings.Trim(xToken, " ") == "" && xUserId != 0 {
        logger.Logger.Errorf("_privilegeUserController.Login auto login params error")
        c.JSON(http.StatusOK, api.JsonData(def.AdminCheckParamErr, "", nil))
        return
    }

    roleKey := ""
    menuItems := make([]*pbs.MenuItem, 0)

    // 用户密码密码登录
    if strings.Trim(xToken, " ") == "" && xUserId == 0 {
        type jsonInput struct {
            LoginName string `json:"login_name"`
            Password  string `json:"password"`
            Os        string `json:"os"`
        }

        var in *jsonInput
        if err := c.ShouldBindJSON(&in); err != nil {
            logger.Logger.Errorf("_privilegeUserController.Login bind json error: %v", err)
            c.JSON(http.StatusOK, api.JsonData(def.AdminCheckParamErr, "", nil))
            return
        }

        if strings.Trim(in.LoginName, " ") == "" || strings.Trim(in.Password, " ") == "" {
            logger.Logger.Errorf("_privilegeUserController.Login params error : username or password empty (username = %s, password = %s)", in.LoginName, in.Password)
            c.JSON(http.StatusOK, api.JsonData(def.AdminCheckParamErr, "", nil))
            return
        }

        users, err := model.PrivilegeAdmin.Find(&pbs.PrivilegeAdminFindReq{
            Page: 1,
            Size: 1,
            Where: &pbs.PrivilegeAdminFindWhere{
                LoginName: in.LoginName,
            },
        })

        // 查询错误
        if err != nil {
            logger.Logger.Errorf("_privilegeUserController.Login model.PrivilegeAdmin.Find(%s) error: %v", in.LoginName, err)
            c.JSON(http.StatusOK, api.JsonData(def.AdminSystemErr, "", nil))
            return
        }

        // 用户不存在
        if users.GetTotalSize() == 0 {
            logger.Logger.Info("_privilegeUserController.Login model.PrivilegeAdmin.Find( ", in.LoginName, ") login name not exist")
            c.JSON(http.StatusOK, api.JsonData(def.AdminLoginUserNotFoundErr, "", nil))
            return
        }

        userInfo := users.GetList()[0]
        salt := userInfo.GetSalt()
        hash := library.MD5(in.Password + salt)
        if userInfo.GetPassword() != hash {
            logger.Logger.Info("_privilegeUserController.Login model.PrivilegeAdmin.Find( ", in.LoginName, ") error: password is wrong")
            c.JSON(http.StatusOK, api.JsonData(def.AdminLoginUserPwdErr, "", nil))
            return
        }

        userCache, err := cache.UserCache.SetUserCache(in.Os, userInfo.GetId())
        if err != nil {
            logger.Logger.Errorf("_privilegeUserController.Login cache.UserCache.SetUserCache(%d) error: %v", userInfo.GetId(), err)
            c.JSON(http.StatusOK, api.JsonData(def.AdminLoginSystemErr, "", nil))
            return
        }
        roleKey = userInfo.GetRoleKey()
        // 计算authorization
        authUserInfo := &service.AuthorizationData{
            Token:  userCache.Token,
            Os:     in.Os,
            UserId: userInfo.GetId(),
        }
        authorization, err := service.CreateAuthorizationData(authUserInfo)
        if err != nil {
            logger.Logger.Errorf("_privilegeUserController.Login service.CreateAuthorizationData(%#v) error: %v", *authUserInfo, err)
            c.JSON(http.StatusOK, api.JsonData(def.AdminLoginSystemErr, "", nil))
            return
        }

        // 无任何角色，即无任何后台权限
        if roleKey == "" {
            logger.Logger.Info("_privilegeUserController.Login 无任何权限，login_name = ", in.LoginName)
            _ = api.WrapOutput(c, def.SuccessCode, "success", &pbs.LoginInfo{
                //Id:    userInfo.GetId(),
                //Os:    in.Os,
                //Token: userCache.Token,
                Authorization: authorization,
                Menus:         menuItems,
                UserInfo: &pbs.UserInfo{
                    LoginName: userInfo.GetLoginName(),
                    Username:  userInfo.GetUsername(),
                },
            })
            return
        }

        // 只要包含超级管理员角色即超管
        isAdmin := service.RoleService.CheckIsAdmin(roleKey)
        // 超级管理员
        if isAdmin {
            // 计算菜单
            menuSearch, err := service.MenuService.GetMenu(c, userInfo.GetId())
            if err != nil {
                logger.Logger.Errorf("_privilegeUserController.Login service.MenuService.GetMenu error: %v", err)
                c.JSON(http.StatusOK, api.ResetJsonData(def.AdminSystemErr, "获取菜单错误", nil))
                return
            }
            menuItems = menuSearch.GetItems()
            _ = api.WrapOutput(c, def.SuccessCode, "success", &pbs.LoginInfo{
                //Id:    userInfo.GetId(),
                //Os:    in.Os,
                //Token: userCache.Token,
                Authorization: authorization,
                Menus:         menuItems,
                UserInfo: &pbs.UserInfo{
                    LoginName: userInfo.GetLoginName(),
                    Username:  userInfo.GetUsername(),
                },
            })
            return
        }

        // 计算菜单
        menuSearch, err := service.MenuService.GetMenu(c, userInfo.GetId())
        if err != nil {
            logger.Logger.Errorf("_privilegeUserController.Login common service.MenuService.GetMenu error: %v", err)
            c.JSON(http.StatusOK, api.ResetJsonData(def.AdminSystemErr, "获取菜单错误", nil))
            return
        }
        menuItems = menuSearch.GetItems()

        _ = api.WrapOutput(c, def.SuccessCode, "success", &pbs.LoginInfo{
            //Id:    userInfo.GetId(),
            //Os:    in.Os,
            //Token: userCache.Token,
            Authorization: authorization,
            Menus:         menuItems,
            UserInfo: &pbs.UserInfo{
                LoginName: userInfo.GetLoginName(),
                Username:  userInfo.GetUsername(),
            },
        })
        return
    }

    // token登录
    token := strings.Trim(xToken, " ")
    userId := xUserId
    ck := cache.UserCache.CheckToken(xOs, int32(userId), token)
    if !ck {
        logger.Logger.Errorf("_privilegeUserController.Login cache.UserCache.CheckToken(%v) false: %v", token)
        c.JSON(http.StatusOK, api.JsonData(def.AdminLoginTokenCheckErr, "", nil))
        return
    }

    users, err := model.PrivilegeAdmin.Find(&pbs.PrivilegeAdminFindReq{
        Page: 1,
        Size: 1,
        Where: &pbs.PrivilegeAdminFindWhere{
            Id: []int32{
                int32(userId),
            },
        },
    })

    // 查询错误
    if err != nil {
        logger.Logger.Errorf("_privilegeUserController.Login <auto login> model.PrivilegeAdmin.Find(%d) error: %v", userId, err)
        c.JSON(http.StatusOK, api.JsonData(def.AdminSystemErr, "", nil))
        return
    }

    // 用户不存在
    if users.GetTotalSize() == 0 {
        logger.Logger.Info("_privilegeUserController.Login <auto login> model.PrivilegeAdmin.Find( ", userId, ") login id not exist")
        c.JSON(http.StatusOK, api.JsonData(def.AdminLoginUserNotFoundErr, "", nil))
        return
    }
    userInfo := users.GetList()[0]
    roleKey = userInfo.GetRoleKey()

    // 计算authorization
    authUserInfo := &service.AuthorizationData{
        Token:  token,
        Os:     xOs,
        UserId: userInfo.GetId(),
    }
    authorization, err := service.CreateAuthorizationData(authUserInfo)
    if err != nil {
        logger.Logger.Errorf("_privilegeUserController.Login auto login service.CreateAuthorizationData(%#v) error: %v", *authUserInfo, err)
        c.JSON(http.StatusOK, api.JsonData(def.AdminLoginSystemErr, "", nil))
        return
    }

    // 无任何角色，即无任何后台权限
    if roleKey == "" {
        logger.Logger.Info("_privilegeUserController.Login 无任何权限，user_id = ", userId)
        _ = api.WrapOutput(c, def.SuccessCode, "success", &pbs.LoginInfo{
            //Id:            userInfo.GetId(),
            //Os:            xOs,
            //Token:         token,
            Authorization: authorization,
            Menus:         menuItems,
            UserInfo: &pbs.UserInfo{
                LoginName: userInfo.GetLoginName(),
                Username:  userInfo.GetUsername(),
            },
        })
        return
    }

    // 只要包含超级管理员角色即超管
    roleKeys := strings.Split(roleKey, ",")
    isAdmin := false
    for _, v := range roleKeys {
        if v == def.AdminRoleKey {
            isAdmin = true
            break
        }
    }
    // 超级管理员
    if isAdmin {
        // 计算菜单
        menuSearch, err := service.MenuService.GetMenu(c, userInfo.GetId())
        if err != nil {
            logger.Logger.Errorf("_privilegeUserController.Login auto login service.MenuService.GetMenu error: %v", err)
            c.JSON(http.StatusOK, api.ResetJsonData(def.AdminSystemErr, "获取菜单错误", nil))
            return
        }
        menuItems = menuSearch.GetItems()
        _ = api.WrapOutput(c, def.SuccessCode, "success", &pbs.LoginInfo{
            //Id:            userInfo.GetId(),
            //Os:            xOs,
            //Token:         token,
            Authorization: authorization,
            Menus:         menuItems,
            UserInfo: &pbs.UserInfo{
                LoginName: userInfo.GetLoginName(),
                Username:  userInfo.GetUsername(),
            },
        })
        return
    }

    // 计算菜单
    menuSearch, err := service.MenuService.GetMenu(c, userInfo.GetId())
    if err != nil {
        logger.Logger.Errorf("_privilegeUserController.Login service.MenuService.GetMenu error: %v", err)
        c.JSON(http.StatusOK, api.ResetJsonData(def.AdminSystemErr, "获取菜单错误", nil))
        return
    }
    menuItems = menuSearch.GetItems()
    _ = api.WrapOutput(c, def.SuccessCode, "success", &pbs.LoginInfo{
        //Id:            userInfo.GetId(),
        //Os:            xOs,
        //Token:         token,
        Authorization: authorization,
        Menus:         menuItems,
        UserInfo: &pbs.UserInfo{
            LoginName: userInfo.GetLoginName(),
            Username:  userInfo.GetUsername(),
        },
    })
    return
}

// AddAdminUser add cms admin user
func (*_privilegeUserController) Add(c *gin.Context) {
    type params struct {
        LoginName string `json:"login_name"`
        Password  string `json:"password"`
        Username  string `json:"username"`
        RoleKey   string `json:"role_key"`
    }
    var in *params
    if err := c.ShouldBindJSON(&in); err != nil {
        logger.Logger.Errorf("_privilegeUserController.Add c.ShouldBindJSON(&in) error: %v", err)
        c.JSON(http.StatusOK, api.JsonData(def.AdminCheckParamErr, "", nil))
        return
    }

    resp, err := model.PrivilegeAdmin.Save(&pbs.PrivilegeAdminSaveReq{
        LoginName: in.LoginName,
        Password:  in.Password,
        Username:  in.Username,
        RoleKey:   in.RoleKey,
    })
    if err != nil {
        logger.Logger.Errorf("_privilegeUserController.Add model.PrivilegeAdmin.Save(%v) error: %v", *in, err)
        c.JSON(http.StatusOK, api.JsonData(def.AdminSystemErr, "", nil))
        return
    }

    _ = api.WrapOutput(c, def.SuccessCode, "success", resp)
    return
}

// UpdateAdminUser update admin user info
func (*_privilegeUserController) Update(c *gin.Context) {
    type params struct {
        Id        int32  `json:"id"`
        LoginName string `json:"login_name"`
        Password  string `json:"password"`
        Username  string `json:"username"`
        RoleKey   string `json:"role_key"`
        IsAdmin   int32  `json:"is_admin"`
    }
    var in *params
    if err := c.ShouldBindJSON(&in); err != nil {
        logger.Logger.Errorf("_privilegeUserController.Update c.ShouldBindJSON(&in) error: %v", err)
        c.JSON(http.StatusOK, api.JsonData(def.AdminCheckParamErr, "", nil))
        return
    }

    userList, err := model.PrivilegeAdmin.Find(&pbs.PrivilegeAdminFindReq{
        Page: 1,
        Size: 1,
        Where: &pbs.PrivilegeAdminFindWhere{
            Id: []int32{in.Id},
        },
    })
    if err != nil {
        logger.Logger.Errorf("_privilegeUserController.Update model.PrivilegeAdmin.Find(%v) error: %v", *in, err)
        c.JSON(http.StatusOK, api.JsonData(def.AdminSystemErr, "", nil))
        return
    }

    if userList == nil || userList.GetTotalSize() == 0 {
        logger.Logger.Errorf("_privilegeUserController.Update model.PrivilegeAdmin.Find(%v) error: %v", *in, err)
        c.JSON(http.StatusOK, api.ResetJsonData(def.AdminSystemErr, "管理员账号不存在", nil))
        return
    }

    if service.RoleService.CheckIsAdmin(userList.GetList()[0].GetRoleKey()) {
        logger.Logger.Info("超级管理员账号不允许修改: admin info = ", userList.GetList()[0])
        c.JSON(http.StatusOK, api.ResetJsonData(def.AdminSystemErr, "超级管理员账号不允许修改", nil))
        return
    }

    resp, err := model.PrivilegeAdmin.Save(&pbs.PrivilegeAdminSaveReq{
        Id:        in.Id,
        LoginName: in.LoginName,
        Password:  in.Password,
        Username:  in.Username,
        RoleKey:   in.RoleKey,
        IsAdmin:   in.IsAdmin,
    })
    if err != nil {
        logger.Logger.Errorf("_privilegeUserController.Update model.PrivilegeAdmin.Save(%v) error: %v", *in, err)
        c.JSON(http.StatusOK, api.JsonData(def.AdminSystemErr, "", nil))
        return
    }

    _ = api.WrapOutput(c, def.SuccessCode, "success", resp)
    return
}

// List user list
func (*_privilegeUserController) List(c *gin.Context) {
    type params struct {
        Page      int32  `json:"page"`
        Size      int32  `json:"size"`
        LoginName string `json:"login_name"`
        Username  string `json:"username"`
        IsAdmin   int32  `json:"is_admin"`
        RoleKey   string `json:"role_key"`
        Id        int32  `json:"id"`
    }
    var in *params
    if err := c.ShouldBindJSON(&in); err != nil {
        logger.Logger.Errorf("_privilegeUserController.List c.ShouldBindJSON(&in) error: %v", err)
        c.JSON(http.StatusOK, api.JsonData(def.AdminCheckParamErr, "", nil))
        return
    }

    where := &pbs.PrivilegeAdminFindWhere{
        LoginName:   in.LoginName,
        Username:    in.Username,
        IsAdmin:     in.IsAdmin,
        RoleKey:     in.RoleKey,
        FilterAdmin: true,
    }
    if in.Id > 0 {
        where.Id = []int32{in.Id}
    }
    resp, err := model.PrivilegeAdmin.Find(&pbs.PrivilegeAdminFindReq{
        Page:  in.Page,
        Size:  in.Size,
        Where: where,
    })
    if err != nil {
        logger.Logger.Errorf("_privilegeUserController.List model.PrivilegeAdmin.Find(%+v) error: %v", *in, err)
        c.JSON(http.StatusOK, api.JsonData(def.AdminSystemErr, "", nil))
        return
    }

    roleKeyMap := make(map[string]struct{})
    lists := make([]*pbs.UserAdminInfo, 0)
    for k, v := range resp.GetList() {
        v.Password = ""
        v.Salt = ""
        resp.List[k] = v
        ks := strings.Split(v.GetRoleKey(), ",")
        for _, kr := range ks {
            roleKeyMap[kr] = struct{}{}
        }

        var admin string
        if v.GetIsAdmin() > 0 {
            admin = "超级管理员"
        } else {
            admin = "普通管理员"
        }
        lists = append(lists, &pbs.UserAdminInfo{
            Id:         v.GetId(),
            LoginName:  v.GetLoginName(),
            Username:   v.GetUsername(),
            RoleKey:    v.GetRoleKey(),
            IsAdmin:    v.GetIsAdmin(),
            Admin:      admin,
            CreateTime: v.GetCreateTime(),
            UpdateTime: v.GetUpdateTime(),
        })
    }
    roleKeys := make([]string, 0)
    for k := range roleKeyMap {
        roleKeys = append(roleKeys, k)
    }
    roleList, err := model.RoleModel.Find(&pbs.PrivilegeRoleFindReq{
        Page: 1,
        Size: int32(len(roleKeys)),
        Where: &pbs.PrivilegeRoleFindWhere{
            RoleKeys: roleKeys,
            Status: []int32{
                def.RoleStart,
            },
        },
    })
    if err != nil {
        logger.Logger.Errorf("_privilegeUserController.List model.RoleModel.Find(%+v) error: %v", *in, err)
        c.JSON(http.StatusOK, api.ResetJsonData(def.AdminSystemErr, "角色批量查询错误", nil))
        return
    }
    roleMap := make(map[string]string)
    for _, v := range roleList.GetList() {
        roleMap[v.GetRoleKey()] = v.GetName()
    }
    getRoles := func(key string) []string {
        if service.RoleService.CheckIsAdmin(key) {
            return []string{
                "超级管理员",
            }
        }

        keys := strings.Split(key, ",")
        rs := make([]string, 0)
        for _, k := range keys {
            if _, ok := roleMap[k]; ok {
                rs = append(rs, roleMap[k])
            }
        }
        return rs
    }
    for k, v := range lists {
        v.Roles = getRoles(v.GetRoleKey())
        lists[k] = v
    }

    userList := &pbs.UserList{
        Page:      resp.GetPage(),
        Size:      resp.GetSize(),
        TotalSize: resp.GetTotalSize(),
        TotalPage: resp.GetTotalPage(),
        List:      lists,
    }

    _ = api.WrapOutput(c, def.SuccessCode, "success", userList)
    return
}

// Logout
// 管理员退出登录
func (*_privilegeUserController) Logout(c *gin.Context) {
    resp, err := service.GetLoginInfoByHeader(c)
    if err != nil {
        logger.Logger.Errorf("_privilegeUserController.Logout service.GetLoginInfoByHeader error: %#v", err)
        c.JSON(http.StatusOK, api.JsonData(def.AdminSystemErr, "", nil))
        return
    }

    xToken := resp.Token
    xUserId := resp.UserId
    xOs := resp.Os
    if strings.Trim(xToken, " ") == "" || xUserId == 0 || strings.Trim(xOs, " ") == "" {
        logger.Logger.Errorf("_privilegeUserController.Logout params error")
        c.JSON(http.StatusOK, api.JsonData(def.AdminCheckParamErr, "", nil))
        return
    }

    token := strings.Trim(xToken, " ")
    userId := xUserId
    clean := cache.UserCache.CleanToken(xOs, userId, token)
    if !clean {
        logger.Logger.Errorf("_privilegeUserController.Logout cache.UserCache.CleanToken(%v) false", token)
        c.JSON(http.StatusOK, api.JsonData(def.AdminLoginTokenCheckErr, "", nil))
        return
    }

    // logout
    c.JSON(http.StatusOK, api.JsonData(def.SuccessCode, "success", nil))
    return
}

// Delete
// Warning 物理删除
func (*_privilegeUserController) Delete(c *gin.Context) {
    type params struct {
        LoginName string `json:"login_name"`
        Username  string `json:"username"`
        IsAdmin   int32  `json:"is_admin"`
        RoleKey   string `json:"role_key"`
        Id        int32  `json:"id"`
    }
    var in *params
    if err := c.ShouldBindJSON(&in); err != nil {
        logger.Logger.Errorf("_privilegeUserController.Delete c.ShouldBindJSON(&in) error: %v", err)
        c.JSON(http.StatusOK, api.JsonData(def.AdminCheckParamErr, "", nil))
        return
    }

    if service.RoleService.CheckIsAdmin(in.RoleKey) || in.IsAdmin == 1 {
        logger.Logger.Errorf("_privilegeUserController.Delete 超级管理员不允许删除 in = %+v", *in)
        c.JSON(http.StatusOK, api.ResetJsonData(def.AdminCheckParamErr, "超级管理员不允许删除,", nil))
        return
    }

    where := &pbs.PrivilegeAdminDeleteReq{
        LoginName: in.LoginName,
        Username:  in.Username,
        IsAdmin:   in.IsAdmin,
        RoleKey:   in.RoleKey,
    }
    if in.Id > 0 {
        where.Id = []int32{in.Id}
        userList, err := model.PrivilegeAdmin.Find(&pbs.PrivilegeAdminFindReq{
            Page: 1,
            Size: 1,
            Where: &pbs.PrivilegeAdminFindWhere{
                Id: []int32{in.Id},
            },
        })
        if err != nil {
            logger.Logger.Errorf("_privilegeUserController.Delete model.PrivilegeAdmin.Find in = %+v, error = %v", *in, err)
            c.JSON(http.StatusOK, api.JsonData(def.AdminCheckParamErr, "", nil))
        }
        if userList != nil && userList.GetTotalSize() > 0 {
            if service.RoleService.CheckIsAdmin(userList.GetList()[0].GetRoleKey()) {
                logger.Logger.Errorf("_privilegeUserController.Delete 超级管理员不允许删除 in = %+v， user = %+v", *in, userList.GetList()[0])
                c.JSON(http.StatusOK, api.ResetJsonData(def.AdminCheckParamErr, "超级管理员不允许删除,", nil))
                return
            }
        }
    }
    resp, err := model.PrivilegeAdmin.Delete(where)
    if err != nil {
        logger.Logger.Errorf("_privilegeUserController.List model.PrivilegeAdmin.Delete(%+v) error: %v", *where, err)
        c.JSON(http.StatusOK, api.JsonData(def.AdminSystemErr, "", nil))
        return
    }

    _ = api.WrapOutput(c, def.SuccessCode, "success", resp)
    return
}
