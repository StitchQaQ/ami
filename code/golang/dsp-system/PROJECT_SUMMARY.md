# DSP System 项目摘要

## 项目概述

这是一个基于 **Gin + gRPC** 架构的 DSP（需求方平台）广告竞价系统示例项目。该系统实现了 OpenRTB 2.5 协议，支持实时竞价（RTB）功能。

## 技术架构

### 分层架构

```
┌─────────────────────────────────────┐
│          HTTP/JSON API              │  (ADX 请求入口)
└───────────┬─────────────────────────┘
            │
┌───────────▼─────────────────────────┐
│       Handler 层 (Gin)              │  HTTP 请求处理
│    - rtb_handler.go                 │
└───────────┬─────────────────────────┘
            │
┌───────────▼─────────────────────────┐
│       Service 层 (业务逻辑)         │  竞价决策核心
│    - bid_service.go                 │
│    - ad_select.go                   │
└──────┬─────────────┬────────────────┘
       │             │
       │             │
┌──────▼──────┐  ┌──▼─────────────────┐
│   RPC 层    │  │   Repository 层    │
│ (gRPC)      │  │  (数据访问)        │
│             │  │                    │
│ - user      │  │ - Redis Cache      │
│ - budget    │  │ - ClickHouse       │
└─────────────┘  └────────────────────┘
```

### 目录结构

```
dsp-system/
├── api/                      # OpenRTB 协议定义
│   └── openrtb.go            # BidRequest, BidResponse 等数据结构
│
├── handler/                  # HTTP 请求处理层
│   └── rtb_handler.go        # 处理 ADX 的 RTB 竞价请求
│
├── service/                  # 业务逻辑层
│   ├── bid_service.go        # 竞价决策核心
│   └── ad_select.go          # 广告素材匹配算法
│
├── rpc/                      # gRPC 客户端
│   ├── user_client.go        # 用户画像服务客户端
│   └── budget_client.go      # 预算服务客户端
│
├── repository/               # 数据访问层
│   ├── redis_cache.go        # Redis 缓存操作
│   └── clickhouse_repo.go    # ClickHouse 日志存储
│
├── config/                   # 配置管理
│   └── config.go             # 环境变量配置加载
│
├── build/                    # 编译输出目录
│   └── dsp-system            # 可执行文件
│
├── main.go                   # 应用入口
├── go.mod                    # Go 模块定义
├── go.sum                    # 依赖版本锁定
├── Makefile                  # 构建和运行脚本
├── Dockerfile                # Docker 镜像构建
├── .gitignore                # Git 忽略规则
├── README.md                 # 项目文档
├── EXAMPLE.md                # 使用示例
└── test_request.json         # 测试请求数据
```

## 核心功能

### 1. OpenRTB 竞价接口

- **接口**: `POST /bid`
- **功能**: 接收 ADX 发送的竞价请求，返回竞价响应
- **超时要求**: 100ms 内完成处理
- **协议**: OpenRTB 2.5

### 2. 用户画像匹配

- 通过 gRPC 调用用户画像服务
- 获取用户标签（年龄、性别、兴趣等）
- Redis 缓存用户画像数据（5分钟过期）

### 3. 广告智能匹配

- 根据用户标签匹配广告
- 多维度评分算法：
  - 出价权重：40%
  - 标签匹配度：30%
  - 随机因子：10%

### 4. 预算控制

- 通过 gRPC 调用预算服务
- 实时校验活动预算是否充足
- 支持批量预算检查

### 5. 日志记录

- 异步记录竞价日志到 ClickHouse
- 支持赢标、曝光、点击、转化等日志
- 性能监控数据收集

## 关键技术特性

### 1. 高性能

- **并发处理**: 用户画像查询与广告匹配并行执行
- **连接池**: gRPC 长连接复用
- **缓存策略**: Redis 缓存热点数据
- **异步日志**: 非阻塞日志写入

### 2. 可扩展性

- **分层架构**: 各层职责清晰，易于扩展
- **接口设计**: 使用接口抽象，便于替换实现
- **配置化**: 通过环境变量灵活配置

