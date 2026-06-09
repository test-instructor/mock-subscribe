package router

import (
	"github.com/gin-gonic/gin"
)

type contract struct{}

func (r *contract) Init(public *gin.RouterGroup, private *gin.RouterGroup) {
	_ = public
	group := private.Group("mockSubscribeContract")
	group.GET("getContractList", apiInfo.Contract.GetContractList)
	group.GET("findContract", apiInfo.Contract.FindContract)
	group.PUT("updateContractStatus", apiInfo.Contract.UpdateContractStatus)
	group.GET("getContractRecordList", apiInfo.Contract.GetContractRecordList)
}
