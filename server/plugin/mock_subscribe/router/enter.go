package router

import "github.com/flipped-aurora/gin-vue-admin/server/plugin/mock_subscribe/api"

var (
	Router  = new(router)
	apiInfo = api.Api
)

type router struct {
	Merchant       merchant
	Contract       contract
	Deduct         deduct
	Wechat         wechat
	Callback       callback
	DeductCallback deductCallback
}
