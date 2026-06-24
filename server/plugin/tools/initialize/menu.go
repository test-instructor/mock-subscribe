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
			Path:      "toolsEnvironment",
			Name:      "ToolsEnvironment",
			Hidden:    false,
			Component: "plugin/tools/view/environment.vue",
			Sort:      1,
			Meta:      model.Meta{Title: "环境配置", Icon: "setting", KeepAlive: true},
		},
		{
			ParentId:  0,
			Path:      "toolsUserRelation",
			Name:      "ToolsUserRelation",
			Hidden:    false,
			Component: "plugin/tools/view/userRelation.vue",
			Sort:      2,
			Meta:      model.Meta{Title: "用户数据", Icon: "user", KeepAlive: true},
		},
		{
			ParentId:  0,
			Path:      "toolsFanFollow",
			Name:      "ToolsFanFollow",
			Hidden:    false,
			Component: "plugin/tools/view/fanFollow.vue",
			Sort:      3,
			Meta:      model.Meta{Title: "粉丝/关注/好友", Icon: "star", KeepAlive: true},
		},
		{
			ParentId:  0,
			Path:      "toolsSendChat",
			Name:      "ToolsSendChat",
			Hidden:    false,
			Component: "plugin/tools/view/sendChat.vue",
			Sort:      4,
			Meta:      model.Meta{Title: "发送公屏消息", Icon: "chat", KeepAlive: true},
		},
	}
	utils.RegisterMenus(entities...)
}
