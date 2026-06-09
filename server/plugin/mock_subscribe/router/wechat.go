package router

import "github.com/gin-gonic/gin"

type wechat struct{}

func (r *wechat) Init(public *gin.RouterGroup, private *gin.RouterGroup) {
	_ = private

	public.POST("papay/preentrustweb", apiInfo.Wechat.ContractSign)
	public.POST("papay/querycontract", apiInfo.Wechat.QueryContract)
	public.POST("papay/deletecontract", apiInfo.Wechat.TerminateContract)
	public.POST("pay/pappayapply", apiInfo.Wechat.ApplyDeduct)
	public.POST("transit/queryorder", apiInfo.Wechat.QueryDeduct)
	public.POST("v3/papay/contracts/:contract_id/notify", apiInfo.Wechat.PreDeductNotify)
}
