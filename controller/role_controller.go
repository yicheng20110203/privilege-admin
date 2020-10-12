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

type _roleController struct {
}

var (
    RoleController *_roleController
)

// 角色列表
func (*_roleController) List(c *gin.Context) {
    type params struct {
        Page         int32  `json:"page"`
        Size         int32  `json:"size"`
        Name         string `json:"name"`
        FuzzyName    int32  `json:"fuzzy_name"`
        Desc         string `json:"desc"`
        Status       int32  `json:"status"`
        RoleKey      string `json:"role_key"`
        FuzzyRoleKey int32  `json:"fuzzy_role_key"`
        Id           int32  `json:"id"`
    }
    var in *params
    if err := c.ShouldBindJSON(&in); err != nil {
        logger.Logger.Errorf("_roleController.List c.ShouldBindJSON(&in) error: %v", err)
        c.JSON(http.StatusOK, mp.JsonData(def.AdminCheckParamErr, "", nil))
        return
    }

    where := &pbs.PrivilegeRoleFindWhere{
        Name:         in.Name,
        FuzzyName:    in.FuzzyName > 0,
        RoleKey:      in.RoleKey,
        FuzzyRoleKey: in.FuzzyRoleKey > 0,
        Desc:         in.Desc,
    }
    if in.Id > 0 {
        where.Id = []int32{in.Id}
    }

    if in.Status > 0 {
        where.Status = []int32{in.Status}
    }
    resp, err := model.RoleModel.Find(&pbs.PrivilegeRoleFindReq{
        Page:  in.Page,
        Size:  in.Size,
        Where: where,
    })
    if err != nil {
        logger.Logger.Errorf("_roleController.List model.RoleModel.Find(%+v) error: %v", *in, err)
        c.JSON(http.StatusOK, mp.JsonData(def.AdminSystemErr, "", nil))
        return
    }

    roles := make([]*pbs.Role, 0)
    for _, v := range resp.GetList() {
        roles = append(roles, &pbs.Role{
            Id:         v.GetId(),
            Name:       v.GetName(),
            Desc:       v.GetDesc(),
            RoleKey:    v.GetRoleKey(),
            Status:     v.GetStatus(),
            StatusDesc: service.RoleService.GetStatusDesc(v.GetStatus()),
            CreateTime: v.GetCreateTime(),
            UpdateTime: v.GetUpdateTime(),
        })
    }

    _ = mp.WrapOutput(c, def.SuccessCode, "success", &pbs.RoleList{
        Page:      resp.GetPage(),
        Size:      resp.GetSize(),
        TotalPage: resp.GetTotalPage(),
        TotalSize: resp.GetTotalSize(),
        List:      roles,
    })
    return
}

// Add 角色添加
func (*_roleController) Add(c *gin.Context) {
    type params struct {
        Name          string   `json:"name"`
        Desc          string   `json:"desc"`
        ParentRoleKey string   `json:"parent_role_key"`
        MenuKeys      []string `json:"menu_keys"`
    }
    var in *params
    if err := c.ShouldBindJSON(&in); err != nil {
        logger.Logger.Errorf("_roleController.Add c.ShouldBindJSON(&in) error: %v", err)
        c.JSON(http.StatusOK, mp.JsonData(def.AdminCheckParamErr, "", nil))
        return
    }

    if len(in.MenuKeys) == 0 || in.Name == "" {
        logger.Logger.Errorf("_roleController.Add params(&in) check error: %v", "参数检查错误")
        c.JSON(http.StatusOK, mp.JsonData(def.AdminCheckParamErr, "", nil))
        return
    }

    // 保存role
    resp, err := model.RoleModel.Save(&pbs.PrivilegeRoleSaveReq{
        Name:    in.Name,
        Desc:    in.Desc,
        RoleKey: in.ParentRoleKey,
        Status:  def.StatusStarted,
    })
    if err != nil {
        logger.Logger.Errorf("_roleController.Add model.RoleModel.Save(%+v) error: %v", *in, err)
        c.JSON(http.StatusOK, mp.JsonData(def.AdminSystemErr, "", nil))
        return
    }

    // 获取role info
    roleInfos, err := model.RoleModel.Find(&pbs.PrivilegeRoleFindReq{
        Page: 1,
        Size: 1,
        Where: &pbs.PrivilegeRoleFindWhere{
            Id: []int32{
                resp.GetId(),
            },
        },
    })
    if err != nil || roleInfos == nil {
        logger.Logger.Errorf("_roleController.Add model.RoleModel.Find(%d) error: %v", resp.GetId(), err)
        c.JSON(http.StatusOK, mp.JsonData(def.AdminSystemErr, "", nil))
        return
    }

    rs, err := service.RoleService.AddMultiPower(c, roleInfos.GetList()[0].GetRoleKey(), in.MenuKeys)
    if err != nil {
        logger.Logger.Errorf("_roleController.Add service.RoleService.AddMultiPower(%d, %+v) error: %v", roleInfos.GetList()[0].GetRoleKey(), in.MenuKeys, err)
        c.JSON(http.StatusOK, mp.JsonData(def.AdminSystemErr, "", nil))
    }

    _ = mp.WrapOutput(c, def.SuccessCode, "success", rs)
    return
}

