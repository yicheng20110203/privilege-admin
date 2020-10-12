package client

import (
    "github.com/go-redis/redis/v7"
    "gitlab.ceibsmoment.com/c/mp/config"
    "sync"
    "time"
)

var (
    _mux sync.RWMutex
)

func GetRedisClient(db int) *redis.Client {
    _mux.Lock()
    cfg := config.Cfg
    client := redis.NewClient(&redis.Options{
        Addr:     cfg.Redis.Host + ":" + cfg.Redis.Port,
        Password: cfg.Redis.Password,
        DB:       db,
    })
    defer _mux.Unlock()
    return client
}

func RedisSet(db int, key string, data interface{}, expireMinute int) (err error) {
    cl := GetRedisClient(db)
    defer cl.Close()
    err = cl.Set(key, data, time.Second*time.Duration(expireMinute)).Err()
    return
}

func RedisGet(db int, key string) (data string, err error) {
    cl := GetRedisClient(db)
    defer cl.Close()
    data, err = cl.Get(key).Result()
    return
}

func RedisDel(db int, key string) (err error) {
    cl := GetRedisClient(db)
    defer cl.Close()
    _, err = cl.Del(key).Result()
    return
}
