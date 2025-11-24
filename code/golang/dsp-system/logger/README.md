# Logger 日志系统

## 功能特性

- ✅ 支持多种日志级别：DEBUG, INFO, WARN, ERROR, FATAL
- ✅ 自动日志轮转（按大小和时间）
- ✅ 自动清理过期日志（默认保留7天）
- ✅ 同时输出到文件和控制台
- ✅ 支持压缩旧日志文件
- ✅ 与 Gin 框架无缝集成

## 配置说明

在 `config/config.go` 中配置日志参数：

```go
type LogConfig struct {
    Level      string // 日志级别: debug, info, warn, error
    FilePath   string // 日志文件路径
    MaxSize    int    // 单个文件最大大小(MB)
    MaxBackups int    // 保留的旧日志文件数量
    MaxAge     int    // 保留天数
    Compress   bool   // 是否压缩旧日志
}
```

### 环境变量配置

```bash
# 日志级别 (debug/info/warn/error)
export LOG_LEVEL=info

# 日志文件路径
export LOG_FILE_PATH=logs/dsp-system.log
```

## 使用方法

### 1. 初始化日志系统

```go
import (
    "dsp-system/config"
    "dsp-system/logger"
)

func main() {
    cfg := config.LoadConfig()
    
    // 初始化日志系统
    if err := logger.Init(&cfg.Log); err != nil {
        panic(err)
    }
}
```

### 2. 基本日志输出

```go
// Debug 级别
logger.Debug("这是调试信息")
logger.Debugf("用户ID: %s", userID)

// Info 级别
logger.Info("服务启动成功")
logger.Infof("监听端口: %d", port)

// Warn 级别
logger.Warn("预算即将耗尽")
logger.Warnf("剩余预算: %.2f", budget)

// Error 级别
logger.Error("数据库连接失败")
logger.Errorf("连接错误: %v", err)

// Fatal 级别（会退出程序）
logger.Fatal("致命错误")
logger.Fatalf("启动失败: %v", err)
```

### 3. 结构化日志

```go
import "github.com/sirupsen/logrus"

// 单个字段
logger.WithField("request_id", requestID).Info("处理请求")

// 多个字段
logger.WithFields(logrus.Fields{
    "request_id": requestID,
    "user_id":    userID,
    "duration":   duration,
}).Info("请求完成")
```

### 4. 在 Handler 中使用

```go
func (h *RTBHandler) HandleBidRequest(c *gin.Context) {
    logger.WithFields(logrus.Fields{
        "request_id": bidRequest.ID,
        "imp_count":  len(bidRequest.Imp),
    }).Info("收到竞价请求")
    
    // 业务逻辑...
    
    if err != nil {
        logger.WithField("request_id", bidRequest.ID).
               Errorf("竞价处理失败: %v", err)
        return
    }
}
```

## 日志级别说明

| 级别  | 用途 | 示例场景 |
|------|------|---------|
| DEBUG | 调试信息 | 详细的变量值、执行流程 |
| INFO  | 常规信息 | 服务启动、请求处理、正常业务流程 |
| WARN  | 警告信息 | 预算不足、性能问题、可恢复的错误 |
| ERROR | 错误信息 | 数据库错误、RPC调用失败、业务异常 |
| FATAL | 致命错误 | 启动失败、无法恢复的错误（会退出程序）|

## 日志文件管理

### 日志轮转规则

- **按大小轮转**：单个日志文件达到 MaxSize (默认100MB) 时自动创建新文件
- **按时间清理**：自动删除超过 MaxAge (默认7天) 的旧日志
- **备份数量限制**：最多保留 MaxBackups (默认7个) 个旧日志文件

### 日志文件命名

```
logs/
├── dsp-system.log              # 当前日志文件
├── dsp-system-2024-11-19.log   # 旧日志（未压缩）
└── dsp-system-2024-11-18.log.gz # 旧日志（已压缩）
```

### 查看日志

```bash
# 实时查看日志
tail -f logs/dsp-system.log

# 查看最近100行
tail -n 100 logs/dsp-system.log

# 搜索错误日志
grep "ERROR" logs/dsp-system.log

# 统计错误数量
grep -c "ERROR" logs/dsp-system.log
```

## 性能优化

### 生产环境建议

```bash
# 设置为 info 或 warn 级别
export LOG_LEVEL=info

# 开启压缩节省磁盘空间
# 在 config.go 中设置 Compress: true
```

### 开发环境建议

```bash
# 设置为 debug 级别查看详细信息
export LOG_LEVEL=debug
```

## 与 Gin 集成

日志系统已与 Gin 框架集成，所有 HTTP 请求都会自动记录：

```
[GIN] 2024/11/20 - 15:04:05 | 200 |     2.345ms |  127.0.0.1 | POST "/bid"
```

## 注意事项

1. **避免在高频代码中使用 Debug 日志**，会影响性能
2. **敏感信息脱敏**：不要记录密码、token 等敏感信息
3. **日志级别控制**：生产环境使用 info 或 warn 级别
4. **磁盘空间监控**：定期检查日志占用的磁盘空间

## 示例

完整示例请参考 `main.go` 和 `handler/rtb_handler.go`。

