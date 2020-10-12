package middleware

import (
    "github.com/gin-gonic/gin"
    "gitlab.ceibsmoment.com/c/mp"
    "gitlab.ceibsmoment.com/c/mp/config"
    def "gitlab.ceibsmoment.com/c/mp/global"
    "gitlab.ceibsmoment.com/c/mp/logger"
    "gitlab.ceibsmoment.com/c/mp/service"
    "net/http"
)

func Auth() gin.HandlerFunc {
    return func(c *gin.Context) {
        var skip bool
        skip = config.Cfg.Auth.Skip == "true"
        authPass, err := service.Auth(c, "privilege", skip)
        if err != nil {
            logger.Logger.Errorf("middleware.Auth service.Auth error: %v", err)
            c.JSON(http.StatusOK, mp.JsonData(def.PrivilegeAuthValidErr, "", nil))
            c.Abort()
            return
        }

        if authPass {
            c.Next()
            return
        }

        c.JSON(http.StatusOK, mp.JsonData(def.PrivilegeAuthValidErr, "", nil))
        c.Abort()
        return
    }
}
