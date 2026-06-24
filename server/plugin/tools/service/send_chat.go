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
	"sync"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	commonReq "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	toolsModel "github.com/flipped-aurora/gin-vue-admin/server/plugin/tools/model"
	toolsReq "github.com/flipped-aurora/gin-vue-admin/server/plugin/tools/model/request"
	"go.uber.org/zap"
)

type sendChat struct{}

var chineseChars = []rune{
	'我', '你', '他', '她', '的', '是', '在', '有', '和', '了',
	'人', '这', '中', '大', '为', '上', '个', '国', '我', '们',
	'来', '到', '时', '说', '要', '就', '不', '会', '才', '可',
	'下', '过', '子', '也', '得', '着', '看', '发', '后', '作',
	'用', '里', '出', '道', '去', '行', '所', '然', '家', '种',
	'事', '成', '方', '多', '经', '么', '好', '小', '还',
	'感', '谢', '啊', '哈', '嗯', '呀', '哦', '喂', '嘿', '嘻',
}

func randomChineseMsg() string {
	n := rand.Intn(16) + 5
	runes := make([]rune, n)
	for i := range runes {
		runes[i] = chineseChars[rand.Intn(len(chineseChars))]
	}
	return string(runes)
}

type chatGlobal struct {
	Accept         string `json:"Accept"`
	AcceptEncoding string `json:"Accept-Encoding"`
	AcceptLanguage string `json:"Accept-Language"`
	Appid          string `json:"appid"`
	Application    string `json:"application"`
	Client         string `json:"client"`
	Connection     string `json:"Connection"`
	Deviceid       string `json:"deviceid"`
	Host           string `json:"Host"`
	MarketFlag     string `json:"market_flag"`
	PackageName    string `json:"packageName"`
	Roomid         string `json:"roomid"`
	Uid            string `json:"uid"`
	UserAgent      string `json:"User-Agent"`
	Version        string `json:"version"`
}

type chatBody struct {
	Param  interface{} `json:"param"`
	Global chatGlobal  `json:"global"`
}

func (s *sendChat) CreateSendChatTask(req toolsReq.SendChatCreate) (uint, error) {
	if req.RoomID == "" {
		return 0, errors.New("房间ID不能为空")
	}
	if req.EnvironmentKey == "" {
		return 0, errors.New("环境Key不能为空")
	}
	if req.AccountCount <= 0 {
		return 0, errors.New("账号数量必须大于0")
	}
	if req.MsgCountPerAccount <= 0 {
		return 0, errors.New("每账号消息数必须大于0")
	}
	if req.MsgInterval < 0 {
		req.MsgInterval = 0
	}

	task := toolsModel.SendChatTask{
		RoomID:             req.RoomID,
		EnvironmentKey:     req.EnvironmentKey,
		AccountCount:       req.AccountCount,
		MsgCountPerAccount: req.MsgCountPerAccount,
		MsgInterval:        req.MsgInterval,
		Status:             "running",
		SuccessCount:       0,
	}
	if err := global.GVA_DB.Create(&task).Error; err != nil {
		return 0, err
	}

	go s.runSendChatTask(task.ID, req)

	return task.ID, nil
}

