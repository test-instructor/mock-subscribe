package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// 定义请求体结构
type RequestBody struct {
	Param  Param  `json:"param"`
	Global Global `json:"global"`
}

type Param struct {
	Id   int64 `json:"id"`
	Oper int   `json:"oper"`
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
	XToken         string `json:"X-Token"`
}

func main() {
	// 读取ids.txt文件
	file, err := os.Open("/Users/taylor/Documents/HS/server/finds/ids.txt")
	if err != nil {
		fmt.Printf("无法打开文件: %v\n", err)
		return
	}
	defer file.Close()
	// 43.138.249.156:29349		预发
	// 1.116.109.241:28196		测试
	// 创建HTTP客户端
	client := &http.Client{}
	url := "http://43.138.249.156:29349/testProxy/FriendExtObj/FriendOper"
	method := "POST"

	// 初始化基础请求体（除了id之外的固定部分）
	baseBody := RequestBody{
		Param: Param{
			Oper: 1, // 固定操作值
			Id:   12121,
		},
		Global: Global{
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
			Roomid:         "0",
			Uid:            "5596158",
			UserAgent:      "yyzz/5.27.0 (com.youzuo.miko.inner; build:96; iOS 15.4.1) Alamofire/5.5.0",
			Version:        "5.27.0",
		},
	}

	// 逐行读取ID并处理
	scanner := bufio.NewScanner(file)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		idStr := strings.TrimSpace(scanner.Text())
		if idStr == "" { // 跳过空行
			continue
		}

		// 转换ID为整数
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			fmt.Printf("第%d行ID格式错误: %v\n", lineNum, err)
			continue
		}

		// 替换当前ID
		//baseBody.Param.Id = id
		baseBody.Global.Uid = idStr
		// 转换为JSON
		jsonData, err := json.MarshalIndent(baseBody, "", "    ")
		if err != nil {
			fmt.Printf("第%d行JSON序列化失败: %v\n", lineNum, err)
			continue
		}

		// 创建请求
		payload := strings.NewReader(string(jsonData))
		req, err := http.NewRequest(method, url, payload)
		if err != nil {
			fmt.Printf("第%d行创建请求失败: %v\n", lineNum, err)
			continue
		}

		// 设置请求头
		req.Header.Add("appid", "1")
		req.Header.Add("client", "ios;15.4.1;iPhone11,8;5.27.0;5.27.0.96")
		req.Header.Add("deviceid", "1BDB215B-4CF5-4540-9BDE-C7B4D8574731")
		req.Header.Add("packagename", "com.youzuo.miko.inner")
		req.Header.Add("user-agent", "yyzz/5.27.0 (com.youzuo.miko.inner; build:96; iOS 15.4.1) Alamofire/5.5.0")
		req.Header.Add("version", "5.27.0")
		req.Header.Add("uid", "200426370")
		req.Header.Add("application", "myyw")
		req.Header.Add("Content-Type", "application/json")

		// 发送请求
		res, err := client.Do(req)
		if err != nil {
			fmt.Printf("第%d行请求失败: %v\n", lineNum, err)
			continue
		}

		// 读取响应
		body, err := io.ReadAll(res.Body)
		res.Body.Close() // 及时关闭body
		if err != nil {
			fmt.Printf("第%d行读取响应失败: %v\n", lineNum, err)
			continue
		}

		// 输出结果
		fmt.Printf("第%d行(ID: %d) 响应: %s\n", lineNum, id, string(body))
		time.Sleep(3 * time.Second)
	}

	// 检查扫描错误
	if err := scanner.Err(); err != nil {
		fmt.Printf("读取文件错误: %v\n", err)
	}
}
