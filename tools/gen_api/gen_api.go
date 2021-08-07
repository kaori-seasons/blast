package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var (
	method string
)

func init() {
	flag.StringVar(&method, "m", "", "method")
}

func main() {
	flag.Parse()

	if method == "" {
		flag.PrintDefaults()
		return
	}

	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	packName := filepath.Base(dir)

	words := strings.Split(method, "_")
	for i, word := range words {
		words[i] = strings.Title(word)
	}
	methodUpper := strings.Join(words, "")

	methodFile := method + ".go"
	err = ioutil.WriteFile(methodFile,
		[]byte(fmt.Sprintf(handleCont,
			packName, methodUpper, methodUpper, methodUpper)), 0666)
	if err != nil {
		panic(err)
	}
	fmt.Println(methodFile + " is generated")

	methodMsgFile := method + "_msg.go"
	err = ioutil.WriteFile(methodMsgFile,
		[]byte(fmt.Sprintf(msgCont,
			packName, methodUpper, methodUpper, methodUpper,
			methodUpper, methodUpper, methodUpper,
			methodUpper, methodUpper, methodUpper)), 0666)
	if err != nil {
		panic(err)
	}
	fmt.Println(methodMsgFile + " is generated")

	cfgFile := packName + "_cfg.go"
	fi, err := os.Stat(cfgFile)
	if fi != nil && err == nil {
		fmt.Println(cfgFile + " is exist. skip")
	} else {
		if os.IsNotExist(err) {
			err = ioutil.WriteFile(cfgFile,
				[]byte(fmt.Sprintf(cfgCont,
					packName, packName, packName, packName, packName)), 0666)
			if err != nil {
				panic(err)
			}

			fmt.Println(cfgFile + " is generated")
		}
	}

}

var msgCont = `package %s

import "github.com/complone/blast/http_api"

type %sRequest struct {
	// todo add member here
}

func New%sResponse() *%sResponse {
	return &%sResponse{Response: http_api.NewResponse()}
}

func NewErr%sResponse(code http_api.ErrorCode, err error) *%sResponse {
	resp := New%sResponse()
	resp.Code = code
	resp.Msg = code.String(err)
	return resp
}

func (s *%sResponse) GetResponse() http_api.Response {
	return s.Response
}

type %sResponse struct {
	http_api.Response

	// todo add member here
}
`

var handleCont = `package %s

import (
	"github.com/complone/blast/common"
	"github.com/complone/blast/http_api"
	"github.com/gin-gonic/gin"
)

func %sHandle(body []byte, c *gin.Context, resp *http_api.IResponse, ctx *http_api.APIContext) int {
	var req %sRequest
	e := common.DecodeJson(body, &req)
	if e != nil {
		*resp = NewErr%sResponse(http_api.INVALID_PARAM, e)
		return 0
	}
	
	// todo add code here
	return 0
}
`

var cfgCont = `
package %s

import (
	"github.com/complone/blast/common"
)

type %sCfg struct {
	common.Config

	// todo add member here
}

func (s *%sCfg) GetConfig() common.Config {
	return s.Config
}

var My%sCfg %sCfg

func Init() error {
	// todo add member here
	return nil
}
`
