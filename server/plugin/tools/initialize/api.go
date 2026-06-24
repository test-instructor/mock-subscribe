package initialize

import (
	"context"

	model "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/plugin-tool/utils"
)

func Api(ctx context.Context) {
	_ = ctx
	entities := []model.SysApi{
		{Path: "/toolsEnvironment/createEnvironment", Description: "创建环境配置", ApiGroup: "Tools-环境", Method: "POST"},
		{Path: "/toolsEnvironment/updateEnvironment", Description: "更新环境配置", ApiGroup: "Tools-环境", Method: "PUT"},
		{Path: "/toolsEnvironment/deleteEnvironment", Description: "删除环境配置", ApiGroup: "Tools-环境", Method: "DELETE"},
		{Path: "/toolsEnvironment/findEnvironment", Description: "获取环境配置详情", ApiGroup: "Tools-环境", Method: "GET"},
		{Path: "/toolsEnvironment/getEnvironmentList", Description: "获取环境配置列表", ApiGroup: "Tools-环境", Method: "GET"},

		{Path: "/toolsUserRelation/createUserRelation", Description: "批量创建用户关联", ApiGroup: "Tools-用户关联", Method: "POST"},
		{Path: "/toolsUserRelation/findUserRelation", Description: "获取用户关联详情", ApiGroup: "Tools-用户关联", Method: "GET"},
		{Path: "/toolsUserRelation/getUserRelationList", Description: "获取用户关联列表", ApiGroup: "Tools-用户关联", Method: "GET"},
		{Path: "/toolsUserRelation/deleteUserRelation", Description: "删除用户关联", ApiGroup: "Tools-用户关联", Method: "DELETE"},
		{Path: "/toolsUserRelation/getUserIdsByEnvironment", Description: "根据环境获取用户ID列表", ApiGroup: "Tools-用户关联", Method: "GET"},

		{Path: "/toolsFanFollow/createFanFollow", Description: "执行粉丝/关注/好友操作", ApiGroup: "Tools-粉丝关注好友", Method: "POST"},
		{Path: "/toolsFanFollow/getFanFollowList", Description: "获取粉丝/关注/好友操作记录", ApiGroup: "Tools-粉丝关注好友", Method: "GET"},

		{Path: "/toolsSendChat/createSendChatTask", Description: "创建并启动发送任务", ApiGroup: "Tools-发送公屏消息", Method: "POST"},
		{Path: "/toolsSendChat/getSendChatTaskList", Description: "获取发送任务列表", ApiGroup: "Tools-发送公屏消息", Method: "GET"},
		{Path: "/toolsSendChat/stopSendChatTask", Description: "停止发送任务", ApiGroup: "Tools-发送公屏消息", Method: "PUT"},
	}
	utils.RegisterApis(entities...)
}
