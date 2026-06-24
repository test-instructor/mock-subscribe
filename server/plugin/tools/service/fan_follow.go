package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	commonReq "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	toolsModel "github.com/flipped-aurora/gin-vue-admin/server/plugin/tools/model"
	toolsReq "github.com/flipped-aurora/gin-vue-admin/server/plugin/tools/model/request"
	"go.uber.org/zap"
)

type fanFollow struct{}

type friendOperBody struct {
	Param  map[string]interface{} `json:"param"`
	Global map[string]string      `json:"global"`
}

var friendHeaders = map[string]string{
	"appid":       "1",
	"client":      "ios;15.4.1;iPhone11,8;5.27.0;5.27.0.96",
	"deviceid":    "1BDB215B-4CF5-4540-9BDE-C7B4D8574731",
	"packagename": "com.youzuo.miko.inner",
	"user-agent":  "yyzz/5.27.0 (com.youzuo.miko.inner; build:96; iOS 15.4.1) Alamofire/5.5.0",
	"version":     "5.27.0",
	"uid":         "200426370",
	"application": "myyw",
}

func buildFriendGlobal(uid string) map[string]string {
	return map[string]string{
		"Accept":          "*/*",
		"Accept-Encoding": "br;q=1.0, gzip;q=0.9, deflate;q=0.8",
		"Accept-Language": "zh-Hans-CN;q=1.0",
		"Appid":           "1",
		"Application":     "zanzan",
		"Client":          "ios;15.4.1;iPhone11,8;5.27.0;5.27.0.96",
		"Connection":      "keep-alive",
		"Deviceid":        "1BDB215B-4CF5-4540-9BDE-C7B4D8574731",
		"Host":            "api-test.miyafm.com:443",
		"MarketFlag":      "1",
		"PackageName":     "com.youzuo.miko.inner",
		"Roomid":          "0",
		"Uid":             uid,
		"User-Agent":      "yyzz/5.27.0 (com.youzuo.miko.inner; build:96; iOS 15.4.1) Alamofire/5.5.0",
		"Version":         "5.27.0",
	}
}

func doFriendOper(client *http.Client, baseURL string, targetID int64, oper int, global map[string]string) error {
	body := friendOperBody{
		Param: map[string]interface{}{
			"id":   targetID,
			"oper": oper,
		},
		Global: global,
	}
	jsonData, err := json.Marshal(body)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", baseURL+"/testProxy/FriendExtObj/FriendOper", bytes.NewReader(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	for k, v := range friendHeaders {
		req.Header.Set(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = io.ReadAll(resp.Body)
	return err
}

func (s *fanFollow) CreateFanFollow(req toolsReq.FanFollowCreate) (uint, int, error) {
	if req.EnvironmentKey == "" {
		return 0, 0, errors.New("环境Key不能为空")
	}
	if req.UserID == 0 {
		return 0, 0, errors.New("用户ID不能为空")
	}
	if req.Operation != "fans" && req.Operation != "follow" && req.Operation != "friend" {
		return 0, 0, errors.New("不支持的操作类型")
	}
	if req.Count <= 0 {
		return 0, 0, errors.New("执行次数必须大于0")
	}

	var env toolsModel.Environment
	if err := global.GVA_DB.Where("key = ?", req.EnvironmentKey).First(&env).Error; err != nil {
		return 0, 0, errors.New("未找到对应环境配置")
	}

	allIDs, err := (&userRelation{}).GetUserIdsByEnvironmentKey(req.EnvironmentKey, 0)
	if err != nil {
		return 0, 0, errors.New("获取环境用户数据失败")
	}
	if len(allIDs) == 0 {
		return 0, 0, errors.New("该环境未维护用户数据")
	}
	rand.Shuffle(len(allIDs), func(i, j int) { allIDs[i], allIDs[j] = allIDs[j], allIDs[i] })
	if len(allIDs) > req.Count {
		allIDs = allIDs[:req.Count]
	}

	record := toolsModel.FanFollowRecord{
		EnvironmentKey: req.EnvironmentKey,
		UserID:         req.UserID,
		Operation:      req.Operation,
		Count:          req.Count,
		SuccessCount:   0,
		Status:         "running",
	}
	if err := global.GVA_DB.Create(&record).Error; err != nil {
		return 0, 0, err
	}

	baseURL := fmt.Sprintf("%s:%d", env.Domain, env.Port)
	client := &http.Client{Timeout: 10 * time.Second}
	successCount := 0

	for _, uid := range allIDs {
		uidStr := strconv.FormatUint(uid, 10)
		globalFields := buildFriendGlobal(uidStr)
		switch req.Operation {
		case "fans":
			if err := doFriendOper(client, baseURL, int64(req.UserID), 1, globalFields); err == nil {
				successCount++
			}
		case "follow":
			if err := doFriendOper(client, baseURL, int64(req.UserID), 0, globalFields); err == nil {
				successCount++
			}
		case "friend":
			if err := doFriendOper(client, baseURL, int64(req.UserID), 0, globalFields); err == nil {
				successCount++
			}
			if err := doFriendOper(client, baseURL, int64(req.UserID), 1, globalFields); err == nil {
				successCount++
			}
		}
		global.GVA_DB.Model(&toolsModel.FanFollowRecord{}).
			Where("id = ?", record.ID).
			Update("success_count", successCount)
		time.Sleep(500 * time.Millisecond)
	}

	finalStatus := "completed"
	global.GVA_DB.Model(&toolsModel.FanFollowRecord{}).
		Where("id = ?", record.ID).
		Updates(map[string]interface{}{
			"success_count": successCount,
			"status":        finalStatus,
		})

	global.GVA_LOG.Info("FanFollow operation completed",
		zap.String("environmentKey", req.EnvironmentKey),
		zap.Uint64("userId", req.UserID),
		zap.String("operation", req.Operation),
		zap.Int("count", req.Count),
		zap.Int("successCount", successCount))

	return record.ID, successCount, nil
}

func (s *fanFollow) GetFanFollowList(info toolsReq.FanFollowSearch) ([]toolsModel.FanFollowRecord, int64, error) {
	var list []toolsModel.FanFollowRecord
	var total int64

	db := global.GVA_DB.Model(&toolsModel.FanFollowRecord{})
	if info.EnvironmentKey != "" {
		db = db.Where("environment_key = ?", info.EnvironmentKey)
	}
	if info.Operation != "" {
		db = db.Where("operation = ?", info.Operation)
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := db.Scopes((&commonReq.PageInfo{Page: info.Page, PageSize: info.PageSize}).Paginate()).Order("id desc").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}
