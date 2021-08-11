# blast
A network communication framework

> Architecture


- common
    - assert
    - concurrent_limiter 
    - scheduler
    - set
    - domain

- components 
    - redis
        - Exception automatic retry
        - invoke monitor


>Getting started

```sh
$ cat go.mod
module echo
go 1.14
require (
        github.com/complone/blast v0.0.0-20210711135543-1a32326de9d6
        github.com/gin-gonic/gin v1.6.3
        github.com/sirupsen/logrus v1.7.0
)
$ vi go.mod
module echo
go 1.14
require (
-       github.com/complone/blast v0.0.0-20210711135543-1a32326de9d6
        github.com/gin-gonic/gin v1.6.3
        github.com/sirupsen/logrus v1.7.0
)
$ sh build.sh
go: finding module for package github.com/complone/blast/http_api
go: finding module for package github.com/complone/blast
go: finding module for package github.com/complone/blast/common
go: downloading github.com/complone/blast v0.0.0-20210712032730-
ba212fa0c2f9
go: found github.com/complone/blast in github.com/complone/blast v0.0.0-20210712032730-ba212fa0c2f9
$ cat go.mod
module echo
go 1.14
require (
        github.com/complone/blast v0.0.0-20210712032730-ba212fa0c2f9
        github.com/gin-gonic/gin v1.6.3
        github.com/sirupsen/logrus v1.7.0
)
```

> Example

```go
package api
import (
        "github.com/complone/blast/common"
        "github.com/complone/blast/http_api"
        "github.com/gin-gonic/gin"
)
func EchoHandle(body []byte, c *gin.Context, resp *module.IResponse, ctx
*module.APIContext) int {
        var req EchoRequest
        e := common.DecodeJson(body, &req)
        if e != nil {
                *resp = NewErrEchoResponse(module.INVALID_PARAM, e)
                return 0
}}


// todo add code here
res := NewEchoResponse()
if req.Msg != "" {
    res.Msg = req.Msg
} else {
res.Msg = MyapiCfg.Option.DefValue
       ctx.ToLog("reply msg: %s", res.Msg)
       *resp = res
        return 0
}
```
