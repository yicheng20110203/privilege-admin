package controller

import (
    "github.com/gin-gonic/gin"
    "gitlab.ceibsmoment.com/c/mp"
    def "gitlab.ceibsmoment.com/c/mp/global"
    "gitlab.ceibsmoment.com/c/mp/logger"
    "gitlab.ceibsmoment.com/c/mp/model"
    pbs "gitlab.ceibsmoment.com/c/mp/mppbs"
    "gitlab.ceibsmoment.com/c/mp/service"
    "net/http"
)

type _menuController struct {
}

var (
    MenuController *_menuController
)

// 菜单列表
func (*_menuController) List(c *gin.Context) {
    type params struct {
        Id           int32    `json:"id"`
        Page         int32    `json:"page"`
        Size         int32    `json:"size"`
        Path         string   `json:"path"`
        Component    string   `json:"component"`
        Title        string   `json:"title"`
        Name         string   `json:"name"`
        Icon         string   `json:"icon"`
        MenuKey      string   `json:"menu_key"`
        FuzzyMenuKey int32    `json:"fuzzy_menu_key"`
        Level        int32    `json:"level"`
        MenuKeys     []string `json:"menu_keys"`
        IsDelete     int32    `json:"is_delete"`
    }
    var in *params
    if err := c.ShouldBindJSON(&in); err != nil {
        logger.Logger.Errorf("_menuController.List c.ShouldBindJSON(&in) error: %v", err)
        c.JSON(http.StatusOK, mp.JsonData(def.AdminCheckParamErr, "", nil))
        return
    }

    in.IsDelete = def.IsDeleteNormal
    where := &pbs.PrivilegeMenuFindWhere{
        Path:         in.Path,
        Title:        in.Title,
        Name:         in.Name,
        Icon:         in.Icon,
        MenuKey:      in.MenuKey,
        FuzzyMenuKey: in.FuzzyMenuKey > 0,
    }
    if in.Id > 0 {
        where.Id = []int32{in.Id}
    }

    if in.Component != "" {
        where.Component = []string{in.Component}
    }

    if in.Level > 0 {
        where.Level = []int32{in.Level}
    }

    if len(in.MenuKeys) > 0 {
        where.MenuKeys = in.MenuKeys
    }

    if in.IsDelete != 0 {
        where.IsDelete = []int32{in.IsDelete}
    } else {
        where.IsDelete = []int32{def.IsDeleteNormal, def.IsDeleteDel}
    }

    resp, err := model.MenuModel.Find(&pbs.PrivilegeMenuFindReq{
        Page:  in.Page,
        Size:  in.Size,
        Where: where,
    })
    if err != nil {
        logger.Logger.Errorf("_menuController.List model.MenuModel.Find(%+v) error: %v", *in, err)
        c.JSON(http.StatusOK, mp.JsonData(def.AdminSystemErr, "", nil))
        return
    }

    _ = mp.WrapOutput(c, def.SuccessCode, "success", resp)
    return
}

// Add 菜单添加
func (*_menuController) Add(c *gin.Context) {
    type params struct {
        Path         string `json:"path"`
        Component    string `json:"component"`
        Title        string `json:"title"`
        Name         string `json:"name"`
        Icon         string `json:"icon"`
        MenuKey      string `json:"menu_key"`
        DisplayOrder int32  `json:"display_order"`
    }
    var in *params
    if err := c.ShouldBindJSON(&in); err != nil {
        logger.Logger.Errorf("_menuController.Add c.ShouldBindJSON(&in) error: %v", err)
        c.JSON(http.StatusOK, mp.JsonData(def.AdminCheckParamErr, "", nil))
        return
    }

    resp, err := model.MenuModel.Save(&pbs.PrivilegeMenuSaveReq{
        Path:         in.Path,
        Component:    in.Component,
        Title:        in.Title,
        Name:         in.Name,
        Icon:         in.Icon,
        MenuKey:      in.MenuKey,
        DisplayOrder: in.DisplayOrder,
        IsDelete:     def.IsDeleteNormal,
        IsHidden:     def.MenuIsHiddenFalse,
    })
    if err != nil {
        logger.Logger.Errorf("_menuController.Add model.MenuModel.Save(%+v) error: %v", *in, err)
        c.JSON(http.StatusOK, mp.JsonData(def.AdminSystemErr, "", nil))
        return
    }

    _ = mp.WrapOutput(c, def.SuccessCode, "success", resp)
    return
}

