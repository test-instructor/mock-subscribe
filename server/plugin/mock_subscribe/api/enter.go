package api

import "github.com/flipped-aurora/gin-vue-admin/server/plugin/mock_subscribe/service"

var (
	Api         = new(api)
	serviceInfo = service.Service
)

type api struct {
	Merchant       merchant
	Contract       contract
	Deduct         deduct
	Wechat         wechat
	Callback       callback
	DeductCallback deductCallback
}
