package controller

import (
    "github.com/gin-gonic/gin"
    "gitlab.ceibsmoment.com/c/mp"
    def "gitlab.ceibsmoment.com/c/mp/global"
    "gitlab.ceibsmoment.com/c/mp/logger"
    "gitlab.ceibsmoment.com/c/mp/model"
    pbs "gitlab.ceibsmoment.com/c/mp/mppbs"
    "net/http"
    "sync"
)

type _roleMenuController struct {
}

var (
    RoleMenuController *_roleMenuController
)

// 角色前台路由关系列表
func (*_roleMenuController) List(c *gin.Context) {
    type params struct {
        Id           int32    `json:"id"`
        Page         int32    `json:"page"`
        Size         int32    `json:"size"`
        RoleKey      string   `json:"role_key"`
        FuzzyRoleKey int32    `json:"fuzzy_role_key"`
        MenuKey      string   `json:"menu_key"`
        FuzzyMenuKey int32    `json:"fuzzy_menu_key"`
        RoleKeys     []string `json:"role_keys"`
        MenuKeys     []string `json:"menu_keys"`
        IsDelete     []int32  `json:"is_delete"`
    }
    var in *params
    if err := c.ShouldBindJSON(&in); err != nil {
        logger.Logger.Errorf("_roleMenuController.List c.ShouldBindJSON(&in) error: %v", err)
        c.JSON(http.StatusOK, mp.JsonData(def.AdminCheckParamErr, "", nil))
        return
    }

    in.IsDelete = []int32{def.IsDeleteNormal}
    where := &pbs.PrivilegeRoleMenuFindWhere{
        RoleKey:      in.RoleKey,
        MenuKey:      in.MenuKey,
        FuzzyRoleKey: in.FuzzyRoleKey > 0,
        FuzzyMenuKey: in.FuzzyMenuKey > 0,
    }
    if in.Id > 0 {
        where.Id = []int32{in.Id}
    }

    if len(in.RoleKeys) > 0 {
        where.RoleKeys = in.RoleKeys
    }

    if len(in.MenuKeys) > 0 {
        where.MenuKeys = in.MenuKeys
    }

    if len(in.IsDelete) > 0 {
        where.IsDelete = in.IsDelete
    }

    resp, err := model.RoleMenuModel.Find(&pbs.PrivilegeRoleMenuFindReq{
        Page:  in.Page,
        Size:  in.Size,
        Where: where,
    })
    if err != nil {
        logger.Logger.Errorf("_roleMenuController.List model.RoleMenuModel.Find(%+v) error: %v", *in, err)
        c.JSON(http.StatusOK, mp.JsonData(def.AdminSystemErr, "", nil))
        return
    }

    _ = mp.WrapOutput(c, def.SuccessCode, "success", resp)
    return
}

// Add 角色前台路由关系添加
func (*_roleMenuController) Add(c *gin.Context) {
    type params struct {
        RoleKey string `json:"role_key"`
        MenuKey string `json:"menu_key"`
    }
    var in *params
    if err := c.ShouldBindJSON(&in); err != nil {
        logger.Logger.Errorf("_roleMenuController.Add c.ShouldBindJSON(&in) error: %v", err)
        c.JSON(http.StatusOK, mp.JsonData(def.AdminCheckParamErr, "", nil))
        return
    }

    resp, err := model.RoleMenuModel.Save(&pbs.PrivilegeRoleMenuSaveReq{
        RoleKey: in.RoleKey,
        MenuKey: in.MenuKey,
    })
    if err != nil {
        logger.Logger.Errorf("_roleMenuController.Add model.RoleMenuModel.Save(%+v) error: %v", *in, err)
        c.JSON(http.StatusOK, mp.JsonData(def.AdminSystemErr, "", nil))
        return
    }

    _ = mp.WrapOutput(c, def.SuccessCode, "success", resp)
    return
}

