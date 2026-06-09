package router

import "github.com/gin-gonic/gin"

type callback struct{}

func (r *callback) Init(public *gin.RouterGroup, private *gin.RouterGroup) {
	_ = private

	public.POST("papay/notify", apiInfo.Callback.ReceiveContract)
	public.GET("papay/callback-records", apiInfo.Callback.GetCallbackRecordList)
	public.GET("papay/callback-record", apiInfo.Callback.FindCallbackRecord)
}
