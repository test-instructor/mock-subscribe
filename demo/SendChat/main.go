package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	idsFile       = "/Users/taylor/Documents/HS/server/SendChat/ids.txt"
	maxGoroutines = 20
	//roomid        = "2606076"
	roomid    = "2617821"
	sendCount = 3000

	baseURL = "http://1.116.109.241:28196/testProxy"
)

type RequestBody struct {
	Param  interface{} `json:"param"`
	Global Global      `json:"global"`
}

type Global struct {
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

var chineseChars = []rune{
	'我', '你', '他', '她', '的', '是', '在', '有', '和', '了',
	'人', '这', '中', '大', '为', '上', '个', '国', '我', '们',
	'来', '到', '时', '说', '要', '就', '不', '会', '才', '可',
	'下', '过', '子', '也', '得', '着', '看', '发', '后', '作',
	'用', '里', '出', '道', '去', '行', '所', '然', '家', '种',
	'事', '成', '方', '多', '经', '么', '去', '好', '小', '还',
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

func buildGlobal(uid, roomid string) Global {
	return Global{
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
		Roomid:         roomid,
		Uid:            uid,
		UserAgent:      "yyzz/5.27.0 (com.youzuo.miko.inner; build:96; iOS 15.4.1) Alamofire/5.5.0",
		Version:        "5.27.0",
	}
}

func doRequest(client *http.Client, url string, body RequestBody, headers map[string]string) (string, error) {
	jsonData, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("json marshal error: %v", err)
	}
	payload := bytes.NewReader(jsonData)
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return "", fmt.Errorf("new request error: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	res, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("do request error: %v", err)
	}
	defer res.Body.Close()

	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("read body error: %v", err)
	}
	return string(respBody), nil
}

func buildHeaders(g Global) map[string]string {
	return map[string]string{
		"appid":       g.Appid,
		"client":      g.Client,
		"deviceid":    g.Deviceid,
		"packagename": g.PackageName,
		"user-agent":  g.UserAgent,
		"version":     g.Version,
		"uid":         g.Uid,
		"application": g.Application,
	}
}

func worker(id string, wg *sync.WaitGroup) {
	defer wg.Done()

	client := &http.Client{}
	roomidInt, _ := strconv.ParseInt(roomid, 10, 64)
	playerIdInt, _ := strconv.ParseInt(id, 10, 64)

	global := buildGlobal(id, roomid)

	enterURL := baseURL + "/RoomExtObj/EnterRoom2"
	enterBody := RequestBody{
		Param: map[string]interface{}{
			"room_id":        roomidInt,
			"device_type":    20,
			"is_long_socket": true,
		},
		Global: global,
	}
	resp, err := doRequest(client, enterURL, enterBody, buildHeaders(global))
	if err != nil {
		fmt.Printf("[%s] 进入房间失败: %v\n", id, err)
		return
	}
	fmt.Printf("[%s] 进入房间成功: %s\n", id, resp)

	collectionURL := baseURL + "/RoomExtObj/CollectionRoom"
	for _, status := range []int{1, 0} {
		collectionBody := RequestBody{
			Param: map[string]interface{}{
				"playerId": playerIdInt,
				"roomId":   roomidInt,
				"status":   status,
			},
			Global: global,
		}
		resp, err := doRequest(client, collectionURL, collectionBody, buildHeaders(global))
		if err != nil {
			fmt.Printf("[%s] CollectionRoom status=%d 失败: %v\n", id, status, err)
			continue
		}
		fmt.Printf("[%s] CollectionRoom status=%d 成功: %s\n", id, status, resp)
	}

	sendURL := baseURL + "/RoomExtObj/SendChat"
	for i := 1; i <= sendCount; i++ {
		msg := randomChineseMsg()
		options := map[string]interface{}{
			"type":      0,
			"emojiId":   0,
			"toId":      0,
			"giftNum":   0,
			"gameGlory": "",
			"giftId":    0,
			"createAt":  time.Now().Unix(),
			"goldVoice": 0,
		}
		data, err := json.Marshal(options)
		if err != nil {
			fmt.Printf("序列化失败: %v\n", err)
			return
		}
		sendBody := RequestBody{
			Param: map[string]interface{}{
				"content":    msg,
				"is_private": false,
				"options":    data,
			},
			Global: global,
		}
		resp, err := doRequest(client, sendURL, sendBody, buildHeaders(global))
		if err != nil {
			fmt.Printf("[%s] 第%d条发送失败: %v\n", id, i, err)
			continue
		}
		fmt.Printf("[%s] 第%d条发送成功 [%s]: %s\n", id, i, msg, resp)
		time.Sleep(100 * time.Millisecond)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	file, err := os.Open(idsFile)
	if err != nil {
		fmt.Printf("无法打开文件: %v\n", err)
		return
	}
	defer file.Close()

	var ids []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		id := strings.TrimSpace(scanner.Text())
		if id != "" {
			ids = append(ids, id)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("读取文件错误: %v\n", err)
		return
	}

	if maxGoroutines <= 0 || maxGoroutines > len(ids) {
		panic(fmt.Sprintf("maxGoroutines 值无效: %d (有效范围: 1~%d)", maxGoroutines, len(ids)))
	}

	selectedIDs := ids[:maxGoroutines]
	fmt.Printf("共读取到 %d 个 ID，将处理前 %d 个\n", len(ids), maxGoroutines)

	var wg sync.WaitGroup
	for _, id := range selectedIDs {
		wg.Add(1)
		go worker(id, &wg)
	}

	wg.Wait()
	fmt.Println("全部完成，程序退出")
}
