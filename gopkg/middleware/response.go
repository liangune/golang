package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type response struct {
	Code    int         `json:"code"`
	Message string      `json:"msg"`
	Data    interface{} `json:"data"`
}

func ResponseJSON(ctx *gin.Context, code int, msg string, data interface{}) {
	ctx.JSON(http.StatusOK, response{
		Code:    code,
		Message: msg,
		Data:    data,
	})
	return
}
