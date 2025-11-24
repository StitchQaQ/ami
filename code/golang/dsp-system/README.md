# DSP System - 广告竞价系统

基于 Gin + gRPC 的 DSP（需求方平台）广告竞价系统示例项目。

## 架构说明

```
dsp-system/
├── api/                     # OpenRTB 协议定义
│   └── openrtb.go           # BidRequest, BidResponse 等数据结构
├── handler/                 # HTTP 请求处理层
│   └── rtb_handler.go       # 处理ADX的RTB竞价请求
├── service/                 # 业务逻辑层
│   ├── bid_service.go       # 竞价决策核心（调用算法+预算校验）
│   └── ad_select.go         # 广告素材匹配（基于用户标签）
├── rpc/                     # gRPC 客户端
│   ├── user_client.go       # 调用用户画像服务（获取用户标签）
│   └── budget_client.go     # 调用预算服务（校验剩余预算）
├── repository/              # 数据访问层
│   ├── redis_cache.go       # Redis操作封装
│   └── clickhouse_repo.go   # 日志存储
├── config/                  # 配置管理
│   └── config.go            # 环境变量配置
├── types/                   # 共享类型定义
│   └── bid_request.go       # 跨层数据类型
├── build/                   # 编译输出
│   └── dsp-system           # 可执行文件
├── main.go                  # 入口（初始化Gin+注册路由）
├── Makefile                 # 构建脚本
├── Dockerfile               # Docker 镜像
├── README.md                # 项目说明
├── EXAMPLE.md               # 使用示例
├── PROJECT_SUMMARY.md       # 项目摘要
└── test_request.json        # 测试请求数据
```

## 功能特性

- ✅ OpenRTB 2.5 协议支持
- ✅ 用户画像服务集成（gRPC）
- ✅ 预算服务集成（gRPC）
- ✅ Redis 缓存层
- ✅ ClickHouse 日志存储
- ✅ 广告智能匹配算法
- ✅ 高性能竞价处理（100ms 超时要求）

## 快速开始

### 1. 安装依赖

```bash
go mod tidy
```

### 2. 配置环境变量

```bash
# 复制环境变量模板
cp .env.example .env

# 编辑配置
vim .env
```

### 3. 启动服务

```bash
# 开发模式
make run

# 或直接运行
go run main.go
```

### 4. 测试接口

```bash
# 健康检查
curl http://localhost:8080/health

# 发送竞价请求
curl -X POST http://localhost:8080/bid \
  -H "Content-Type: application/json" \
  -d @test_request.json
```

## API 接口

### 1. 竞价接口

**POST /bid**

接收 ADX 发来的 OpenRTB 竞价请求，返回竞价响应。

请求示例：

```json
{
  "id": "req_12345",
  "imp": [
    {
      "id": "imp_001",
      "banner": {
        "w": 728,
        "h": 90
      },
      "bidfloor": 2.5
    }
  ],
  "user": {
    "id": "user_67890"
  },
  "device": {
    "ua": "Mozilla/5.0...",
    "ip": "192.168.1.1",
    "devicetype": 1,
    "os": "iOS"
  }
}
```

响应示例：

```json
{
  "id": "req_12345",
  "cur": "CNY",
  "seatbid": [
    {
      "bid": [
        {
          "id": "bid_001",
          "impid": "imp_001",
          "price": 5.5,
          "adid": "ad_001",
          "cid": "campaign_001",
          "crid": "creative_001",
          "w": 728,
          "h": 90
        }
      ]
    }
  ]
}
```

### 2. 健康检查

**GET /health**

返回服务健康状态。

### 3. 统计信息

**GET /stats**

返回 DSP 系统统计信息（QPS、成功率等）。

### 4. 赢标通知

**GET /win?price=${AUCTION_PRICE}&bidid=xxx**

接收 ADX 发来的赢标通知。

### 5. 计费通知

**GET /bill?bidid=xxx**

接收 ADX 发来的计费通知。

## 技术栈

- **Web 框架**: Gin
- **RPC 框架**: gRPC
- **缓存**: Redis
- **日志存储**: ClickHouse
- **配置管理**: 环境变量

## 核心流程

### 竞价流程

1. **接收请求**: ADX 发送 OpenRTB 竞价请求
2. **解析请求**: 解析广告位、设备、用户信息
3. **获取画像**: 通过 gRPC 调用用户画像服务
4. **广告匹配**: 根据用户标签匹配候选广告
5. **预算校验**: 通过 gRPC 调用预算服务检查预算
6. **出价计算**: 基于 eCPM 算法计算出价
7. **返回响应**: 构建 OpenRTB 响应返回给 ADX
8. **记录日志**: 异步记录竞价日志到 ClickHouse

### 赢标流程

1. **接收通知**: ADX 发送赢标通知（含最终价格）
2. **扣减预算**: 调用预算服务扣减预算
3. **记录日志**: 更新竞价日志状态为"赢标"
4. **返回确认**: 返回 200 OK

## 开发指南

### 添加新的广告匹配规则

编辑 `service/ad_select.go`，在 `SelectAds` 方法中添加匹配逻辑。

### 添加新的 gRPC 客户端

在 `rpc/` 目录下创建新的客户端文件，参考 `user_client.go` 的实现。

### 添加新的缓存操作

在 `repository/redis_cache.go` 中添加新的方法。

### 添加新的日志类型

在 `repository/clickhouse_repo.go` 中添加新的日志结构和方法。

## 性能优化建议

1. **并行调用**: 用户画像和预算检查可以并行调用
2. **缓存优化**: 频繁访问的数据（用户画像、广告素材）使用 Redis 缓存
3. **连接池**: gRPC 连接使用连接池管理
4. **批量操作**: 多个广告位的预算检查可以批量调用
5. **超时控制**: 严格控制每个步骤的超时时间（总共不超过 100ms）

## 监控指标

- QPS（每秒请求数）
- 竞价成功率
- 赢标率
- 平均响应时间
- 预算使用率
- 缓存命中率

## 注意事项

⚠️ **本项目为示例代码**，生产环境使用需要补充：

1. **gRPC 接口定义**: 需要编写 `.proto` 文件并生成代码
2. **数据库连接**: ClickHouse 和 Redis 的实际连接代码
3. **错误处理**: 更完善的错误处理和重试机制
4. **日志系统**: 结构化日志和日志收集
5. **监控告警**: Prometheus + Grafana 监控
6. **单元测试**: 各模块的单元测试
7. **压力测试**: 验证系统性能指标
8. **安全机制**: API 认证、限流、防刷等

## License

MIT

