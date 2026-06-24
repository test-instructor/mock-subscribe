package service

import (
	"errors"
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	commonReq "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	toolsModel "github.com/flipped-aurora/gin-vue-admin/server/plugin/tools/model"
	toolsReq "github.com/flipped-aurora/gin-vue-admin/server/plugin/tools/model/request"
)

type environment struct{}

func (s *environment) CreateEnvironment(info *toolsModel.Environment) error {
	if strings.TrimSpace(info.Name) == "" {
		return errors.New("名称不能为空")
	}
	if strings.TrimSpace(info.Key) == "" {
		return errors.New("Key不能为空")
	}
	var count int64
	global.GVA_DB.Model(&toolsModel.Environment{}).Where("key = ? AND id <> ?", info.Key, info.ID).Count(&count)
	if count > 0 {
		return errors.New("Key已存在")
	}
	return global.GVA_DB.Create(info).Error
}

func (s *environment) UpdateEnvironment(info *toolsModel.Environment) error {
	if info.ID == 0 {
		return errors.New("ID不能为空")
	}
	if strings.TrimSpace(info.Name) == "" {
		return errors.New("名称不能为空")
	}
	if strings.TrimSpace(info.Key) == "" {
		return errors.New("Key不能为空")
	}
	var current toolsModel.Environment
	if err := global.GVA_DB.Where("id = ?", info.ID).First(&current).Error; err != nil {
		return errors.New("环境配置不存在")
	}
	if current.Key != info.Key {
		return errors.New("Key不允许修改")
	}
	return global.GVA_DB.Model(&toolsModel.Environment{}).
		Where("id = ?", info.ID).
		Updates(map[string]interface{}{
			"name":   info.Name,
			"domain": info.Domain,
			"port":   info.Port,
			"remark": info.Remark,
		}).Error
}

func (s *environment) DeleteEnvironment(id uint) error {
	if id == 0 {
		return errors.New("ID不能为空")
	}
	return global.GVA_DB.Delete(&toolsModel.Environment{}, id).Error
}

func (s *environment) GetEnvironment(id uint) (toolsModel.Environment, error) {
	var env toolsModel.Environment
	err := global.GVA_DB.Where("id = ?", id).First(&env).Error
	return env, err
}

func (s *environment) GetEnvironmentByKey(key string) (toolsModel.Environment, error) {
	var env toolsModel.Environment
	err := global.GVA_DB.Where("key = ?", key).First(&env).Error
	return env, err
}

func (s *environment) GetEnvironmentList(info toolsReq.EnvironmentSearch) ([]toolsModel.Environment, int64, error) {
	var list []toolsModel.Environment
	var total int64

	db := global.GVA_DB.Model(&toolsModel.Environment{})
	if info.Name != "" {
		db = db.Where("name LIKE ?", "%"+info.Name+"%")
	}
	if info.Key != "" {
		db = db.Where("key LIKE ?", "%"+info.Key+"%")
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := db.Scopes((&commonReq.PageInfo{Page: info.Page, PageSize: info.PageSize}).Paginate()).Order("id desc").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}
