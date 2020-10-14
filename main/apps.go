package main

import (
    "gitlab.ceibsmoment.com/c/mp/logger"
)

type DyContext struct {
    Name string `json:"name"`
}

type HandleFunc func(c *DyContext)

type Middleware HandleFunc

func main() {
    c := &DyContext{
        Name: "cyy",
    }
    LoggerMiddleware()(c)
    RecoverMiddleware()(c)
}

func LoggerMiddleware() Middleware {
    return func(c *DyContext) {
        logger.Logger.Infof("newLogger middleware with name = %s", c.Name)
    }
}

func RecoverMiddleware() Middleware {
    return func(c *DyContext) {
        defer func() {
            if err := recover(); err != nil {
                logger.Logger.Errorf("catch panic error in recover: %v", err)
            }
        }()
        c.Next()
    }
}

func (c *DyContext) Next() {
    logger.Logger.Infof("DyContext exec .... %s", c.Name)
    panic("c.Next panic error")
}
