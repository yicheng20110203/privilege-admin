package global

const (
    RegPlatformAndroid    int32 = 1
    RegPlatformIos        int32 = 2
    RegPlatformH5         int32 = 3
    RegPlatformPc         int32 = 4
    RegPlatformMinProgram int32 = 5

    RegTypeUserPwd    int32 = 1 //用户名+密码注册
    RegTypeMobileCode int32 = 2 //手机号+验证码注册
    RegTypeWx         int32 = 3 //微信注册
    RegTypeWeibo      int32 = 4 //微博注册
    RegTypeH5Login    int32 = 5 //h5注册
    RegTypeEmail      int32 = 6 //邮箱注册
)

func RegPlatforms() map[int32]string {
    return map[int32]string{
        RegPlatformAndroid:    "android",
        RegPlatformIos:        "ios",
        RegPlatformH5:         "h5",
        RegPlatformPc:         "pc",
        RegPlatformMinProgram: "小程序",
    }
}

func RegTypes() map[int32]string {
    return map[int32]string{
        RegTypeUserPwd:    "账号+密码",
        RegTypeMobileCode: "手机号+验证码",
        RegTypeWx:         "微信",
        RegTypeWeibo:      "微博",
        RegTypeH5Login:    "h5",
        RegTypeEmail:      "邮箱",
    }
}

func RegPlatformDesc(regPlatform int32) string {
    all := RegPlatforms()
    for k, v := range all {
        if k == regPlatform {
            return v
        }
    }

    return ""
}

func RegRegDesc(regType int32) string {
    all := RegTypes()
    for k, v := range all {
        if k == regType {
            return v
        }
    }

    return ""
}