package library

import (
    "github.com/gin-gonic/gin"
    "gitlab.ceibsmoment.com/c/mp/logger"
    "golang.org/x/time/rate"
    "net"
    "sync"
)

type _rateLimiter struct {
    r  rate.Limit
    b  int
    ip string
}

var (
    visitors = make(map[string]*rate.Limiter)
    mu       sync.Mutex
)

func NewRateLimiter(r rate.Limit, b int) *_rateLimiter {
    return &_rateLimiter{
        r: r,
        b: b,
    }
}

func (r *_rateLimiter)GetLimiter(ip string) *rate.Limiter {
    mu.Lock()
    defer mu.Unlock()

    if _, ok := visitors[ip]; !ok {
        visitors[ip] = rate.NewLimiter(r.r, r.b)
    }

    return visitors[ip]
}

func (*_rateLimiter) GetIp(c *gin.Context) (ip string, err error) {
    ip, _, err = net.SplitHostPort(c.Request.RemoteAddr)
    if err != nil {
        logger.Logger.Errorf("_rateLimiter.GetIp error: %v", err)
        return
    }

    return
}


