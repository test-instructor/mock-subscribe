package initialize

import (
	"context"
	"fmt"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/mock_subscribe/model"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func Gorm(ctx context.Context) {
	err := global.GVA_DB.WithContext(ctx).AutoMigrate(
		new(model.Merchant),
		new(model.Contract),
		new(model.ContractStatusRecord),
		new(model.ContractRecord),
		new(model.DeductRecord),
		new(model.CallbackRecord),
	)
	if err != nil {
		err = errors.Wrap(err, "注册微信订阅mock表失败")
		zap.L().Error(fmt.Sprintf("%+v", err))
	}
}
