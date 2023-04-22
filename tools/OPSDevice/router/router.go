package router

import (
	"github.com/gin-gonic/gin"
	"go/gopkg/middleware"
	"go/tools/OPSDevice/controller"
	"net/http"
)

func SetupRouter(engine *gin.Engine) {
	//set up router middleware
	engine.Use(middleware.AccessLog(), middleware.Recovery())

	//404
	engine.NoRoute(func(ctx *gin.Context) {
		middleware.ResponseJSON(ctx, http.StatusNotFound, "请求方法不存在", nil)
	})

	//register
	engine.GET("/opsdevice/ping", controller.PingExecute)

}
