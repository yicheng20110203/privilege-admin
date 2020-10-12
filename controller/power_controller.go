package controller

import (
    "github.com/gin-gonic/gin"
    "gitlab.ceibsmoment.com/c/mp"
    def "gitlab.ceibsmoment.com/c/mp/global"
    "gitlab.ceibsmoment.com/c/mp/logger"
    pbs "gitlab.ceibsmoment.com/c/mp/mppbs"
    "gitlab.ceibsmoment.com/c/mp/service"
    "net/http"
)

type _powerController struct {
}

var (
    PowerController *_powerController
)

// ListTree
// get menu tree
func (obj *_powerController) ListTree(c *gin.Context) {
    authData, err := service.GetLoginInfoByHeader(c)
    if err != nil {
        logger.Logger.Errorf("_menuController.ListAuthTree service.GetLoginInfoByHeader error: %#v", err)
        c.JSON(http.StatusOK, mp.JsonData(def.AdminSystemErr, "", nil))
        return
    }
    userId := authData.UserId

    type params struct {
        ParentRoleKey string `json:"parent_role_key"`
    }
    var in *params
    if err := c.ShouldBindJSON(&in); err != nil {
        logger.Logger.Errorf("_powerController.ListTree c.ShouldBindJSON(&in) error: %v", err)
        c.JSON(http.StatusOK, mp.JsonData(def.AdminCheckParamErr, "", nil))
        return
    }

    resp, err := service.PowerService.GetPower(c, userId)
    if err != nil {
        logger.Logger.Errorf("_powerController.ListTree service.PowerService.GetPower error: %v", err)
        c.JSON(http.StatusOK, mp.JsonData(def.AdminSystemErr, "", nil))
        return
    }

    if in.ParentRoleKey == "" {
        //c.JSON(http.StatusOK, mp.JsonData(def.SuccessCode, "success", resp))
        _ = mp.WrapOutput(c, def.SuccessCode, "success", resp)
        return
    }

    mps, err := service.PowerService.GetPowerByRoleKey(c, in.ParentRoleKey)
    if err != nil {
        logger.Logger.Errorf("_powerController.ListTree service.PowerService.GetPowerByRoleKey error: %v", err)
        c.JSON(http.StatusOK, mp.JsonData(def.AdminCheckParamErr, "", nil))
        return
    }

    resp.Items = obj._process(resp.GetItems(), mps)
    _ = mp.WrapOutput(c, def.SuccessCode, "success", resp)
    return
}

func (obj *_powerController) _process(items []*pbs.PowerItem, mps map[string]string) []*pbs.PowerItem {
    for k, v := range items {
        if _, ok := mps[v.GetMenuKey()]; ok {
            v.Selected = true
            if v.Children != nil && len(v.Children) > 0 {
                v.Children = obj._process(v.Children, mps)
            }
            items[k] = v
        }
    }
    return items
}
