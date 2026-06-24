package router

import "github.com/flipped-aurora/gin-vue-admin/server/plugin/tools/api"

var (
	Router  = new(router)
	apiInfo = api.Api
)

type router struct {
	Environment  environment
	UserRelation userRelation
	FanFollow    fanFollow
	SendChat     sendChat
}
