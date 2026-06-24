package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	toolsModel "github.com/flipped-aurora/gin-vue-admin/server/plugin/tools/model"
	toolsReq "github.com/flipped-aurora/gin-vue-admin/server/plugin/tools/model/request"
	"go.uber.org/zap"
)

type fanFollow struct{}

type friendOperBody struct {
	Param  map[string]interface{} `json:"param"`
	Global map[string]string      `json:"global"`
}

func (s *fanFollow) CreateFanFollow(req toolsReq.FanFollowCreate) (int, error) {
	if req.EnvironmentKey == "" {
		return 0, errors.New("环境Key不能为空")
	}
	if req.UserID == 0 {
		return 0, errors.New("用户ID不能为空")
	}
	if req.Count <= 0 {
		return 0, errors.New("执行次数必须大于0")
	}

	var env toolsModel.Environment
	err := global.GVA_DB.Where("key = ?", req.EnvironmentKey).First(&env)
	if err != nil {
		return 0, errors.New("未找到对应环境配置")
	}

	baseURL := fmt.Sprintf("%s:%d", env.Domain, env.Port)
	headers := map[string]string{
		"appid":       "1",
		"client":      "ios;15.4.1;iPhone11,8;5.27.0;5.27.0.96",
		"deviceid":    "1BDB215B-4CF5-4540-9BDE-C7B4D8574731",
		"packagename": "com.youzuo.miko.inner",
		"user-agent":  "yyzz/5.27.0 (com.youzuo.miko.inner; build:96; iOS 15.4.1) Alamofire/5.5.0",
		"version":     "5.27.0",
		"uid":         "200426370",
		"application": "myyw",
	}

	globalFields := map[string]string{
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
		"Uid":             strconv.FormatUint(req.UserID, 10),
		"User-Agent":      "yyzz/5.27.0 (com.youzuo.miko.inner; build:96; iOS 15.4.1) Alamofire/5.5.0",
		"Version":         "5.27.0",
	}

	client := &http.Client{Timeout: 10 * time.Second}
	successCount := 0

	for i := 0; i < req.Count; i++ {
		var oper int
		switch req.Operation {
		case "fans":
			oper = 1
		case "follow":
			oper = 0
		case "friend":
			oper = 0
		default:
			return 0, errors.New("不支持的操作类型")
		}

		run := func(op int, label string) error {
			body := friendOperBody{
				Param: map[string]interface{}{
					"id":   int64(12121),
					"oper": op,
				},
				Global: globalFields,
			}
			jsonData, err := json.Marshal(body)
			if err != nil {
				return err
			}
			payload := bytes.NewReader(jsonData)
			httpReq, err := http.NewRequest("POST", baseURL+"/testProxy/FriendExtObj/FriendOper", payload)
			if err != nil {
				return err
			}
			httpReq.Header.Set("Content-Type", "application/json")
			for k, v := range headers {
				httpReq.Header.Set(k, v)
			}
			resp, err := client.Do(httpReq)
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			_, err = io.ReadAll(resp.Body)
			return err
		}

		if req.Operation == "friend" {
			if err := run(0, "关注"); err == nil {
				successCount++
			}
			if err := run(1, "粉丝"); err == nil {
				successCount++
			}
		} else {
			if err := run(oper, req.Operation); err == nil {
				successCount++
			}
		}
		time.Sleep(3 * time.Second)
	}

	global.GVA_LOG.Info("FanFollow operation completed",
		zap.String("environmentKey", req.EnvironmentKey),
		zap.Uint64("userId", req.UserID),
		zap.String("operation", req.Operation),
		zap.Int("count", req.Count),
		zap.Int("successCount", successCount))

	return successCount, nil
}
