package model

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// FanFollowRecord 粉丝/关注/好友操作执行记录
type FanFollowRecord struct {
	global.GVA_MODEL

	EnvironmentKey string `json:"environmentKey" form:"environmentKey" gorm:"column:environment_key;comment:环境key;index:idx_env_created"`
	UserID         uint64 `json:"userId" form:"userId" gorm:"column:user_id;comment:被操作目标用户ID"`
	Operation      string `json:"operation" form:"operation" gorm:"column:operation;comment:操作类型 fans/follow/friend"`
	Count          int    `json:"count" form:"count" gorm:"column:count;comment:计划执行数量"`
	SuccessCount   int    `json:"successCount" form:"successCount" gorm:"column:success_count;default:0;comment:实际成功数量"`
	Status         string `json:"status" form:"status" gorm:"column:status;default:running;comment:状态 running/completed/failed"`
}

func (FanFollowRecord) TableName() string {
	return "gva_tools_fan_follow_record"
}
