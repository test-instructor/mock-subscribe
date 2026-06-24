package model

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// Environment 环境配置
type Environment struct {
	global.GVA_MODEL

	Name   string `json:"name" form:"name" gorm:"column:name;comment:中文名称"`
	Key    string `json:"key" form:"key" gorm:"column:key;comment:唯一标识"`
	Domain string `json:"domain" form:"domain" gorm:"column:domain;comment:服务器地址(含http协议)"`
	Port   int    `json:"port" form:"port" gorm:"column:port;comment:端口号"`
	Remark string `json:"remark" form:"remark" gorm:"column:remark;comment:备注"`
}

func (Environment) TableName() string {
	return "gva_tools_environment"
}
