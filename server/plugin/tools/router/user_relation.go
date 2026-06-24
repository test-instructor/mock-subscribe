package router

import (
	"github.com/gin-gonic/gin"
)

type userRelation struct{}

func (r *userRelation) Init(public, private *gin.RouterGroup) {
	group := private.Group("toolsUserRelation")
	group.POST("createUserRelation", apiInfo.UserRelation.CreateUserRelation)
	group.GET("findUserRelation", apiInfo.UserRelation.FindUserRelation)
	group.GET("getUserRelationList", apiInfo.UserRelation.GetUserRelationList)
	group.DELETE("deleteUserRelation", apiInfo.UserRelation.DeleteUserRelation)
	group.GET("getUserIdsByEnvironment", apiInfo.UserRelation.GetUserIdsByEnvironment)
}
