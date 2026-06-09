package router

import "github.com/gin-gonic/gin"

type callback struct{}

func (r *callback) Init(public *gin.RouterGroup, private *gin.RouterGroup) {
	publicGroup := public.Group("mockSubscribeCallback")
	publicGroup.POST("receiveContract", apiInfo.Callback.ReceiveContract)

	privateGroup := private.Group("mockSubscribeCallback")
	privateGroup.GET("getCallbackRecordList", apiInfo.Callback.GetCallbackRecordList)
	privateGroup.GET("findCallbackRecord", apiInfo.Callback.FindCallbackRecord)
}
