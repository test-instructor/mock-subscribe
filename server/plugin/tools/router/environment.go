package router

import (
	"github.com/gin-gonic/gin"
)

type environment struct{}

func (r *environment) Init(public, private *gin.RouterGroup) {
	group := private.Group("toolsEnvironment")
	group.POST("createEnvironment", apiInfo.Environment.CreateEnvironment)
	group.PUT("updateEnvironment", apiInfo.Environment.UpdateEnvironment)
	group.DELETE("deleteEnvironment", apiInfo.Environment.DeleteEnvironment)
	group.GET("findEnvironment", apiInfo.Environment.FindEnvironment)
	group.GET("getEnvironmentList", apiInfo.Environment.GetEnvironmentList)
}
