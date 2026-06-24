package model

import "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"

type FanFollowCreate struct {
	EnvironmentKey string `json:"environmentKey" form:"environmentKey"`
	UserID         uint64 `json:"userId" form:"userId"`
	Operation      string `json:"operation" form:"operation"`
	Count          int    `json:"count" form:"count"`
}

type FanFollowSearch struct {
	request.PageInfo
	EnvironmentKey string `form:"environmentKey"`
	Operation      string `form:"operation"`
}
