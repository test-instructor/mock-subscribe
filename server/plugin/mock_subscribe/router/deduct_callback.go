package router

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type deductCallback struct{}

func (r *deductCallback) Init(public *gin.RouterGroup, private *gin.RouterGroup) {
	public.POST("pay/pappaynotify", apiInfo.DeductCallback.ReceiveDeduct)

	group := private.Group("mockSubscribeDeductCallback").Use(middleware.OperationRecordWithUserID(3))
	group.GET("getDeductCallbackRecordList", apiInfo.DeductCallback.GetDeductCallbackRecordList)
	group.GET("findDeductCallbackRecord", apiInfo.DeductCallback.FindDeductCallbackRecord)
}
