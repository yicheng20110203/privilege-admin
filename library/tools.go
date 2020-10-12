package library

import (
	"crypto/md5"
	"crypto/tls"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"gitlab.ceibsmoment.com/c/mp/logger"
	"strings"
	"time"
)

var (
	JsonpbMarshaler   *jsonpb.Marshaler
	JsonpbUnmarshaler *jsonpb.Unmarshaler
)

// 引入包时自动初始化
func init() {
	JsonpbMarshaler = &jsonpb.Marshaler{
		EnumsAsInts:  true,
		OrigName:     true,
		EmitDefaults: true,
	}

	JsonpbUnmarshaler = &jsonpb.Unmarshaler{
		AllowUnknownFields: true,
	}
}

func MD5(s string) string {
	r := fmt.Sprintf("%x", md5.Sum([]byte(s)))
	r = strings.ToLower(r)
	return r
}

func Marshaler(message proto.Message) (s string, err error) {
	s, err = JsonpbMarshaler.MarshalToString(message)
	if err != nil {
		logger.Logger.Errorf("library.Marshaler <%v> error: <%v>", message, err)
	}
	return
}

func Unmarshaler(s string, message *proto.Message) (err error) {
	r := strings.NewReader(s)
	err = JsonpbUnmarshaler.Unmarshal(r, *message)
	return
}

type RestyRequest struct {
	Url        string                 `json:"url"`
	Method     string                 `json:"method"`
	Params     map[string]interface{} `json:"params"`
	Header     map[string]string      `json:"header"`
	RetryTimes int                    `json:"retry_times"`
	Timeout    time.Duration          `json:"timeout"`
	SkipVerify bool                   `json:"insecure_skip_verify"`
}

func HttpRequest(req *RestyRequest) (resp []byte, err error) {
	logger.Logger.Info("-------------- Library HttpRequest() Begin ------------")
	logger.Logger.Info("Library HttpRequest() request params: ", *req)
	var response *resty.Response

	defer func() {
		if err == nil {
			ti := response.Request.TraceInfo()
			traceInfo := map[string]interface{}{
				"DNSLookup":     ti.DNSLookup,
				"ConnTime":      ti.ConnTime,
				"TLSHandshake":  ti.TLSHandshake,
				"ServerTime":    ti.ServerTime,
				"ResponseTime":  ti.ResponseTime,
				"TotalTime":     ti.TotalTime,
				"IsConnReused":  ti.IsConnReused,
				"IsConnWasIdle": ti.IsConnWasIdle,
				"ConnIdleTime":  ti.ConnIdleTime,
			}
			logger.Logger.Info("Library HttpRequest() trace info : ", traceInfo)
		}

		logger.Logger.Info("Library HttpRequest() response: ", string(resp))
		logger.Logger.Info("-------------- Library HttpRequest() End ------------")
	}()

	c := resty.New()
	if req.RetryTimes > 0 {
		c = c.SetRetryCount(req.RetryTimes).SetRetryWaitTime(req.Timeout)
	} else {
		c = c.SetTimeout(req.Timeout)
	}

	// skip verify
	if req.SkipVerify {
		c.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}

	r := c.R().SetHeaders(req.Header).EnableTrace()
	r = r.SetBody(req.Params)

	// GET request
	if strings.ToUpper(req.Method) == "GET" {
		response, err = r.Get(req.Url)
		if err != nil {
			logger.Logger.Errorf("Library HttpRequest() r.Get error: <%#v>", err)
			return nil, err
		}

		resp = response.Body()
		return resp, err
	}

	// POST request
	response, err = r.Post(req.Url)
	if err != nil {
		logger.Logger.Errorf("Library HttpRequest() r.Post error: %#v", err)
		return nil, err
	}
	resp = response.Body()

	return
}
