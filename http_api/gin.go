package http_api

import (
	"crypto/md5"
	"fmt"
	"github.com/complone/blast/common"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var globalGinEngine *gin.Engine
var ginHandles map[string]HandleFunc

func RunGin(addr string) error {
	gin.SetMode(gin.ReleaseMode)
	globalGinEngine = gin.Default()

	err := initHandle()
	if err != nil {
		return err
	}

	// default handle
	indexHtml := fmt.Sprintf(`
<h2>WelCome to %s</h2>
<table border="1">
  <tr>
    <th>InnerPath</th>
  </tr>
  <tr>
    <td><a href="/openapi">directory of %s</a></td>
  </tr>
  <tr>
	<td><a href="/metrics">metrics of %s</a></td>
  </tr>
</table>
`, common.AppName, common.AppName, common.AppName)
	globalGinEngine.GET("/", func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")

		c.String(http.StatusOK, indexHtml)
	})

	// monitor
	globalGinEngine.GET("metrics", gin.WrapH(promhttp.Handler()))

	path, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return err
	}

	dir := http.Dir(path + "/..")
	log.Infof("set static_fs to %s", dir)
	globalGinEngine.StaticFS("openapi", dir)

	// scheduler
	common.StartSchedule()

	return globalGinEngine.Run(addr)
}

func initHandle() error {
	ginHandles = make(map[string]HandleFunc)

	apis := strings.Split(common.GlobalConfig.OpenAPI.APISwitch, ",")
	for _, api := range apis {
		if handles, ok := ModuleHandleContainer[api]; ok {
			if err := handles.Init(); err != nil {
				return err
			}
			for method, handle := range handles.Handles {
				globalGinEngine.POST(method, globalGinHandle)

				if _, ok := ginHandles[method]; ok {
					return fmt.Errorf("duplicated handle for api.method: %s.%s", api, method)
				}
				ginHandles[method] = handle
			}
		} else {
			return fmt.Errorf("not found Init handle for api: %s", api)
		}
	}
	return nil
}

func globalGinHandle(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		resp := NewResponse()
		resp.Code = INVALID_PARAM
		resp.Msg = resp.Code.String(nil)
		c.JSON(http.StatusOK, resp)
		return
	}

	method := c.Request.URL.Path

	// /v1/query => v1_query
	schemaPath := strings.ReplaceAll(method[1:], "/", "_")
	schema := GetSchema(schemaPath)
	if schema != nil {
		log.Infof("schema(%s) for %s found. run validateJson", schemaPath, method)
		err = ValidateJson(schema, body)
		if err != nil {
			resp := NewResponse()
			resp.Code = SCHEMA_VALIDATE_FAILED
			resp.Msg = resp.Code.String(err)
			c.JSON(http.StatusOK, resp)
			return
		}
	} else {
		log.Infof("schema(%s) for %s not found. skip validateJson", schemaPath, method)
	}

	// schemaCheck
	handle, ok := ginHandles[method]
	common.True(ok)

	logtag := fmt.Sprintf("%s_%x", method, md5.Sum(body))

	var resp IResponse
	ctx := NewAPIContext(method, logtag, len(body))
	defer func() {
		ctx.EndRequest(resp.GetResponse().Code, "%s", resp.GetResponse().Msg)
	}()

	ret := handle(body, c, &resp, ctx)
	if ret == 0 {
		c.JSON(http.StatusOK, resp)
	}
}
