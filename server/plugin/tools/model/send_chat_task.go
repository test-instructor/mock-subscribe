package model

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// SendChatTask 发送公屏消息任务
type SendChatTask struct {
	global.GVA_MODEL

	RoomID             string `json:"roomId" form:"roomId" gorm:"column:room_id;comment:房间ID"`
	EnvironmentKey     string `json:"environmentKey" form:"environmentKey" gorm:"column:environment_key;comment:环境key"`
	AccountCount       int    `json:"accountCount" form:"accountCount" gorm:"column:account_count;comment:发送账号数量"`
	MsgCountPerAccount int    `json:"msgCountPerAccount" form:"msgCountPerAccount" gorm:"column:msg_count_per_account;comment:每账号消息数"`
	MsgInterval        int    `json:"msgInterval" form:"msgInterval" gorm:"column:msg_interval;comment:消息间隔(毫秒)"`
	Status             string `json:"status" form:"status" gorm:"column:status;default:running;comment:任务状态(running/completed/stopped)"`
	SuccessCount       int    `json:"successCount" form:"successCount" gorm:"column:success_count;default:0;comment:成功数量"`
}

func (SendChatTask) TableName() string {
	return "gva_tools_send_chat_task"
}
