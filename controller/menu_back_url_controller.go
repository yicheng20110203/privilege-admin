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

type _menuBackUrlController struct {
}

var (
    MenuBackUrlController *_menuBackUrlController
)

// 前台路由与后台路由映射关系列表
func (*_menuBackUrlController) List(c *gin.Context) {
    type params struct {
        Id           int32   `json:"id"`
        Page         int32   `json:"page"`
        Size         int32   `json:"size"`
        MenuKey      string  `json:"menu_key"`
        FuzzyMenuKey int32   `json:"fuzzy_menu_key"`
        BackUrl      string  `json:"back_url"`
        FuzzyBackUrl int32   `json:"fuzzy_back_url"`
        Desc         string  `json:"desc"`
        IsDelete     []int32 `json:"is_delete"`
    }
    var in *params
    if err := c.ShouldBindJSON(&in); err != nil {
        logger.Logger.Errorf("_menuBackUrlController.List c.ShouldBindJSON(&in) error: %v", err)
        c.JSON(http.StatusOK, mp.JsonData(def.AdminCheckParamErr, "", nil))
        return
    }

    where := &pbs.PrivilegeMenuBackUrlFindWhere{
        BackUrl:      in.BackUrl,
        FuzzyBackUrl: in.FuzzyBackUrl > 0,
        MenuKey:      in.MenuKey,
        FuzzyMenuKey: in.FuzzyMenuKey > 0,
        Desc:         in.Desc,
    }
    if in.Id > 0 {
        where.Id = []int32{in.Id}
    }

    if len(in.IsDelete) > 0 {
        where.IsDelete = in.IsDelete
    }

    resp, err := model.MenuBackUrlModel.Find(&pbs.PrivilegeMenuBackUrlFindReq{
        Page:  in.Page,
        Size:  in.Size,
        Where: where,
    })
    if err != nil {
        logger.Logger.Errorf("_menuBackUrlController.List model.MenuBackUrlModel.Find(%+v) error: %v", *in, err)
        c.JSON(http.StatusOK, mp.JsonData(def.AdminSystemErr, "", nil))
        return
    }

    _ = mp.WrapOutput(c, def.SuccessCode, "success", resp)
    return
}

// Add 前台路由与后台路由映射关系添加
func (*_menuBackUrlController) Add(c *gin.Context) {
    type params struct {
        MenuKey string `json:"menu_key"`
        BackUrl string `json:"back_url"`
        Desc    string `json:"desc"`
    }
    var in *params
    if err := c.ShouldBindJSON(&in); err != nil {
        logger.Logger.Errorf("_menuBackUrlController.Add c.ShouldBindJSON(&in) error: %v", err)
        c.JSON(http.StatusOK, mp.JsonData(def.AdminCheckParamErr, "", nil))
        return
    }

    resp, err := model.MenuBackUrlModel.Save(&pbs.PrivilegeMenuBackUrlSaveReq{
        MenuKey: in.MenuKey,
        BackUrl: in.BackUrl,
        Desc:    in.Desc,
    })
    if err != nil {
        logger.Logger.Errorf("_menuBackUrlController.Add model.MenuBackUrlModel.Save(%+v) error: %v", *in, err)
        c.JSON(http.StatusOK, mp.JsonData(def.AdminSystemErr, "", nil))
        return
    }

    _ = mp.WrapOutput(c, def.SuccessCode, "success", resp)
    return
}

