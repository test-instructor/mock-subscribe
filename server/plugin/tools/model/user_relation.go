package model

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// UserRelation 用户-环境关联
type UserRelation struct {
	global.GVA_MODEL

	EnvironmentKey string `json:"environmentKey" form:"environmentKey" gorm:"column:environment_key;comment:关联环境key;index:idx_env_user"`
	UserID         uint64 `json:"userId" form:"userId" gorm:"column:user_id;comment:用户ID;index:idx_env_user"`
}

func (UserRelation) TableName() string {
	return "gva_tools_user_relation"
}
