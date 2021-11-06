package resp

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"pkg.deepin.com/service/lib/response"
)

// Ok ...
func Ok(ctx *gin.Context, data interface{}) {
	ctx.AbortWithStatusJSON(http.StatusOK, response.Response{
		Result: true,
		Code:   http.StatusOK,
		Data:   data,
	})
}

// Fail ...
func Fail(ctx *gin.Context, code int, err error, data ...interface{}) {
	if code == http.StatusUnauthorized {
		ctx.AbortWithStatus(code)
		return
	}
	resp := &response.Response{Code: code, Message: err.Error()}
	if len(data) > 0 {
		resp.Err = data[0]
	}
	ctx.AbortWithStatusJSON(http.StatusOK, resp)
}
