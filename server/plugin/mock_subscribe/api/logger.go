package api

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ctxKey 用于在 gin.Context 中传递 traceID
const (
	traceIDKey    = "mock_subscribe_trace_id"
	handlerKey    = "mock_subscribe_handler"
	traceIDHeader = "X-Trace-Id"
)

// newTraceID 生成一个 16 字节的十六进制 traceID
func newTraceID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		// 退化方案：使用时间戳
		return time.Now().Format("20060102150405.000000")
	}
	return hex.EncodeToString(b)
}

// GetTraceID 从 gin.Context 中获取当前请求的 traceID，不存在则生成一个
func GetTraceID(c *gin.Context) string {
	if v, ok := c.Get(traceIDKey); ok {
		if s, ok := v.(string); ok && s != "" {
			return s
		}
	}
	if v := c.GetHeader(traceIDHeader); v != "" {
		c.Set(traceIDKey, v)
		return v
	}
	tid := newTraceID()
	c.Set(traceIDKey, tid)
	return tid
}

// LogRequest 记录请求入口的日志，包含请求内容、调用链起点、客户端信息等
func LogRequest(c *gin.Context, handler string, request any) {
	traceID := GetTraceID(c)
	c.Set(handlerKey, handler)
	c.Writer.Header().Set(traceIDHeader, traceID)
	global.GVA_LOG.Info("【请求入口】mock_subscribe API",
		zap.String("trace_id", traceID),
		zap.String("handler", handler),
		zap.String("method", c.Request.Method),
		zap.String("path", c.Request.URL.Path),
		zap.String("query", c.Request.URL.RawQuery),
		zap.String("client_ip", c.ClientIP()),
		zap.Any("request", request),
	)
}

// LogServiceCall 记录调用 service 层的日志，用于串联调用链
func LogServiceCall(c *gin.Context, service, action string, extra ...zap.Field) {
	traceID := GetTraceID(c)
	fields := []zap.Field{
		zap.String("trace_id", traceID),
		zap.String("handler", c.GetString(handlerKey)),
		zap.String("service", service),
		zap.String("action", action),
	}
	fields = append(fields, extra...)
	global.GVA_LOG.Info("【调用链路】service 调用", fields...)
}

// LogError 记录处理过程中的错误
func LogError(c *gin.Context, action string, err error, extra ...zap.Field) {
	if err == nil {
		return
	}
	traceID := GetTraceID(c)
	fields := []zap.Field{
		zap.String("trace_id", traceID),
		zap.String("handler", c.GetString(handlerKey)),
		zap.String("action", action),
		zap.Error(err),
	}
	fields = append(fields, extra...)
	global.GVA_LOG.Error("【错误内容】处理失败", fields...)
}

// LogResponse 记录请求出口的日志，包含响应内容与耗时
func LogResponse(c *gin.Context, action string, response any, start time.Time) {
	traceID := GetTraceID(c)
	cost := time.Since(start)
	global.GVA_LOG.Info("【请求出口】mock_subscribe API",
		zap.String("trace_id", traceID),
		zap.String("handler", c.GetString(handlerKey)),
		zap.String("action", action),
		zap.Any("response", response),
		zap.Duration("cost", cost),
	)
}
