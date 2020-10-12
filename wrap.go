package mp

import (
    "github.com/gin-gonic/gin"
    "github.com/golang/protobuf/proto"
    "github.com/golang/protobuf/ptypes"
    "github.com/golang/protobuf/ptypes/any"
    def "gitlab.ceibsmoment.com/c/mp/global"
    "gitlab.ceibsmoment.com/c/mp/library"
    "gitlab.ceibsmoment.com/c/mp/logger"
    pbs "gitlab.ceibsmoment.com/c/mp/mppbs"
)

// WrapOutput response wrapper
// wrap response data for client
func WrapOutput(ctx *gin.Context, code def.ErrorCode, msg string, data proto.Message) (err error) {
    var pbAny *any.Any
    pbAny, err = ptypes.MarshalAny(data)
    if err != nil {
        logger.Logger.Errorf("wrap.WrapOutput ptypes.MarshalAny(%v) error: (%v)", data, err)
        code = def.ProtoAnyErr
        msg = _errorMsgWithDefault(def.ProtoAnyErr, "server error")
        data = nil
    }

    message := &pbs.ApiResponse{
        Code: int32(code),
        Msg:  msg,
        Data: pbAny,
    }
    s, _ := library.Marshaler(message)
    _, err = ctx.Writer.Write([]byte(s))
    if err != nil {
        logger.Logger.Errorf("wrap.WrapOutput ctx.Writer.Write error(%v)", err)
    }
    return
}

// 优先使用定义的code对应msg
func JsonData(code def.ErrorCode, msg string, data interface{}) gin.H {
    return gin.H{
        "code": int32(code),
        "msg":  _errorMsgWithDefault(code, msg),
        "data": data,
    }
}

func ResetJsonData(code def.ErrorCode, msg string, data interface{}) gin.H {
    return gin.H{
        "code": int32(code),
        "msg":  msg,
        "data": data,
    }
}

func _errorMsgWithDefault(code def.ErrorCode, defaultVal string) (msg string) {
    errors := def.Errors()
    if v, ok := errors[code]; ok {
        msg = v
        return msg
    }

    return defaultVal
}
