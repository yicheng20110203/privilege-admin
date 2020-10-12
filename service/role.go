package service

import (
    "github.com/gin-gonic/gin"
    "gitlab.ceibsmoment.com/c/mp/global"
    "gitlab.ceibsmoment.com/c/mp/logger"
    "gitlab.ceibsmoment.com/c/mp/model"
    pbs "gitlab.ceibsmoment.com/c/mp/mppbs"
    "strings"
    "sync"
)

type _roleService struct{}

var (
    RoleService *_roleService
)

// 检查一个角色是否是超级管理员
func (*_roleService) CheckIsAdmin(roleKeys string) bool {
    if roleKeys == "" {
        return false
    }

    keys := strings.Split(roleKeys, ",")
    for _, key := range keys {
        if key == global.AdminRoleKey {
            return true
        }
    }

    return false
}

func (*_roleService) GetStatusDesc(status int32) string {
    switch status {
    case global.RoleUnuse:
        return "未启用"
    case global.RoleStart:
        return "已启用"
    case global.RoleDel:
        return "已删除"
    default:
        return "未知"
    }
}

func (*_roleService) AddMultiPower(c *gin.Context, roleKey string, menuKeys []string) (resp *pbs.RolePowerMultiAddResult, err error) {
    _, err = model.RoleMenuModel.Delete(&pbs.PrivilegeRoleMenuDeleteReq{
        RoleKey:      roleKey,
        FuzzyRoleKey: false,
    })
    if err != nil {
        logger.Logger.Errorf("_roleService.AddMultiPower(%s, %+v) error: %v", roleKey, menuKeys, err)
        return nil, err
    }

    var wg sync.WaitGroup
    wg.Add(len(menuKeys))

    resp = &pbs.RolePowerMultiAddResult{
        Ids: make([]int32, 0),
    }
    for _, menuKey := range menuKeys {
        go func(wg *sync.WaitGroup, menuKey string) {
            rs, errs := model.RoleMenuModel.Save(&pbs.PrivilegeRoleMenuSaveReq{
                RoleKey:  roleKey,
                MenuKey:  menuKey,
                IsDelete: global.IsDeleteNormal,
            })
            defer wg.Done()
            if errs != nil {
                logger.Logger.Error("_roleService.AddMultiPower(%s, %s) error: %v", roleKey, menuKey, errs)
                return
            }
            resp.Ids = append(resp.Ids, rs.GetId())
        }(&wg, menuKey)
    }
    wg.Wait()

    return resp, nil
}
