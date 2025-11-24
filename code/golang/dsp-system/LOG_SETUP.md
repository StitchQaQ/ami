# 日志系统配置完成 ✅

## 🎉 功能概述

已成功封装日志系统，具备以下功能：

- ✅ **多级别日志**：DEBUG, INFO, WARN, ERROR, FATAL
- ✅ **自动轮转**：按大小（100MB）和时间自动切分
- ✅ **自动清理**：默认保留 7 天，自动删除过期日志
- ✅ **自动压缩**：旧日志文件自动压缩节省空间
- ✅ **双重输出**：同时输出到文件和控制台
- ✅ **Gin 集成**：自动记录所有 HTTP 请求

## 📁 文件结构

```
dsp-system/
├── config/
│   └── config.go          # 新增 LogConfig 配置
├── logger/
│   ├── logger.go          # 日志系统核心代码
│   ├── example.go         # 使用示例
│   └── README.md          # 详细文档
├── logs/                  # 日志文件目录（自动创建）
│   └── dsp-system.log     # 当前日志文件
└── main.go                # 已集成日志系统
```

## 🚀 快速开始

### 1. 日志配置（环境变量）

```bash
# 设置日志级别
export LOG_LEVEL=info              # debug, info, warn, error

# 设置日志文件路径
export LOG_FILE_PATH=logs/dsp-system.log
```

### 2. 运行服务

```bash
cd /Users/shengli/Code/www/blog/code/golang/dsp-system

# 使用默认配置（info级别）
go run main.go

# 使用 debug 级别（查看详细日志）
LOG_LEVEL=debug go run main.go

# 使用 warn 级别（只看警告和错误）
LOG_LEVEL=warn go run main.go
```

### 3. 查看日志

```bash
# 实时查看日志
tail -f logs/dsp-system.log

# 查看错误日志
grep "level=error" logs/dsp-system.log

# 查看警告和错误
grep -E "level=(warn|error)" logs/dsp-system.log

# 统计不同级别的日志数量
grep -c "level=info" logs/dsp-system.log
grep -c "level=error" logs/dsp-system.log
```

## 💻 代码使用示例

### 基础用法

```go
import "dsp-system/logger"

// 基本日志
logger.Info("服务启动成功")
logger.Infof("监听端口: %d", 8080)

// 错误日志
logger.Error("数据库连接失败")
logger.Errorf("错误详情: %v", err)

// 警告日志
logger.Warn("预算不足")
logger.Warnf("剩余预算: %.2f CNY", budget)

// 调试日志（生产环境不输出）
logger.Debug("调试信息")
logger.Debugf("变量值: %+v", data)
```

### 结构化日志（推荐）

```go
import "github.com/sirupsen/logrus"

// 单个字段
logger.WithField("request_id", requestID).Info("处理请求")

// 多个字段
logger.WithFields(logrus.Fields{
    "request_id": requestID,
    "user_id":    userID,
    "duration":   duration.Milliseconds(),
    "status":     "success",
}).Info("请求完成")
```

### 在 Handler 中使用

```go
func (h *RTBHandler) HandleBidRequest(c *gin.Context) {
    // 记录请求
    logger.WithFields(logrus.Fields{
        "request_id": bidRequest.ID,
        "imp_count":  len(bidRequest.Imp),
    }).Info("收到竞价请求")
    
    // 错误处理
    if err != nil {
        logger.WithFields(logrus.Fields{
            "request_id": bidRequest.ID,
            "error":      err.Error(),
        }).Error("竞价处理失败")
        return
    }
    
    // 成功日志
    logger.WithFields(logrus.Fields{
        "request_id":  bidRequest.ID,
        "bid_count":   len(bids),
        "duration_ms": time.Since(startTime).Milliseconds(),
    }).Info("竞价成功")
}
```

## 📊 日志输出示例

```
time="2024-11-20 16:40:47" level=info msg="日志系统初始化完成: Level=info, File=logs/dsp-system.log"
time="2024-11-20 16:40:47" level=info msg="配置加载完成"
time="2024-11-20 16:40:47" level=info msg="RPC客户端初始化完成"
time="2024-11-20 16:40:47" level=info msg="DSP服务启动: http://localhost:8080"
```

## 🔧 日志配置详解

### config/config.go 中的配置

