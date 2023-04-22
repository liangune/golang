package middleware

import (
	"github.com/gin-gonic/gin"
	"go/gopkg/logger/vglog"
	"net/http"
	"runtime/debug"
)

func Recovery() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				vglog.Error("debug stack warn: %v", string(debug.Stack()))
				ResponseJSON(ctx, http.StatusBadRequest, "server panic", nil)
				return
			}
		}()
		ctx.Next()
	}
}