// MultiAdd 批量添加
func (*_roleMenuController) MultiAdd(c *gin.Context) {
    type params struct {
        RoleKey string   `json:"role_key"`
        MenuKey []string `json:"menu_key"`
    }
    var in *params
    if err := c.ShouldBindJSON(&in); err != nil {
        logger.Logger.Errorf("_roleMenuController.MultiAdd c.ShouldBindJSON(&in) error: %v", err)
        c.JSON(http.StatusOK, mp.JsonData(def.AdminCheckParamErr, "", nil))
        return
    }

    var wg sync.WaitGroup
    wg.Add(len(in.MenuKey))
    var affects []int32
    for _, v := range in.MenuKey {
        go func(v string) {
            resp, err := model.RoleMenuModel.Save(&pbs.PrivilegeRoleMenuSaveReq{
                RoleKey: in.RoleKey,
                MenuKey: v,
            })

            if err != nil {
                logger.Logger.Errorf("_roleMenuController.MultiAdd model.RoleMenuModel.Save(%s, %s) error: %v", in.RoleKey, v, err)
                return
            }
            affects = append(affects, resp.GetId())
            defer wg.Done()
        }(v)
    }
    wg.Wait()

    c.JSON(http.StatusOK, mp.JsonData(def.SuccessCode, "success", affects))
    return
}

// Update 角色前台路由关系修改
func (*_roleMenuController) Update(c *gin.Context) {
    type params struct {
        Id       int32  `json:"id"`
        RoleKey  string `json:"role_key"`
        MenuKey  string `json:"menu_key"`
        IsDelete int32  `json:"is_delete"`
    }
    var in *params
    if err := c.ShouldBindJSON(&in); err != nil {
        logger.Logger.Errorf("_roleMenuController.Update c.ShouldBindJSON(&in) error: %v", err)
        c.JSON(http.StatusOK, mp.JsonData(def.AdminCheckParamErr, "", nil))
        return
    }

    if in.Id == 0 {
        c.JSON(http.StatusOK, mp.JsonData(def.AdminCheckParamErr, "", nil))
        return
    }

    resp, err := model.RoleMenuModel.Save(&pbs.PrivilegeRoleMenuSaveReq{
        Id:       in.Id,
        RoleKey:  in.RoleKey,
        MenuKey:  in.MenuKey,
        IsDelete: in.IsDelete,
    })
    if err != nil {
        logger.Logger.Errorf("_roleMenuController.Update model.RoleMenuModel.Save(%+v) error: %v", *in, err)
        c.JSON(http.StatusOK, mp.JsonData(def.AdminSystemErr, "", nil))
        return
    }

    _ = mp.WrapOutput(c, def.SuccessCode, "success", resp)
    return
}

// Delete
// Warning 物理删除
func (*_roleMenuController) Delete(c *gin.Context) {
    type params struct {
        Id           int32    `json:"id"`
        RoleKey      string   `json:"role_key"`
        FuzzyRoleKey int32    `json:"fuzzy_role_key"`
        MenuKey      string   `json:"menu_key"`
        FuzzyMenuKey int32    `json:"fuzzy_menu_key"`
        RoleKeys     []string `json:"role_keys"`
        MenuKeys     []string `json:"menu_keys"`
        IsDelete     []int32  `json:"is_delete"`
    }
    var in *params
    if err := c.ShouldBindJSON(&in); err != nil {
        logger.Logger.Errorf("_roleMenuController.Delete c.ShouldBindJSON(&in) error: %v", err)
        c.JSON(http.StatusOK, mp.JsonData(def.AdminCheckParamErr, "", nil))
        return
    }

    where := &pbs.PrivilegeRoleMenuDeleteReq{
        RoleKey:      in.RoleKey,
        MenuKey:      in.MenuKey,
        FuzzyRoleKey: in.FuzzyRoleKey > 0,
        FuzzyMenuKey: in.FuzzyMenuKey > 0,
    }
    if in.Id > 0 {
        where.Id = []int32{in.Id}
    }

    if len(in.RoleKeys) > 0 {
        where.RoleKeys = in.RoleKeys
    }

    if len(in.MenuKeys) > 0 {
        where.MenuKeys = in.MenuKeys
    }

    if len(in.IsDelete) > 0 {
        where.IsDelete = in.IsDelete
    }
    resp, err := model.RoleMenuModel.Delete(where)
    if err != nil {
        logger.Logger.Errorf("_roleMenuController.List model.RoleMenuModel.Delete(%+v) error: %v", *where, err)
        c.JSON(http.StatusOK, mp.JsonData(def.AdminSystemErr, "", nil))
        return
    }

    _ = mp.WrapOutput(c, def.SuccessCode, "success", resp)
    return
}
