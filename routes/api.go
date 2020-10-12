package routes

import (
    "github.com/gin-gonic/gin"
    "gitlab.ceibsmoment.com/c/mp/controller"
)

func RegisterApi(g *gin.RouterGroup) {
    // 用户
    g.POST("user/login", controller.UserController.Login)
    g.POST("user/add", controller.UserController.Add)
    g.POST("user/update", controller.UserController.Update)
    g.POST("user/list", controller.UserController.List)
    g.POST("user/logout", controller.UserController.Logout)
    g.POST("user/delete", controller.UserController.Delete)

    // 角色
    g.POST("role/list", controller.RoleController.List)
    g.POST("role/add", controller.RoleController.Add)
    g.POST("role/update", controller.RoleController.Update)
    g.POST("role/delete", controller.RoleController.Delete)

    // 菜单
    g.POST("menu/list", controller.MenuController.List)
    g.POST("menu/add", controller.MenuController.Add)
    g.POST("menu/update", controller.MenuController.Update)
    g.POST("menu/delete", controller.MenuController.Delete)
    g.POST("menu/list/tree/auth", controller.MenuController.ListAuthTree)
    g.POST("menu/list/tree/all", controller.MenuController.ListAllTree)
    g.POST("power/list/tree/auth", controller.PowerController.ListTree)

    // 角色与菜单映射关系
    g.POST("role/menu/list", controller.RoleMenuController.List)
    g.POST("role/menu/add", controller.RoleMenuController.Add)
    g.POST("role/menu/add/multi", controller.RoleMenuController.MultiAdd)
    g.POST("role/menu/update", controller.RoleMenuController.Update)
    g.POST("role/menu/delete", controller.RoleMenuController.Delete)

    // 前台菜单 & 后台路由映射关系保存
    g.POST("menu/back/url/list", controller.MenuBackUrlController.List)
    g.POST("menu/back/url/add", controller.MenuBackUrlController.Add)
    g.POST("menu/back/url/add/multi", controller.MenuBackUrlController.MultiAdd)
    g.POST("menu/back/url/update", controller.MenuBackUrlController.Update)
    g.POST("menu/back/url/delete", controller.MenuBackUrlController.Delete)
}
