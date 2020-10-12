package global

import errors2 "errors"

type ErrorCode int32

const (
    IsDeleteNormal int32 = 2
    IsDeleteDel    int32 = 1

    //启动状态-未知
    StatusUnknown int32 = 0
    //未启动
    StatusUnStart int32 = 1
    //已启动
    StatusStarted int32 = 2
    //已下线
    StatusOffline int32 = 3

    //菜单可见性 - 可见
    MenuIsHiddenFalse int32 = 2
    //菜单可见性 - 不可见
    MenuIsHiddenYes int32 = 1

    SuccessCode        ErrorCode = 0 // 成功状态码
    AdminCheckParamErr ErrorCode = 10000
    AdminSystemErr     ErrorCode = 50000
    AdminSystemWarning ErrorCode = 53000

    // login error code
    AdminLoginSystemErr       ErrorCode = 51000
    AdminLoginUserPwdErr      ErrorCode = 51001
    AdminLoginUserNotFoundErr ErrorCode = 51002
    AdminLoginTokenCheckErr   ErrorCode = 51003

    // pb any error
    ProtoAnyErr ErrorCode = 60000

    // 鉴权错误
    PrivilegeAuthValidErr = 70000
    // 鉴权Header参数错误
    PrivilegeAuthParamsErr = 70001
)

func Errors() map[ErrorCode]string {
    return map[ErrorCode]string{
        AdminCheckParamErr:        "请求参数校验错误",
        AdminSystemErr:            "系统异常",
        AdminSystemWarning:        "系统警告",
        AdminLoginSystemErr:       "登录系统错误",
        AdminLoginUserPwdErr:      "账号密码错误",
        AdminLoginUserNotFoundErr: "管理员账号不存在",
        ProtoAnyErr:               "server error: [proto to any error]",
        AdminLoginTokenCheckErr:   "token验证错误",
        PrivilegeAuthValidErr:     "接口鉴权错误",
        PrivilegeAuthParamsErr:    "接口鉴权Header参数错误",
    }
}

func IsDeletes() map[int32]string {
    return map[int32]string{
        IsDeleteNormal: "正常",
        IsDeleteDel:    "删除",
    }
}

func AllStatus() map[int32]string {
    return map[int32]string{
        StatusUnknown: "--",
        StatusUnStart: "未上架",
        StatusStarted: "已上架",
        StatusOffline: "已下架",
    }
}

func StatusMsg(status int32) (msg string) {
    mps := AllStatus()
    if v, ok := mps[status]; ok {
        msg = v
    }
    return
}

func WrapError(code ErrorCode) error {
    msg := "Undefined error"
    errors := Errors()
    if v, ok := errors[code]; ok {
        msg = v
    }
    return errors2.New(msg)
}

func ErrMsg(code ErrorCode) (msg string) {
    errors := Errors()
    if v, ok := errors[code]; ok {
        msg = v
    }
    return
}

func IsDeleteMsg(code int32) (msg string) {
    mps := IsDeletes()
    if v, ok := mps[code]; ok {
        msg = v
    }
    return
}

