package main

import (
    "flag"
    "github.com/gin-gonic/gin"
    "gitlab.ceibsmoment.com/c/mp/config"
    "gitlab.ceibsmoment.com/c/mp/global"
    "gitlab.ceibsmoment.com/c/mp/logger"
    "gitlab.ceibsmoment.com/c/mp/middleware"
    "gitlab.ceibsmoment.com/c/mp/routes"
)

func main() {
    var env string
    ss := flag.String("env", global.EnvLocal, "")
    flag.Parse()
    env = *ss

    err := config.LoadCfg(env)
    if err != nil {
        panic(err)
    }

    router := gin.New()
    g := router.Group("/privilege")
    g.Use(gin.Logger(), gin.Recovery())

    // 捕获错误
    defer func() {
        if recv := recover(); recv != nil {
            logger.Logger.Errorf("recover panic error: %v", recv)
        }
    }()

    // 鉴权
    g.Use(middleware.Auth())

    // 注册Api
    routes.RegisterApi(g)

    if err := router.Run(":8080"); err != nil {
        panic("start error!")
    }
}
