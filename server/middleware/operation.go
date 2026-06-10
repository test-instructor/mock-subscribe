package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/utils"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var respPool sync.Pool
var bufferSize = 1024

func init() {
	respPool.New = func() interface{} {
		return make([]byte, bufferSize)
	}
}

// OperationRecord 创建操作记录中间件，写入 sys_operation_records 表。
// 其中 UserID 优先从 JWT 解析，其次从请求头 x-user-id 解析，仍未取到则记为 0。
func OperationRecord() gin.HandlerFunc {
	return func(c *gin.Context) {
		recordOperation(c, resolveOperationUserID(c))
	}
}

// OperationRecordWithUserID 创建操作记录中间件，强制把 UserID 写为指定值。
// 适用于 mock / 系统级等没有真实登录用户的接口，便于把操作历史归到固定系统用户。
func OperationRecordWithUserID(forcedUserID int) gin.HandlerFunc {
	return func(c *gin.Context) {
		recordOperation(c, forcedUserID)
	}
}

// resolveOperationUserID 优先从 JWT 取登录用户 ID，其次从 x-user-id 请求头解析。
func resolveOperationUserID(c *gin.Context) int {
	claims, _ := utils.GetClaims(c)
	if claims != nil && claims.BaseClaims.ID != 0 {
		return int(claims.BaseClaims.ID)
	}
	id, err := strconv.Atoi(c.Request.Header.Get("x-user-id"))
	if err != nil {
		return 0
	}
	return id
}

func recordOperation(c *gin.Context, userId int) {
	var body []byte
	if c.Request.Method != http.MethodGet {
		var err error
		body, err = io.ReadAll(c.Request.Body)
		if err != nil {
			global.GVA_LOG.Error("read body from request error:", zap.Error(err))
		} else {
			c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		}
	} else {
		query := c.Request.URL.RawQuery
		query, _ = url.QueryUnescape(query)
		split := strings.Split(query, "&")
		m := make(map[string]string)
		for _, v := range split {
			kv := strings.Split(v, "=")
			if len(kv) == 2 {
				m[kv[0]] = kv[1]
			}
		}
		body, _ = json.Marshal(&m)
	}
	record := system.SysOperationRecord{
		Ip:     c.ClientIP(),
		Method: c.Request.Method,
		Path:   c.Request.URL.Path,
		Agent:  c.Request.UserAgent(),
		Body:   "",
		UserID: userId,
	}

	// 上传文件时候 中间件日志进行裁断操作
	if strings.Contains(c.GetHeader("Content-Type"), "multipart/form-data") {
		record.Body = "[文件]"
	} else {
		if len(body) > bufferSize {
			record.Body = "[超出记录长度]"
		} else {
			record.Body = string(body)
		}
	}

	writer := responseBodyWriter{
		ResponseWriter: c.Writer,
		body:           &bytes.Buffer{},
	}
	c.Writer = writer
	now := time.Now()

	c.Next()

	latency := time.Since(now)
	record.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()
	record.Status = c.Writer.Status()
	record.Latency = latency
	record.Resp = writer.body.String()

	if strings.Contains(c.Writer.Header().Get("Pragma"), "public") ||
		strings.Contains(c.Writer.Header().Get("Expires"), "0") ||
		strings.Contains(c.Writer.Header().Get("Cache-Control"), "must-revalidate, post-check=0, pre-check=0") ||
		strings.Contains(c.Writer.Header().Get("Content-Type"), "application/force-download") ||
		strings.Contains(c.Writer.Header().Get("Content-Type"), "application/octet-stream") ||
		strings.Contains(c.Writer.Header().Get("Content-Type"), "application/vnd.ms-excel") ||
		strings.Contains(c.Writer.Header().Get("Content-Type"), "application/download") ||
		strings.Contains(c.Writer.Header().Get("Content-Disposition"), "attachment") ||
		strings.Contains(c.Writer.Header().Get("Content-Transfer-Encoding"), "binary") {
		if len(record.Resp) > bufferSize {
			// 截断
			record.Body = "超出记录长度"
		}
	}
	// SysOperationRecord 上的 User SysUser 是 belongs-to 关系，GORM v2 在 Create 阶段会
	// 主动去装载 / 关联这个 User，UserID 强制为 3 时（sys_users 中可能并不存在）会触发额外
	// 的 SELECT / 关联 INSERT，导致 sys_operation_records 写入静默失败。用 Omit("User")
	// 切断这条路径，只写 user_id 外键列。
	if err := global.GVA_DB.Omit("User").Create(&record).Error; err != nil {
		global.GVA_LOG.Error("create operation record error:",
			zap.Error(err),
			zap.String("path", record.Path),
			zap.Int("user_id", record.UserID),
		)
	}
}

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}
