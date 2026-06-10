package router

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type merchant struct{}

func (r *merchant) Init(public *gin.RouterGroup, private *gin.RouterGroup) {
	_ = public
	group := private.Group("mockSubscribeMerchant").Use(middleware.OperationRecordWithUserID(3))
	group.POST("createMerchant", apiInfo.Merchant.CreateMerchant)
	group.PUT("updateMerchant", apiInfo.Merchant.UpdateMerchant)
	group.DELETE("deleteMerchant", apiInfo.Merchant.DeleteMerchant)
	group.GET("findMerchant", apiInfo.Merchant.FindMerchant)
	group.GET("getMerchantList", apiInfo.Merchant.GetMerchantList)
}