// Update 菜单修改
func (*_menuController) Update(c *gin.Context) {
    type params struct {
        Id           int32  `json:"id"`
        Path         string `json:"path"`
        Component    string `json:"component"`
        Title        string `json:"title"`
        Name         string `json:"name"`
        Icon         string `json:"icon"`
        MenuKey      string `json:"menu_key"`
        DisplayOrder int32  `json:"display_order"`
        IsHidden     int32  `json:"is_hidden"`
    }
    var in *params
    if err := c.ShouldBindJSON(&in); err != nil {
        logger.Logger.Errorf("_menuController.Update c.ShouldBindJSON(&in) error: %v", err)
        c.JSON(http.StatusOK, mp.JsonData(def.AdminCheckParamErr, "", nil))
        return
    }

    if in.Id == 0 {
        c.JSON(http.StatusOK, mp.JsonData(def.AdminCheckParamErr, "", nil))
        return
    }

    resp, err := model.MenuModel.Save(&pbs.PrivilegeMenuSaveReq{
        Id:           in.Id,
        Path:         in.Path,
        Title:        in.Title,
        Name:         in.Name,
        Icon:         in.Icon,
        MenuKey:      in.MenuKey,
        DisplayOrder: in.DisplayOrder,
        IsHidden:     in.IsHidden,
    })
    if err != nil {
        logger.Logger.Errorf("_menuController.Update model.MenuModel.Save(%+v) error: %v", *in, err)
        c.JSON(http.StatusOK, mp.JsonData(def.AdminSystemErr, "", nil))
        return
    }

    _ = mp.WrapOutput(c, def.SuccessCode, "success", resp)
    return
}

// ListTree
// get menu tree
func (*_menuController) ListAuthTree(c *gin.Context) {
    authData, err := service.GetLoginInfoByHeader(c)
    if err != nil {
        logger.Logger.Errorf("_menuController.ListAuthTree service.GetLoginInfoByHeader error: %#v", err)
        c.JSON(http.StatusOK, mp.JsonData(def.AdminSystemErr, "", nil))
        return
    }

    userId := authData.UserId
    resp, err := service.MenuService.GetMenu(c, int32(userId))
    if err != nil {
        logger.Logger.Errorf("_menuController.ListAuthTree service.MenuService.GetMenu error: %v", err)
        c.JSON(http.StatusOK, mp.JsonData(def.AdminSystemErr, "", nil))
        return
    }

    c.JSON(http.StatusOK, mp.JsonData(def.SuccessCode, "success", resp))

    //_ = mp.WrapOutput(c, def.SuccessCode, "success", resp)
    return
}

func (*_menuController) ListAllTree(c *gin.Context) {
    resp, err := service.MenuService.GetAllMenu(c)
    if err != nil {
        logger.Logger.Errorf("_menuController.ListAllTree service.MenuService.GetAllMenu error: %v", err)
        c.JSON(http.StatusOK, mp.JsonData(def.AdminSystemErr, "", nil))
        return
    }

    c.JSON(http.StatusOK, mp.JsonData(def.SuccessCode, "success", resp))
    return
}

// Delete
// Warning 物理删除
func (*_menuController) Delete(c *gin.Context) {
    type params struct {
        Id           int32    `json:"id"`
        Path         string   `json:"path"`
        Component    string   `json:"component"`
        Title        string   `json:"title"`
        Name         string   `json:"name"`
        Icon         string   `json:"icon"`
        MenuKey      string   `json:"menu_key"`
        FuzzyMenuKey int32    `json:"fuzzy_menu_key"`
        Level        int32    `json:"level"`
        MenuKeys     []string `json:"menu_keys"`
        IsDelete     int32    `json:"is_delete"`
    }
    var in *params
    if err := c.ShouldBindJSON(&in); err != nil {
        logger.Logger.Errorf("_menuController.Delete c.ShouldBindJSON(&in) error: %v", err)
        c.JSON(http.StatusOK, mp.JsonData(def.AdminCheckParamErr, "", nil))
        return
    }

    where := &pbs.PrivilegeMenuDeleteReq{
        Path:         in.Path,
        Title:        in.Title,
        Name:         in.Name,
        Icon:         in.Icon,
        MenuKey:      in.MenuKey,
        FuzzyMenuKey: in.FuzzyMenuKey > 0,
    }
    if in.Id > 0 {
        where.Id = []int32{in.Id}
    }

    if in.Component != "" {
        where.Component = []string{in.Component}
    }

    if in.Level > 0 {
        where.Level = []int32{in.Level}
    }

    if len(in.MenuKeys) > 0 {
        where.MenuKeys = in.MenuKeys
    }

    if in.IsDelete != 0 {
        where.IsDelete = []int32{in.IsDelete}
    } else {
        where.IsDelete = []int32{def.IsDeleteNormal, def.IsDeleteDel}
    }

    resp, err := model.MenuModel.Delete(where)
    if err != nil {
        logger.Logger.Errorf("_menuController.List model.MenuModel.Delete(%+v) error: %v", *where, err)
        c.JSON(http.StatusOK, mp.JsonData(def.AdminSystemErr, "", nil))
        return
    }

    _ = mp.WrapOutput(c, def.SuccessCode, "success", resp)
    return
}
