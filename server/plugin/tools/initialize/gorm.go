package initialize

import (
	"context"
	"fmt"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	toolsModel "github.com/flipped-aurora/gin-vue-admin/server/plugin/tools/model"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func Gorm(ctx context.Context) {
	err := global.GVA_DB.WithContext(ctx).AutoMigrate(
		new(toolsModel.Environment),
		new(toolsModel.UserRelation),
		new(toolsModel.SendChatTask),
		new(toolsModel.FanFollowRecord),
	)
	if err != nil {
		err = errors.Wrap(err, "注册tools表失败")
		zap.L().Error(fmt.Sprintf("%+v", err))
	}
}
