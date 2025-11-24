# DSP System 使用示例

## 快速演示

### 1. 启动服务

```bash
# 进入项目目录
cd dsp-system

# 运行服务
go run main.go
```

输出示例：
```
配置加载完成
2024/01/15 10:00:00 Redis连接失败: dial tcp [::1]:6379: connect: connection refused
2024/01/15 10:00:00 ClickHouse配置加载: localhost:9000/dsp_logs
2024/01/15 10:00:00 RPC客户端初始化完成
2024/01/15 10:00:00 DSP服务启动: http://localhost:8080
======================================
API文档:
  POST /bid        - OpenRTB竞价接口
  GET  /health     - 健康检查
  GET  /stats      - 统计信息
  GET  /win        - 赢标通知
  GET  /bill       - 计费通知
======================================
```

### 2. 测试健康检查

```bash
curl http://localhost:8080/health
```

响应：
```json
{
  "status": "ok",
  "time": 1705286400
}
```

### 3. 发送竞价请求

使用提供的测试数据：

```bash
curl -X POST http://localhost:8080/bid \
  -H "Content-Type: application/json" \
  -d @test_request.json
```

响应示例：
```json
{
  "id": "req_test_20240101_001",
  "cur": "CNY",
  "seatbid": [
    {
      "bid": [
        {
          "id": "bid_req_test_20240101_001_1705286400123456000",
          "impid": "imp_001",
          "price": 6.2,
          "adid": "ad_003",
          "nurl": "http://dsp.example.com/win?price=${AUCTION_PRICE}",
          "burl": "http://dsp.example.com/bill?bidid=ad_003",
          "cid": "campaign_003",
          "crid": "creative_003",
          "adomain": ["tech.com"],
          "w": 728,
          "h": 90
        }
      ],
      "seat": "dsp-seat-001"
    }
  ]
}
```

服务器日志：
```
收到竞价请求: ID=req_test_20240101_001, Imps=1
调用用户服务: UserID=user_test_67890
用户画像: UserID=user_test_67890, Tags=[男性 25-34岁 运动爱好者 科技爱好者]
候选广告数量: 3
标签匹配: UserTags=[男性 25-34岁 运动爱好者 科技爱好者], Matched=3/3
广告选择完成: Total=3
检查预算: CampaignID=campaign_003, BidPrice=6.20
检查预算: CampaignID=campaign_001, BidPrice=5.50
检查预算: CampaignID=campaign_002, BidPrice=4.80
BidLog: {Timestamp:2024-01-15 10:00:00 RequestID:req_test_20240101_001 ImpID:imp_001 UserID:user_test_67890 ...}
竞价完成: Bids=3, Duration=23ms
竞价成功: RequestID=req_test_20240101_001, Bids=3, Time=23ms
```

### 4. 模拟赢标通知

```bash
curl "http://localhost:8080/win?price=6.20&bidid=ad_003"
```

响应：
```
OK
```

服务器日志：
```
赢标通知: BidID=ad_003, Price=6.20
```

### 5. 查看统计信息

```bash
curl http://localhost:8080/stats
```

响应：
```json
{
  "avg_response_ms": 0,
  "bid_rate": 0,
  "qps": 0,
  "total_bids": 0,
  "total_requests": 0,
  "total_wins": 0,
  "win_rate": 0
}
```

## 自定义竞价请求

### 示例 1: Banner 广告

```json
{
  "id": "banner_001",
  "imp": [
    {
      "id": "imp_001",
      "banner": {
        "w": 300,
        "h": 250,
        "pos": 0
      },
      "bidfloor": 1.5
    }
  ],
  "user": {
    "id": "user_001",
    "gender": "F",
    "yob": 1998
  },
  "device": {
    "ip": "192.168.1.1",
    "ua": "Mozilla/5.0...",
    "devicetype": 4,
    "os": "iOS"
  }
}
```

### 示例 2: 移动端广告

```json
{
  "id": "mobile_001",
  "imp": [
    {
      "id": "imp_mobile_001",
      "banner": {
        "w": 320,
        "h": 50
      },
      "bidfloor": 2.0
    }
  ],
  "app": {
    "id": "app_001",
    "name": "购物APP",
    "bundle": "com.example.shop"
  },
  "device": {
    "ifa": "AEBE52E7-03EE-455A-B3C4-E57283966239",
    "devicetype": 4,
    "os": "iOS",
    "osv": "17.0"
  },
  "user": {
    "id": "user_mobile_001"
  }
}
```

### 示例 3: 多广告位请求

```json
{
  "id": "multi_imp_001",
  "imp": [
    {
      "id": "imp_top",
      "banner": {
        "w": 728,
        "h": 90
      },
      "bidfloor": 5.0,
      "tagid": "top_banner"
    },
    {
      "id": "imp_sidebar",
      "banner": {
        "w": 300,
        "h": 600
      },
      "bidfloor": 3.0,
      "tagid": "sidebar_banner"
    },
    {
      "id": "imp_footer",
      "banner": {
        "w": 728,
        "h": 90
      },
      "bidfloor": 2.0,
      "tagid": "footer_banner"
    }
  ],
  "user": {
    "id": "user_multi_001"
  }
}
```

## 使用 Makefile

