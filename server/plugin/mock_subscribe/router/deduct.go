package router

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type deduct struct{}

func (r *deduct) Init(public *gin.RouterGroup, private *gin.RouterGroup) {
	_ = public
	group := private.Group("mockSubscribeDeduct").Use(middleware.OperationRecordWithUserID(3))
	group.GET("getDeductRecordList", apiInfo.Deduct.GetDeductRecordList)
	group.GET("findDeductRecord", apiInfo.Deduct.FindDeductRecord)
}
