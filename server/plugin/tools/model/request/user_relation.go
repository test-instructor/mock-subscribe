package model

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

type UserRelationCreate struct {
	EnvironmentKey string `json:"environmentKey" form:"environmentKey"`
	UserIds        string `json:"userIds" form:"userIds"`
}

type UserRelationSearch struct {
	request.PageInfo
	EnvironmentKey string   `form:"environmentKey"`
	UserID         uint64   `form:"userId"`
	UserIds        []uint64 `form:"userIds[]"`
}
