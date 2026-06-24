package router

import (
	"github.com/gin-gonic/gin"
)

type fanFollow struct{}

func (r *fanFollow) Init(public, private *gin.RouterGroup) {
	group := private.Group("toolsFanFollow")
	group.POST("createFanFollow", apiInfo.FanFollow.CreateFanFollow)
	group.GET("getFanFollowList", apiInfo.FanFollow.GetFanFollowList)
}
