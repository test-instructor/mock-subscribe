package router

import (
	"github.com/gin-gonic/gin"
)

type sendChat struct{}

func (r *sendChat) Init(public, private *gin.RouterGroup) {
	group := private.Group("toolsSendChat")
	group.POST("createSendChatTask", apiInfo.SendChat.CreateSendChatTask)
	group.GET("getSendChatTaskList", apiInfo.SendChat.GetSendChatTaskList)
	group.PUT("stopSendChatTask", apiInfo.SendChat.StopSendChatTask)
}