func (s *sendChat) runSendChatTask(taskID uint, req toolsReq.SendChatCreate) {
	defer func() {
		global.GVA_DB.Model(&toolsModel.SendChatTask{}).Where("id = ?", taskID).
			Update("status", "completed")
	}()

	var env toolsModel.Environment
	if err := global.GVA_DB.Where("key = ?", req.EnvironmentKey).First(&env).Error; err != nil {
		global.GVA_LOG.Error("SendChat task failed: env not found", zap.String("key", req.EnvironmentKey))
		return
	}

	var userIDs []uint64
	if err := global.GVA_DB.Model(&toolsModel.UserRelation{}).
		Where("environment_key = ?", req.EnvironmentKey).
		Order("id asc").
		Limit(req.AccountCount).
		Pluck("user_id", &userIDs).Error; err != nil || len(userIDs) == 0 {
		global.GVA_LOG.Error("SendChat task failed: no user IDs found", zap.String("key", req.EnvironmentKey))
		return
	}

	roomidInt, _ := strconv.ParseInt(req.RoomID, 10, 64)
	interval := time.Duration(req.MsgInterval) * time.Millisecond
	if interval == 0 {
		interval = 100 * time.Millisecond
	}

	baseURL := fmt.Sprintf("%s:%d%s", env.Domain, env.Port, "/testProxy")

	headers := map[string]string{
		"appid":       "1",
		"client":      "ios;15.4.1;iPhone11,8;5.27.0;5.27.0.96",
		"deviceid":    "1BDB215B-4CF5-4540-9BDE-C7B4D8574731",
		"packagename": "com.youzuo.miko.inner",
		"user-agent":  "yyzz/5.27.0 (com.youzuo.miko.inner; build:96; iOS 15.4.1) Alamofire/5.5.0",
		"version":     "5.27.0",
		"application": "zanzan",
	}

	successCount := 0
	var mu sync.Mutex
	var wg sync.WaitGroup
	limiter := make(chan struct{}, 10)

	for _, uid := range userIDs {
		wg.Add(1)
		go func(userID uint64) {
			defer wg.Done()
			limiter <- struct{}{}
			defer func() { <-limiter }()

			uidStr := strconv.FormatUint(userID, 10)
			client := &http.Client{Timeout: 30 * time.Second}
			globalFields := chatGlobal{
				Accept:         "*/*",
				AcceptEncoding: "br;q=1.0, gzip;q=0.9, deflate;q=0.8",
				AcceptLanguage: "zh-Hans-CN;q=1.0",
				Appid:          "1",
				Application:    "zanzan",
				Client:         "ios;15.4.1;iPhone11,8;5.27.0;5.27.0.96",
				Connection:     "keep-alive",
				Deviceid:       "1BDB215B-4CF5-4540-9BDE-C7B4D8574731",
				Host:           "api-test.miyafm.com:443",
				MarketFlag:     "1",
				PackageName:    "com.youzuo.miko.inner",
				Roomid:         req.RoomID,
				Uid:            uidStr,
				UserAgent:      "yyzz/5.27.0 (com.youzuo.miko.inner; build:96; iOS 15.4.1) Alamofire/5.5.0",
				Version:        "5.27.0",
			}

			doReq := func(path string, param map[string]interface{}) error {
				body := chatBody{Param: param, Global: globalFields}
				jsonData, _ := json.Marshal(body)
				payload := bytes.NewReader(jsonData)
				httpReq, _ := http.NewRequest("POST", baseURL+path, payload)
				httpReq.Header.Set("Content-Type", "application/json")
				for k, v := range headers {
					httpReq.Header.Set(k, v)
				}
				httpReq.Header.Set("uid", uidStr)
				resp, err := client.Do(httpReq)
				if err != nil {
					return err
				}
				defer resp.Body.Close()
				_, err = io.ReadAll(resp.Body)
				return err
			}

			enterBody := map[string]interface{}{
				"room_id":        roomidInt,
				"device_type":    20,
				"is_long_socket": true,
			}
			doReq("/RoomExtObj/EnterRoom2", enterBody)

			sendBody := map[string]interface{}{
				"is_private": false,
				"options":    json.RawMessage(`{"type":0,"emojiId":0,"toId":0,"giftNum":0,"gameGlory":"","giftId":0,"createAt":0,"goldVoice":0}`),
			}
			for i := 0; i < req.MsgCountPerAccount; i++ {
				sendBody["content"] = randomChineseMsg()
				if err := doReq("/RoomExtObj/SendChat", sendBody); err == nil {
					mu.Lock()
					successCount++
					global.GVA_DB.Model(&toolsModel.SendChatTask{}).Where("id = ?", taskID).Update("success_count", successCount)
					mu.Unlock()
				}
				time.Sleep(interval)
			}
		}(uid)
	}

	wg.Wait()

	global.GVA_LOG.Info("SendChat task completed",
		zap.Uint("taskId", taskID),
		zap.Int("successCount", successCount))
}

func (s *sendChat) StopSendChatTask(id uint) error {
	if id == 0 {
		return errors.New("任务ID不能为空")
	}
	result := global.GVA_DB.Model(&toolsModel.SendChatTask{}).Where("id = ? AND status = ?", id, "running").
		Update("status", "stopped")
	if result.RowsAffected == 0 {
		return errors.New("任务不存在或已停止")
	}
	return nil
}

func (s *sendChat) GetSendChatTaskList(info toolsReq.SendChatSearch) ([]toolsModel.SendChatTask, int64, error) {
	var list []toolsModel.SendChatTask
	var total int64

	db := global.GVA_DB.Model(&toolsModel.SendChatTask{})
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := db.Scopes((&commonReq.PageInfo{Page: info.Page, PageSize: info.PageSize}).Paginate()).Order("id desc").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}
