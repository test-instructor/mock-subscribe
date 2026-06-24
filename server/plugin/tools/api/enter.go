package api

import "github.com/flipped-aurora/gin-vue-admin/server/plugin/tools/service"

var (
	Api         = new(api)
	serviceInfo = service.Service
)

type api struct {
	Environment  environment
	UserRelation userRelation
	FanFollow    fanFollow
	SendChat     sendChat
}
