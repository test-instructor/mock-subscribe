package initialize

import (
	"context"

	model "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/plugin-tool/utils"
)

func Menu(ctx context.Context) {
	_ = ctx
	entities := []model.SysBaseMenu{
		{
			ParentId:  0,
			Path:      "mockSubscribe",
			Name:      "MockSubscribeRoot",
			Hidden:    false,
			Component: "view/routerHolder.vue",
			Sort:      92,
			Meta:      model.Meta{Title: "Mock订阅服务", Icon: "coin"},
		},
		{
			ParentId:  0,
			Path:      "mockSubscribeMerchant",
			Name:      "MockSubscribeMerchant",
			Hidden:    false,
			Component: "plugin/mock_subscribe/view/merchant.vue",
			Sort:      1,
			Meta:      model.Meta{Title: "商户配置", Icon: "office-building", KeepAlive: true},
		},
		{
			ParentId:  0,
			Path:      "mockSubscribeContract",
			Name:      "MockSubscribeContract",
			Hidden:    false,
			Component: "plugin/mock_subscribe/view/contract.vue",
			Sort:      2,
			Meta:      model.Meta{Title: "用户协议", Icon: "tickets", KeepAlive: true},
		},
		{
			ParentId:  0,
			Path:      "mockSubscribeDeduct",
			Name:      "MockSubscribeDeduct",
			Hidden:    false,
			Component: "plugin/mock_subscribe/view/deduct.vue",
			Sort:      3,
			Meta:      model.Meta{Title: "扣款记录", Icon: "wallet", KeepAlive: true},
		},
		{
			ParentId:  0,
			Path:      "mockSubscribeCallback",
			Name:      "MockSubscribeCallback",
			Hidden:    false,
			Component: "plugin/mock_subscribe/view/callback.vue",
			Sort:      4,
			Meta:      model.Meta{Title: "回调记录", Icon: "notification", KeepAlive: true},
		},
	}
	utils.RegisterMenus(entities...)
}