```go
type LogConfig struct {
    Level      string // 日志级别: debug, info, warn, error
    FilePath   string // 日志文件路径
    MaxSize    int    // 单个文件最大大小(MB) - 默认100MB
    MaxBackups int    // 保留的旧日志文件数量 - 默认7个
    MaxAge     int    // 保留天数 - 默认7天
    Compress   bool   // 是否压缩旧日志 - 默认true
}
```

### 默认配置

- **Level**: `info` - 记录 INFO 及以上级别的日志
- **FilePath**: `logs/dsp-system.log` - 日志文件路径
- **MaxSize**: `100` MB - 单个文件最大 100MB
- **MaxBackups**: `7` - 最多保留 7 个备份
- **MaxAge**: `7` 天 - 保留 7 天
- **Compress**: `true` - 自动压缩旧日志

## 📈 日志级别说明

| 级别 | 用途 | 示例场景 | 生产环境 |
|-----|------|---------|---------|
| DEBUG | 详细调试信息 | 变量值、SQL语句、函数调用 | ❌ 不建议 |
| INFO | 常规业务信息 | 服务启动、请求处理、正常流程 | ✅ 推荐 |
| WARN | 警告信息 | 预算不足、性能问题、降级逻辑 | ✅ 推荐 |
| ERROR | 错误信息 | 数据库错误、RPC失败、业务异常 | ✅ 必须 |
| FATAL | 致命错误 | 启动失败、无法恢复的错误 | ✅ 必须 |

## 🗂️ 日志文件管理

### 自动轮转

当日志文件达到 100MB 时，会自动创建新文件：

```
logs/
├── dsp-system.log              # 当前日志（正在写入）
├── dsp-system-2024-11-20.log   # 昨天的日志
└── dsp-system-2024-11-19.log.gz # 前天的日志（已压缩）
```

### 自动清理

- 超过 7 天的日志文件会被自动删除
- 超过 7 个备份文件会删除最旧的
- 旧日志会自动压缩为 `.gz` 格式

## 🎯 最佳实践

### 1. 生产环境配置

```bash
export LOG_LEVEL=info
export LOG_FILE_PATH=/var/log/dsp-system/dsp-system.log
```

### 2. 开发环境配置

```bash
export LOG_LEVEL=debug
export LOG_FILE_PATH=logs/dsp-system.log
```

### 3. 敏感信息脱敏

```go
// ❌ 不要记录敏感信息
logger.Infof("用户密码: %s", password)

// ✅ 只记录必要信息
logger.WithField("user_id", userID).Info("用户登录成功")
```

### 4. 结构化日志

```go
// ❌ 避免字符串拼接
logger.Info("用户" + userID + "在" + time.Now().String() + "执行了操作")

// ✅ 使用结构化日志
logger.WithFields(logrus.Fields{
    "user_id":   userID,
    "timestamp": time.Now().Unix(),
    "operation": "bid_request",
}).Info("用户操作")
```

### 5. 性能优化

```go
// ❌ 避免在循环中使用 Debug
for _, ad := range ads {
    logger.Debugf("处理广告: %+v", ad) // 可能影响性能
}

// ✅ 记录汇总信息
logger.WithField("ad_count", len(ads)).Debug("开始处理广告")
```

## 📚 更多文档

- 详细使用文档：`logger/README.md`
- 代码示例：`logger/example.go`
- 配置说明：`config/config.go`

## ✅ 测试日志系统

```bash
# 1. 启动服务
go run main.go

# 2. 在另一个终端查看日志
tail -f logs/dsp-system.log

# 3. 发送测试请求
curl http://localhost:8080/health

# 4. 查看日志输出
cat logs/dsp-system.log
```

## 🔍 日志分析命令

```bash
# 统计各级别日志数量
grep -c "level=debug" logs/dsp-system.log
grep -c "level=info" logs/dsp-system.log
grep -c "level=warn" logs/dsp-system.log
grep -c "level=error" logs/dsp-system.log

# 查找特定时间的日志
grep "2024-11-20 16:" logs/dsp-system.log

# 查找包含特定字段的日志
grep "request_id" logs/dsp-system.log

# 查看最新的错误
grep "level=error" logs/dsp-system.log | tail -10
```

## 🎊 完成！

日志系统已经完全集成到项目中，可以直接使用了！

