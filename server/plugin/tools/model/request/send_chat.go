package model

import "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"

type SendChatCreate struct {
	RoomID             string `json:"roomId" form:"roomId"`
	EnvironmentKey     string `json:"environmentKey" form:"environmentKey"`
	AccountCount       int    `json:"accountCount" form:"accountCount"`
	MsgCountPerAccount int    `json:"msgCountPerAccount" form:"msgCountPerAccount"`
	MsgInterval        int    `json:"msgInterval" form:"msgInterval"`
}

type SendChatSearch struct {
	request.PageInfo
}
