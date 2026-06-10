package router

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type contract struct{}

func (r *contract) Init(public *gin.RouterGroup, private *gin.RouterGroup) {
	_ = public
	group := private.Group("mockSubscribeContract").Use(middleware.OperationRecordWithUserID(3))
	group.GET("getContractList", apiInfo.Contract.GetContractList)
	group.GET("findContract", apiInfo.Contract.FindContract)
	group.PUT("updateContractStatus", apiInfo.Contract.UpdateContractStatus)
	group.GET("getContractRecordList", apiInfo.Contract.GetContractRecordList)
}
