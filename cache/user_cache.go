package cache

import (
    "encoding/json"
    "errors"
    "fmt"
    "github.com/gin-gonic/gin"
    "gitlab.ceibsmoment.com/c/mp/client"
    "gitlab.ceibsmoment.com/c/mp/config"
    "gitlab.ceibsmoment.com/c/mp/global"
    "gitlab.ceibsmoment.com/c/mp/library"
    "gitlab.ceibsmoment.com/c/mp/logger"
    "time"
)

const (
    SessionDbIndex = 1
)

type _userCache struct {
}

var (
    UserCache *_userCache
)

// admin login token
func (*_userCache) BuildAdminToken(os string, uid int32) string {
    var s string
    s = fmt.Sprintf("admin:%s:%d:%d", os, uid, time.Now().Unix())
    return library.MD5(s)
}

func (obj *_userCache) buildAdminCacheKey(os string, uid int32) string {
    return fmt.Sprintf("admin:%s:%d", os, uid)
}

func (obj *_userCache) SetUserCache(os string, uid int32) (data global.UserCacheInfo, err error) {
    token := obj.BuildAdminToken(os, uid)
    key := obj.buildAdminCacheKey(os, uid)
    data = global.UserCacheInfo{
        Id:    uid,
        Token: token,
    }
    bs, err := json.Marshal(data)
    err = client.RedisSet(SessionDbIndex, key, string(bs), 86400)
    if err != nil {
        logger.Logger.Errorf("_userCache.SetUserCache RedisSet<%v, %s> error: %v", key, string(bs), err)
        return
    }
    return
}

func (obj *_userCache) GetUserCache(os string, uid int32) (user global.UserCacheInfo, err error) {
    key := obj.buildAdminCacheKey(os, uid)
    r, err := client.RedisGet(SessionDbIndex, key)
    if err != nil {
        logger.Logger.Errorf("_userCache.GetUserCache RedisGet(%v) error: %v", key, err)
        return
    }

    err = json.Unmarshal([]byte(r), &user)
    if err != nil {
        logger.Logger.Errorf("_userCache.GetUserCache Unmarshal error: %v", err)
        return
    }
    return
}

func (obj *_userCache) CheckToken(os string, uid int32, token string) bool {
    user, err := obj.GetUserCache(os, uid)
    if err != nil {
        logger.Logger.Errorf("_userCache.CheckToken error : %v", err)
        return false
    }

    if user.Id != uid || user.Token != token {
        logger.Logger.Errorf("_userCache.CheckToken check error")
        return false
    }

    return true
}

type _loginInfo struct {
    Token  string `json:"token"`
    UserId int32  `json:"user_id"`
    Os     string `json:"os"`
}

func (obj *_userCache) _getLoginInfo(c *gin.Context) (authInfo *_loginInfo, err error) {
    authorization := c.Request.Header.Get(global.KeyAuthorization)
    if authorization == "" {
        logger.Logger.Info("service.Auth info: authorization = %s", authorization)
        return &_loginInfo{}, nil
    }

    decryptData, err := library.AesUtil.DecryptString(config.Cfg.Aes.Key, authorization)
    if err != nil {
        logger.Logger.Errorf("service.Auth error: library.AesUtil.DecryptString(%s) error: %#v", authorization, err)
        return nil, err
    }

    dataBytes, _ := json.Marshal(decryptData)
    var resp *_loginInfo
    err = json.Unmarshal(dataBytes, &resp)
    if err != nil {
        logger.Logger.Errorf("service.Auth service json.Unmarshal(%s) error: %#v", string(dataBytes), err)
        err = errors.New("json反序列化错误")
        return nil, err
    }

    return resp, nil
}

func (obj *_userCache) GetUserFromCacheAndCheck(c *gin.Context) (user global.UserCacheInfo, err error) {
    authData, err := obj._getLoginInfo(c)
    if err != nil {
        logger.Logger.Errorf("_userCache.GetUserFromCacheAndCheck obj._getLoginInfo error: %#v", err)
        return global.UserCacheInfo{}, err
    }
    token := authData.Token
    uid := authData.UserId
    os := authData.Os

    check := obj.CheckToken(os, uid, token)
    if !check {
        return global.UserCacheInfo{}, global.WrapError(global.AdminLoginTokenCheckErr)
    }

    user = global.UserCacheInfo{
        Id:    uid,
        Token: token,
    }
    return
}

func (obj *_userCache) CleanToken(os string, uid int32, token string) bool {
    if ck := obj.CheckToken(os, uid, token); ck {
        key := obj.buildAdminCacheKey(os, uid)
        err := client.RedisDel(SessionDbIndex, key)
        if err != nil {
            logger.Logger.Errorf("_userCache.CleanToken(%s) error: %v", key, err)
            return false
        }
        return true
    }

    return false
}