### 3. 可观测性

- **结构化日志**: 记录关键操作日志
- **性能指标**: 响应时间、成功率等
- **调用链追踪**: 支持 Request ID 追踪

## 依赖服务

### 必需服务

| 服务 | 端口 | 说明 |
|------|------|------|
| Redis | 6379 | 缓存用户画像、预算信息 |
| ClickHouse | 9000 | 存储竞价日志 |
| 用户服务 (gRPC) | 50051 | 提供用户画像 |
| 预算服务 (gRPC) | 50052 | 预算校验和扣减 |

### 可选服务

- Prometheus: 指标收集
- Grafana: 可视化监控
- Jaeger: 分布式追踪

## 快速开始

### 1. 安装依赖

```bash
go mod tidy
```

### 2. 启动服务

```bash
# 使用 Makefile
make run

# 或直接运行
go run main.go
```

### 3. 测试接口

```bash
# 健康检查
curl http://localhost:8080/health

# 发送竞价请求
curl -X POST http://localhost:8080/bid \
  -H "Content-Type: application/json" \
  -d @test_request.json
```

## 性能指标

### 目标性能

- **QPS**: > 10,000 req/s
- **响应时间**: P99 < 50ms
- **竞价成功率**: > 95%
- **可用性**: 99.9%

### 优化建议

1. **使用连接池**: gRPC 和 Redis 连接池
2. **批量处理**: 批量查询用户画像和预算
3. **预加载**: 广告库预加载到内存
4. **限流**: API 限流保护后端服务
5. **降级**: 服务降级策略

## 部署方式

### 1. 二进制部署

```bash
make build
./build/dsp-system
```

### 2. Docker 部署

```bash
docker build -t dsp-system:latest .
docker run -p 8080:8080 --env-file .env dsp-system:latest
```

### 3. Kubernetes 部署

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: dsp-system
spec:
  replicas: 3
  template:
    spec:
      containers:
      - name: dsp-system
        image: dsp-system:latest
        ports:
        - containerPort: 8080
```

## 监控告警

### 关键指标

- **业务指标**:
  - 竞价 QPS
  - 竞价成功率
  - 赢标率
  - 平均出价

- **性能指标**:
  - 响应时间 (P50/P95/P99)
  - 错误率
  - 超时率

- **资源指标**:
  - CPU 使用率
  - 内存使用率
  - 协程数量

### 告警规则

- 竞价成功率 < 90%
- P99 响应时间 > 100ms
- 错误率 > 1%
- Redis 连接失败

## 开发规范

### 代码风格

- 使用 `gofmt` 格式化代码
- 使用 `golangci-lint` 代码检查
- 遵循 Go 标准项目布局

### 提交规范

- feat: 新功能
- fix: Bug 修复
- docs: 文档更新
- refactor: 代码重构
- perf: 性能优化
- test: 测试相关

## 后续优化方向

1. **算法优化**:
   - 引入机器学习模型预测 CTR
   - 动态出价策略
   - 频次控制优化

2. **性能优化**:
   - 引入本地缓存（内存）
   - 优化序列化性能
   - 减少内存分配

3. **功能扩展**:
   - 支持视频广告
   - 支持原生广告
   - A/B 测试框架

4. **运维优化**:
   - 自动化部署流程
   - 灰度发布
   - 故障自愈

## 联系方式

- 项目地址: `/Users/shengli/Code/www/blog/code/golang/dsp-system`
- 文档: 查看 README.md 和 EXAMPLE.md
- 测试: 查看 test_request.json

## 注意事项

⚠️ **本项目为示例代码**，用于学习和演示目的。生产环境使用需要补充：

1. 完整的 gRPC 接口定义（.proto 文件）
2. 实际的数据库连接和操作
3. 完善的错误处理和重试机制
4. 全面的单元测试和集成测试
5. 生产级别的监控和日志系统
6. 安全认证和限流保护
7. 完整的 CI/CD 流程

## 许可证

MIT License

