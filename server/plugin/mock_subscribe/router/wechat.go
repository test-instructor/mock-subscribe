package router

import "github.com/gin-gonic/gin"

type wechat struct{}

func (r *wechat) Init(public *gin.RouterGroup, private *gin.RouterGroup) {
	_ = private
	group := public.Group("mockSubscribeWechat")
	group.POST("contractSign", apiInfo.Wechat.ContractSign)
	group.POST("queryContract", apiInfo.Wechat.QueryContract)
	group.POST("terminateContract", apiInfo.Wechat.TerminateContract)
	group.POST("applyDeduct", apiInfo.Wechat.ApplyDeduct)
	group.POST("queryDeduct", apiInfo.Wechat.QueryDeduct)
	group.POST("preDeductNotify", apiInfo.Wechat.PreDeductNotify)

	contractPublicGroup := public.Group("mockSubscribeContract")
	contractPublicGroup.POST("appPureSign", apiInfo.Wechat.ContractSign)
	contractPublicGroup.POST("queryContract", apiInfo.Wechat.QueryContract)
	contractPublicGroup.POST("terminateContract", apiInfo.Wechat.TerminateContract)

	deductPublicGroup := public.Group("mockSubscribeDeduct")
	deductPublicGroup.POST("applyDeduct", apiInfo.Wechat.ApplyDeduct)
	deductPublicGroup.POST("queryDeduct", apiInfo.Wechat.QueryDeduct)
	deductPublicGroup.POST("preDeductNotify", apiInfo.Wechat.PreDeductNotify)
}
