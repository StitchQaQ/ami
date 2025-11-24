package logger

import (
	"time"

	"github.com/sirupsen/logrus"
)

// Example 展示日志系统的各种使用方式
func Example() {
	// 1. 基本日志输出
	Debug("这是调试信息 - 用于开发环境")
	Info("这是常规信息 - 记录正常流程")
	Warn("这是警告信息 - 需要注意但不影响运行")
	Error("这是错误信息 - 业务异常或错误")

	// 2. 格式化输出
	userID := "user_12345"
	Debugf("处理用户请求: %s", userID)
	
	price := 5.68
	Infof("竞价成功，出价: %.2f CNY", price)
	
	budget := 100.50
	Warnf("预算不足，剩余: %.2f CNY", budget)
	
	err := "database connection timeout"
	Errorf("数据库错误: %v", err)

	// 3. 结构化日志（推荐）
	WithField("request_id", "req_001").Info("开始处理请求")
	
	WithFields(logrus.Fields{
		"request_id": "req_001",
		"user_id":    "user_12345",
		"duration_ms": 45,
	}).Info("请求处理完成")

	// 4. 业务场景示例
	requestID := "bid_req_20241120_001"
	impCount := 3
	
	WithFields(logrus.Fields{
		"request_id": requestID,
		"imp_count":  impCount,
		"timestamp":  time.Now().Unix(),
	}).Info("收到RTB竞价请求")

	// 5. 错误追踪
	campaignID := "campaign_001"
	WithFields(logrus.Fields{
		"campaign_id": campaignID,
		"error":       "insufficient budget",
	}).Warn("预算不足，跳过该广告")

	// 6. 性能监控
	startTime := time.Now()
	// ... 执行业务逻辑 ...
	duration := time.Since(startTime)
	
	WithFields(logrus.Fields{
		"operation":   "ad_selection",
		"duration_ms": duration.Milliseconds(),
		"status":      "success",
	}).Info("广告选择完成")

	// 7. 不同级别的使用场景
	
	// DEBUG: 详细的调试信息（生产环境不输出）
	WithFields(logrus.Fields{
		"sql":    "SELECT * FROM ads WHERE id = ?",
		"params": []interface{}{123},
	}).Debug("执行数据库查询")

	// INFO: 重要的业务节点
	WithField("ad_id", "ad_001").Info("广告投放成功")

	// WARN: 需要关注的异常情况
	WithFields(logrus.Fields{
		"cache_key": "user_profile:12345",
		"reason":    "cache miss",
	}).Warn("缓存未命中，降级到数据库查询")

	// ERROR: 需要处理的错误
	WithFields(logrus.Fields{
		"rpc_service": "user_service",
		"error":       "connection refused",
		"retry_count": 3,
	}).Error("RPC调用失败")
}

// ExampleHTTPRequest HTTP请求日志示例
func ExampleHTTPRequest(method, path string, statusCode int, duration time.Duration) {
	WithFields(logrus.Fields{
		"method":      method,
		"path":        path,
		"status_code": statusCode,
		"duration_ms": duration.Milliseconds(),
	}).Info("HTTP请求")
}

// ExampleBidRequest 竞价请求日志示例
func ExampleBidRequest(requestID string, impCount int, hasUser bool) {
	WithFields(logrus.Fields{
		"request_id": requestID,
		"imp_count":  impCount,
		"has_user":   hasUser,
	}).Info("竞价请求")
}

// ExampleAdSelection 广告选择日志示例
func ExampleAdSelection(requestID string, candidateCount int, selectedCount int, duration time.Duration) {
	WithFields(logrus.Fields{
		"request_id":      requestID,
		"candidate_count": candidateCount,
		"selected_count":  selectedCount,
		"duration_ms":     duration.Milliseconds(),
	}).Info("广告选择")
}

// ExampleError 错误日志示例
func ExampleError(operation string, err error, context map[string]interface{}) {
	fields := logrus.Fields{
		"operation": operation,
		"error":     err.Error(),
	}
	
	// 添加上下文信息
	for k, v := range context {
		fields[k] = v
	}
	
	WithFields(fields).Error("操作失败")
}