// Update 角色修改
func (*_roleController) Update(c *gin.Context) {
    type params struct {
        Id       int32    `json:"id"`
        Name     string   `json:"name"`
        Desc     string   `json:"desc"`
        MenuKeys []string `json:"menu_keys"`
    }
    var in *params
    if err := c.ShouldBindJSON(&in); err != nil {
        logger.Logger.Errorf("_roleController.Update c.ShouldBindJSON(&in) error: %v", err)
        c.JSON(http.StatusOK, mp.JsonData(def.AdminCheckParamErr, "", nil))
        return
    }

    if in.Id == 0 {
        c.JSON(http.StatusOK, mp.JsonData(def.AdminCheckParamErr, "", nil))
        return
    }

    _, err := model.RoleModel.Save(&pbs.PrivilegeRoleSaveReq{
        Id:   in.Id,
        Name: in.Name,
        Desc: in.Desc,
    })
    if err != nil {
        logger.Logger.Errorf("_roleController.Update model.RoleModel.Save(%+v) error: %v", *in, err)
        c.JSON(http.StatusOK, mp.JsonData(def.AdminSystemErr, "", nil))
        return
    }

    // 获取role info
    roleInfos, err := model.RoleModel.Find(&pbs.PrivilegeRoleFindReq{
        Page: 1,
        Size: 1,
        Where: &pbs.PrivilegeRoleFindWhere{
            Id: []int32{in.Id},
        },
    })
    if err != nil || roleInfos == nil {
        logger.Logger.Errorf("_roleController.Update model.RoleModel.Find(%+v) error: %v", *in, err)
        c.JSON(http.StatusOK, mp.JsonData(def.AdminSystemErr, "", nil))
        return
    }

    rs, err := service.RoleService.AddMultiPower(c, roleInfos.GetList()[0].GetRoleKey(), in.MenuKeys)
    if err != nil {
        logger.Logger.Errorf("_roleController.Update model.RoleModel.Find(%+v) error: %v", *in, err)
        c.JSON(http.StatusOK, mp.JsonData(def.AdminSystemErr, "", nil))
        return
    }
    _ = mp.WrapOutput(c, def.SuccessCode, "success", rs)
    return
}

// Delete
// Warning 物理删除
func (*_roleController) Delete(c *gin.Context) {
    type params struct {
        Name         string `json:"name"`
        FuzzyName    int32  `json:"fuzzy_name"`
        Status       int32  `json:"status"`
        RoleKey      string `json:"role_key"`
        FuzzyRoleKey int32  `json:"fuzzy_role_key"`
        Id           int32  `json:"id"`
    }
    var in *params
    if err := c.ShouldBindJSON(&in); err != nil {
        logger.Logger.Errorf("_menuBackUrlController.Delete c.ShouldBindJSON(&in) error: %v", err)
        c.JSON(http.StatusOK, mp.JsonData(def.AdminCheckParamErr, "", nil))
        return
    }

    where := &pbs.PrivilegeRoleDeleteReq{
        Name:         in.Name,
        FuzzyName:    in.FuzzyName > 0,
        RoleKey:      in.RoleKey,
        FuzzyRoleKey: in.FuzzyRoleKey > 0,
    }
    if in.Id > 0 {
        where.Id = []int32{in.Id}
    }

    if in.Status > 0 {
        where.Status = []int32{in.Status}
    }
    resp, err := model.RoleModel.Delete(where)
    if err != nil {
        logger.Logger.Errorf("_roleController.List model.RoleModel.Delete(%+v) error: %v", *where, err)
        c.JSON(http.StatusOK, mp.JsonData(def.AdminSystemErr, "", nil))
        return
    }

    _ = mp.WrapOutput(c, def.SuccessCode, "success", resp)
    return
}
