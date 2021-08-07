package http_api

import (
	"encoding/base64"
	"fmt"
	"github.com/complone/blast/monitor"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Response struct {
	Code    ErrorCode `json:"code"`
	Msg     string    `json:"msg"`
	Version int64     `json:"version"`
}

func NewResponse() Response {
	return Response{
		Code:    OK,
		Msg:     OK.String(nil),
		Version: 1,
	}
}

type IResponse interface {
	GetResponse() Response
}

type APIContext struct {
	Name     string
	RecvTime time.Time
	TraceID  string
}

func NewAPIContext(name, logtag string, reqSize int) *APIContext {
	monitor.ProcessingRequests.WithLabelValues(monitor.GetClutser(), name).Inc()
	monitor.RealTimeRequestBodySize.WithLabelValues(monitor.GetClutser(), name).Set(float64(reqSize))
	ctx := &APIContext{name, time.Now(), logtag}
	ctx.ToLog("NewAPIContext: logtag[%s] reqSize[%d]", logtag, reqSize)
	return ctx
}

func (ctx *APIContext) ToLog(format string, v ...interface{}) {
	log.Infof("%s TraceID[%s] CostSince[%v]",
		fmt.Sprintf(format, v...), ctx.TraceID, time.Since(ctx.RecvTime))
}

func (ctx *APIContext) EndRequest(code ErrorCode, format string, v ...interface{}) {
	ctx.ToLog("EndRequest "+format, v...)

	monitor.ProcessingRequests.WithLabelValues(monitor.GetClutser(), ctx.Name).Dec()
	monitor.TotalRequests.WithLabelValues(monitor.GetClutser(), ctx.Name, fmt.Sprintf("%d", code)).Inc()
	monitor.RequestLatency.WithLabelValues(monitor.GetClutser(), ctx.Name).Observe(float64(time.Since(ctx.RecvTime).Milliseconds()))
	monitor.RealTimeRequestLatency.WithLabelValues(monitor.GetClutser(), ctx.Name).Set(float64(time.Since(ctx.RecvTime).Milliseconds()))
}

// HandleFunc
// return 0: handle process data only and the Frame write the response to client
// return 1: handle process data and write it to client it self
type HandleFunc func(body []byte, c *gin.Context, resp *IResponse, ctx *APIContext) int

// InitFunc
// return error (not nil), program will exit with printing the log
type InitFunc func() error

// return user, passwd
func ParseBasicAuth(c *gin.Context) (string, string) {
	authValue := c.Request.Header.Get("Authorization")
	if authValue == "" {
		realm := "Basic realm=" + strconv.Quote("Authorization Required")
		c.Header("WWW-Authenticate", realm)
		return "", ""
	}

	//"Basic cm9vdDphZG1pbg=="
	auth, err := base64.StdEncoding.DecodeString(authValue[len("Basic "):])
	if err != nil {
		return "", ""
	}

	arr := strings.Split(string(auth), ":")
	if len(arr) != 2 {
		c.AbortWithStatus(http.StatusUnauthorized)
		return "", ""
	}
	return arr[0], arr[1]
}
