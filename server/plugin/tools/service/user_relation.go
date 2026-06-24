package service

import (
	"errors"
	"strconv"
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	commonReq "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	toolsModel "github.com/flipped-aurora/gin-vue-admin/server/plugin/tools/model"
	toolsReq "github.com/flipped-aurora/gin-vue-admin/server/plugin/tools/model/request"
)

type userRelation struct{}

func (s *userRelation) CreateUserRelation(info toolsReq.UserRelationCreate) (int, error) {
	if info.EnvironmentKey == "" {
		return 0, errors.New("环境Key不能为空")
	}
	if info.UserIds == "" {
		return 0, errors.New("用户ID不能为空")
	}

	lines := strings.Split(info.UserIds, "\n")
	created := 0
	for _, line := range lines {
		uidStr := strings.TrimSpace(line)
		if uidStr == "" {
			continue
		}
		uid, err := strconv.ParseUint(uidStr, 10, 64)
		if err != nil {
			continue
		}
		relation := toolsModel.UserRelation{
			EnvironmentKey: info.EnvironmentKey,
			UserID:         uid,
		}
		if err := global.GVA_DB.Create(&relation).Error; err == nil {
			created++
		}
	}
	return created, nil
}

func (s *userRelation) DeleteUserRelation(id uint) error {
	if id == 0 {
		return errors.New("ID不能为空")
	}
	return global.GVA_DB.Delete(&toolsModel.UserRelation{}, id).Error
}

func (s *userRelation) GetUserRelation(id uint) (toolsModel.UserRelation, error) {
	var rel toolsModel.UserRelation
	err := global.GVA_DB.Where("id = ?", id).First(&rel).Error
	return rel, err
}

func (s *userRelation) GetUserRelationList(info toolsReq.UserRelationSearch) ([]toolsModel.UserRelation, int64, error) {
	var list []toolsModel.UserRelation
	var total int64

	db := global.GVA_DB.Model(&toolsModel.UserRelation{})
	if info.EnvironmentKey != "" {
		db = db.Where("environment_key = ?", info.EnvironmentKey)
	}
	if info.UserID != 0 {
		db = db.Where("user_id = ?", info.UserID)
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := db.Scopes((&commonReq.PageInfo{Page: info.Page, PageSize: info.PageSize}).Paginate()).Order("id desc").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (s *userRelation) GetUserIdsByEnvironmentKey(envKey string, limit int) ([]uint64, error) {
	var results []struct {
		UserID uint64 `gorm:"column:user_id"`
	}
	err := global.GVA_DB.Model(&toolsModel.UserRelation{}).
		Where("environment_key = ?", envKey).
		Order("id asc").
		Limit(limit).
		Find(&results).Error
	if err != nil {
		return nil, err
	}
	ids := make([]uint64, len(results))
	for i, r := range results {
		ids[i] = r.UserID
	}
	return ids, nil
}
