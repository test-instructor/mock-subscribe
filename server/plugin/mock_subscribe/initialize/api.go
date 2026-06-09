package initialize

import (
	"context"

	model "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/plugin-tool/utils"
)

func Api(ctx context.Context) {
	_ = ctx
	entities := []model.SysApi{
		{Path: "/mockSubscribeMerchant/createMerchant", Description: "创建商户配置", ApiGroup: "微信订阅Mock-商户", Method: "POST"},
		{Path: "/mockSubscribeMerchant/updateMerchant", Description: "更新商户配置", ApiGroup: "微信订阅Mock-商户", Method: "PUT"},
		{Path: "/mockSubscribeMerchant/deleteMerchant", Description: "删除商户配置", ApiGroup: "微信订阅Mock-商户", Method: "DELETE"},
		{Path: "/mockSubscribeMerchant/findMerchant", Description: "获取商户配置详情", ApiGroup: "微信订阅Mock-商户", Method: "GET"},
		{Path: "/mockSubscribeMerchant/getMerchantList", Description: "获取商户配置列表", ApiGroup: "微信订阅Mock-商户", Method: "GET"},

		{Path: "/mockSubscribeContract/getContractList", Description: "获取用户协议列表", ApiGroup: "微信订阅Mock-协议", Method: "GET"},
		{Path: "/mockSubscribeContract/findContract", Description: "获取用户协议详情", ApiGroup: "微信订阅Mock-协议", Method: "GET"},
		{Path: "/mockSubscribeContract/updateContractStatus", Description: "更新用户协议状态", ApiGroup: "微信订阅Mock-协议", Method: "PUT"},
		{Path: "/mockSubscribeContract/getContractRecordList", Description: "获取协议流水列表", ApiGroup: "微信订阅Mock-协议", Method: "GET"},

		{Path: "/mockSubscribeDeduct/getDeductRecordList", Description: "获取扣款记录列表", ApiGroup: "微信订阅Mock-扣款", Method: "GET"},
		{Path: "/mockSubscribeDeduct/findDeductRecord", Description: "获取扣款记录详情", ApiGroup: "微信订阅Mock-扣款", Method: "GET"},

		{Path: "/mockSubscribeWechat/contractSign", Description: "APP纯签约", ApiGroup: "微信订阅Mock-微信接口", Method: "POST"},
		{Path: "/mockSubscribeWechat/queryContract", Description: "查询签约关系", ApiGroup: "微信订阅Mock-微信接口", Method: "POST"},
		{Path: "/mockSubscribeWechat/terminateContract", Description: "申请解约", ApiGroup: "微信订阅Mock-微信接口", Method: "POST"},
		{Path: "/mockSubscribeWechat/applyDeduct", Description: "申请扣款", ApiGroup: "微信订阅Mock-微信接口", Method: "POST"},
		{Path: "/mockSubscribeWechat/queryDeduct", Description: "查询扣款结果", ApiGroup: "微信订阅Mock-微信接口", Method: "POST"},
		{Path: "/mockSubscribeWechat/preDeductNotify", Description: "预扣费通知API", ApiGroup: "微信订阅Mock-微信接口", Method: "POST"},
	}
	utils.RegisterApis(entities...)
}
