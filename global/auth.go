package global

const (
    JumpTypeDefault = 0
    JumpTypeCanJump = 1
    JumpTypeNotJump = 2

    KeyXToken        = "X-Token"
    KeyXUserId       = "X-User-Id"
    KeyXOs           = "X-Os"
    // 解决跨域自定义header非标字段
    KeyAuthorization = "Authorization"

    LoginURI = "/privilege/user/login"
)

var (
    WhiteList = []string{
        LoginURI,
    }
)
