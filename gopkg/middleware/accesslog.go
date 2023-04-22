/*
** accesss log
 */

package middleware

import (
	"github.com/gin-gonic/gin"
	"go/gopkg/logger/vglog"
	"time"
)

//access log
func accessLog(costTime time.Duration, ctx *gin.Context) {
	//clientIp := GetRemoteIp(ctx.Request)
	diffTm := int64(costTime) / 1000

	now := time.Now()

	vglog.Access(`{"date":"%04d-%02d-%02d %02d:%02d:%02d", "host":"%s", "method":"%s","uri":"%s", "proto":"%s","code":%d, "handleTime":%vμs, "path":"%s", "referer": "%s"}`,
		now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(),
		ctx.ClientIP(), ctx.Request.Method, ctx.Request.RequestURI, ctx.Request.Proto, ctx.Writer.Status(),
		diffTm, ctx.Request.URL.Path, ctx.Request.Referer())

	AddInterfaceAccessDuration(ctx.Request.URL.Path, diffTm)
}

func AccessLog() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		startTime := time.Now()

		// 调用下一个中间件
		ctx.Next()

		costTime := time.Since(startTime)

		accessLog(costTime, ctx)
	}
}