// MultiAdd 批量添加
func (*_menuBackUrlController) MultiAdd(c *gin.Context) {
    type params struct {
        MenuKey string   `json:"menu_key"`
        BackUrl []string `json:"back_url"`
        Desc    string   `json:"desc"`
    }
    var in *params
    if err := c.ShouldBindJSON(&in); err != nil {
        logger.Logger.Errorf("_menuBackUrlController.MultiAdd c.ShouldBindJSON(&in) error: %v", err)
        c.JSON(http.StatusOK, mp.JsonData(def.AdminCheckParamErr, "", nil))
        return
    }

    var wg sync.WaitGroup
    wg.Add(len(in.BackUrl))
    var affects []int32
    for _, v := range in.BackUrl {
        go func(v string) {
            resp, err := model.MenuBackUrlModel.Save(&pbs.PrivilegeMenuBackUrlSaveReq{
                MenuKey: in.MenuKey,
                BackUrl: v,
                Desc:    in.Desc,
            })

            if err != nil {
                logger.Logger.Errorf("_menuBackUrlController.MultiAdd model.MenuBackUrlModel.Save(%s, %s) error: %v", in.MenuKey, v, err)
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

// Update 前台路由与后台路由映射关系修改
func (*_menuBackUrlController) Update(c *gin.Context) {
    type params struct {
        Id       int32  `json:"id"`
        MenuKey  string `json:"menu_key"`
        BackUrl  string `json:"back_url"`
        Desc     string `json:"desc"`
        IsDelete int32  `json:"is_delete"`
    }
    var in *params
    if err := c.ShouldBindJSON(&in); err != nil {
        logger.Logger.Errorf("_menuBackUrlController.Update c.ShouldBindJSON(&in) error: %v", err)
        c.JSON(http.StatusOK, mp.JsonData(def.AdminCheckParamErr, "", nil))
        return
    }

    if in.Id == 0 {
        c.JSON(http.StatusOK, mp.JsonData(def.AdminCheckParamErr, "", nil))
        return
    }

    resp, err := model.MenuBackUrlModel.Save(&pbs.PrivilegeMenuBackUrlSaveReq{
        Id:       in.Id,
        MenuKey:  in.MenuKey,
        BackUrl:  in.BackUrl,
        Desc:     in.Desc,
        IsDelete: in.IsDelete,
    })
    if err != nil {
        logger.Logger.Errorf("_menuBackUrlController.Update model.MenuBackUrlModel.Save(%+v) error: %v", *in, err)
        c.JSON(http.StatusOK, mp.JsonData(def.AdminSystemErr, "", nil))
        return
    }

    _ = mp.WrapOutput(c, def.SuccessCode, "success", resp)
    return
}

// Delete
// Warning 物理删除
func (*_menuBackUrlController) Delete(c *gin.Context) {
    type params struct {
        Id           int32   `json:"id"`
        MenuKey      string  `json:"menu_key"`
        FuzzyMenuKey int32   `json:"fuzzy_menu_key"`
        BackUrl      string  `json:"back_url"`
        FuzzyBackUrl int32   `json:"fuzzy_back_url"`
        Desc         string  `json:"desc"`
        IsDelete     []int32 `json:"is_delete"`
    }
    var in *params
    if err := c.ShouldBindJSON(&in); err != nil {
        logger.Logger.Errorf("_menuBackUrlController.Delete c.ShouldBindJSON(&in) error: %v", err)
        c.JSON(http.StatusOK, mp.JsonData(def.AdminCheckParamErr, "", nil))
        return
    }

    where := &pbs.PrivilegeMenuBackUrlDeleteReq{
        BackUrl:      in.BackUrl,
        FuzzyBackUrl: in.FuzzyBackUrl > 0,
        MenuKey:      in.MenuKey,
        FuzzyMenuKey: in.FuzzyMenuKey > 0,
        Desc:         in.Desc,
    }
    if in.Id > 0 {
        where.Id = []int32{in.Id}
    }

    if len(in.IsDelete) > 0 {
        where.IsDelete = in.IsDelete
    }

    resp, err := model.MenuBackUrlModel.Delete(where)
    if err != nil {
        logger.Logger.Errorf("_menuBackUrlController.Delete model.MenuBackUrlModel.Delete(%+v) error: %v", *where, err)
        c.JSON(http.StatusOK, mp.JsonData(def.AdminSystemErr, "", nil))
        return
    }

    _ = mp.WrapOutput(c, def.SuccessCode, "success", resp)
    return
}
