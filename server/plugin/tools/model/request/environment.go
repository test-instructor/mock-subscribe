package model

import "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"

type EnvironmentSearch struct {
	request.PageInfo
	Name string `form:"name"`
	Key  string `form:"key"`
}