项目提供了便捷的 Makefile 命令：

```bash
# 查看所有可用命令
make help

# 安装依赖
make install

# 运行服务
make run

# 编译项目
make build

# 运行测试
make test

# 代码格式化
make fmt

# 代码检查
make lint

# 发送测试请求
make test-request

# 健康检查
make health

# 查看统计
make stats
```

## Docker 部署

### 构建镜像

```bash
make docker-build
# 或
docker build -t dsp-system:latest .
```

### 运行容器

创建 `.env` 文件：
```bash
SERVER_PORT=8080
REDIS_HOST=redis
REDIS_PORT=6379
CLICKHOUSE_HOST=clickhouse
CLICKHOUSE_PORT=9000
USER_SERVICE_ADDR=user-service:50051
BUDGET_SERVICE_ADDR=budget-service:50052
```

运行容器：
```bash
docker run -p 8080:8080 --env-file .env dsp-system:latest
```

### 使用 docker-compose

创建 `docker-compose.yml`：

```yaml
version: '3.8'

services:
  dsp-system:
    build: .
    ports:
      - "8080:8080"
    environment:
      - SERVER_PORT=8080
      - REDIS_HOST=redis
      - REDIS_PORT=6379
      - CLICKHOUSE_HOST=clickhouse
      - CLICKHOUSE_PORT=9000
      - USER_SERVICE_ADDR=user-service:50051
      - BUDGET_SERVICE_ADDR=budget-service:50052
    depends_on:
      - redis
      - clickhouse

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"

  clickhouse:
    image: clickhouse/clickhouse-server:latest
    ports:
      - "8123:8123"
      - "9000:9000"
    environment:
      - CLICKHOUSE_DB=dsp_logs
```

启动：
```bash
docker-compose up -d
```

## 性能测试

使用 Apache Bench 进行压力测试：

```bash
# 安装 ab
# macOS: brew install httpd
# Ubuntu: sudo apt-get install apache2-utils

# 测试竞价接口（100 并发，共 1000 请求）
ab -n 1000 -c 100 -p test_request.json -T application/json \
  http://localhost:8080/bid
```

使用 wrk 进行压力测试：

```bash
# 安装 wrk
# macOS: brew install wrk
# Ubuntu: 从源码编译

# 创建 lua 脚本 (post.lua)
wrk.method = "POST"
wrk.body   = io.open("test_request.json"):read("*a")
wrk.headers["Content-Type"] = "application/json"

# 运行测试（12 线程，400 连接，持续 30 秒）
wrk -t12 -c400 -d30s -s post.lua http://localhost:8080/bid
```

## 监控和调试

### 查看实时日志

```bash
tail -f /path/to/logfile
```

### 使用 pprof 分析性能

在 `main.go` 中添加：

```go
import _ "net/http/pprof"

// 在 main 函数中
go func() {
    log.Println(http.ListenAndServe("localhost:6060", nil))
}()
```

访问性能分析：
```bash
# CPU 分析
go tool pprof http://localhost:6060/debug/pprof/profile

# 内存分析
go tool pprof http://localhost:6060/debug/pprof/heap

# 协程分析
go tool pprof http://localhost:6060/debug/pprof/goroutine
```

## 常见问题

### Q1: Redis 连接失败

确保 Redis 已启动：
```bash
# 启动 Redis
redis-server

# 或使用 Docker
docker run -d -p 6379:6379 redis:7-alpine
```

### Q2: gRPC 服务连接失败

本示例中的 gRPC 客户端使用模拟数据。实际项目中需要：
1. 编写 `.proto` 文件定义接口
2. 使用 `protoc` 生成代码
3. 启动实际的 gRPC 服务

### Q3: 响应时间过长

优化建议：
1. 使用 Redis 缓存用户画像
2. 并行调用多个 RPC 服务
3. 减少日志输出
4. 使用连接池
5. 优化广告匹配算法

## 扩展开发

### 添加新的竞价策略

编辑 `service/ad_select.go`，修改 `calculateAdScore` 方法。

### 添加频次控制

使用 Redis 实现：

```go
// 检查频次
if !s.redisCache.CheckFrequencyCap(ctx, userID, adID) {
    continue // 跳过该广告
}

// 设置频次（24小时内不再展示）
s.redisCache.SetFrequencyCap(ctx, userID, adID, 24*time.Hour)
```

### 添加实时竞价日志

使用 ClickHouse 存储：

```sql
CREATE TABLE bid_logs (
    timestamp DateTime,
    request_id String,
    imp_id String,
    user_id String,
    ad_id String,
    bid_price Float64,
    bid_status Enum8('bid'=1, 'win'=2, 'lose'=3),
    processing_time UInt32
) ENGINE = MergeTree()
ORDER BY (timestamp, request_id);
```

## 参考资料

- [OpenRTB 2.5 规范](https://www.iab.com/wp-content/uploads/2016/03/OpenRTB-API-Specification-Version-2-5-FINAL.pdf)
- [Gin 文档](https://gin-gonic.com/docs/)
- [gRPC Go 教程](https://grpc.io/docs/languages/go/)
- [Redis Go 客户端](https://redis.uptrace.dev/)
- [ClickHouse Go 客户端](https://clickhouse.com/docs/en/integrations/go)

