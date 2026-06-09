package router

import "github.com/gin-gonic/gin"

type deductCallback struct{}

func (r *deductCallback) Init(public *gin.RouterGroup, private *gin.RouterGroup) {
	public.POST("pay/pappaynotify", apiInfo.DeductCallback.ReceiveDeduct)

	group := private.Group("mockSubscribeDeductCallback")
	group.GET("getDeductCallbackRecordList", apiInfo.DeductCallback.GetDeductCallbackRecordList)
	group.GET("findDeductCallbackRecord", apiInfo.DeductCallback.FindDeductCallbackRecord)
}
