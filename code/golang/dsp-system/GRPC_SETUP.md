# gRPC 服务集成文档

## 🎉 完成功能

已成功实现并集成 protobuf + gRPC 服务，DSP 系统现在使用真实的微服务架构！

### ✅ 已完成

1. **Protobuf 定义**
   - `proto/budget.proto` - 预算服务定义
   - `proto/user.proto` - 用户服务定义

2. **gRPC 服务实现**
   - `grpc_server/user_service.go` - 用户画像服务
   - `grpc_server/budget_service.go` - 预算管理服务

3. **客户端集成**
   - 更新 `rpc/user_client.go` 使用真实 gRPC 调用
   - 更新 `rpc/budget_client.go` 使用真实 gRPC 调用
   - 支持服务降级（gRPC 不可用时使用模拟数据）

4. **启动脚本**
   - `scripts/start_user_service.sh` - 启动用户服务
   - `scripts/start_budget_service.sh` - 启动预算服务
   - `scripts/start_all.sh` - 一键启动所有服务
   - `scripts/stop_all.sh` - 一键停止所有服务

## 🏗️ 架构概览

```
┌─────────────────┐
│   HTTP Client   │
└────────┬────────┘
         │ HTTP/JSON
         ▼
┌─────────────────────────┐
│  DSP HTTP Service       │
│  (port: 8088)          │
└──┬──────────────────┬───┘
   │                  │
   │ gRPC            │ gRPC
   ▼                  ▼
┌──────────────┐  ┌──────────────┐
│ User Service │  │Budget Service│
│ (port:50051) │  │ (port:50052) │
└──────────────┘  └──────────────┘
```

## 📦 服务列表

| 服务 | 端口 | 协议 | 说明 |
|------|------|------|------|
| User Service | 50051 | gRPC | 用户画像服务 |
| Budget Service | 50052 | gRPC | 预算管理服务 |
| DSP HTTP Service | 8088 | HTTP | RTB 竞价服务 |

## 🚀 快速开始

### 方式1：一键启动所有服务

```bash
cd /Users/shengli/Code/www/blog/code/golang/dsp-system

# 启动所有服务
./scripts/start_all.sh

# 停止所有服务
./scripts/stop_all.sh
```

### 方式2：手动逐个启动

```bash
# 1. 启动 User Service
./scripts/start_user_service.sh &

# 2. 启动 Budget Service  
./scripts/start_budget_service.sh &

# 等待 gRPC 服务启动
sleep 2

# 3. 启动 DSP HTTP Service
go run main.go
```

## 📝 Protobuf 服务定义

### User Service (用户画像服务)

```protobuf
service UserService {
  // 获取用户画像
  rpc GetUserProfile(GetUserProfileRequest) returns (GetUserProfileResponse);
  
  // 更新用户行为
  rpc UpdateUserBehavior(UpdateUserBehaviorRequest) returns (UpdateUserBehaviorResponse);
  
  // 批量获取用户画像
  rpc BatchGetUserProfiles(BatchGetUserProfilesRequest) returns (BatchGetUserProfilesResponse);
}
```

**测试数据**：
- `user_001` - 男性，28岁，运动爱好者
- `user_002` - 女性，22岁，购物达人
- `user_003` - 男性，38岁，商务人士
- `user_12345` - 男性，30岁，程序员

### Budget Service (预算管理服务)

```protobuf
service BudgetService {
  // 检查预算
  rpc CheckBudget(CheckBudgetRequest) returns (CheckBudgetResponse);
  
  // 扣减预算
  rpc DeductBudget(DeductBudgetRequest) returns (DeductBudgetResponse);
  
  // 获取预算信息
  rpc GetBudgetInfo(GetBudgetInfoRequest) returns (GetBudgetInfoResponse);
  
  // 退还预算
  rpc RefundBudget(RefundBudgetRequest) returns (RefundBudgetResponse);
}
```

**测试数据**：
- `campaign_001` - 总预算 10000，剩余 8500
- `campaign_002` - 总预算 5000，剩余 3200
- `campaign_003` - 总预算 20000，剩余 18500

## 🧪 测试

### 1. 健康检查

```bash
curl http://localhost:8088/health
```

### 2. 发送竞价请求（会调用真实 gRPC 服务）

```bash
curl -X POST http://localhost:8088/bid \
  -H "Content-Type: application/json" \
  -d '{
    "id": "test_request_001",
    "imp": [{
      "id": "1",
      "banner": {"w": 300, "h": 250},
      "bidfloor": 3.0
    }],
    "user": {"id": "user_12345"},
    "device": {
      "ua": "Mozilla/5.0...",
      "ip": "192.168.1.1"
    }
  }'
```

### 3. 使用测试脚本

```bash
./test_grpc_bid.sh
```

## 📊 查看日志

### 实时查看日志

```bash
# DSP HTTP 服务日志
tail -f logs/dsp-system.log

# User Service 日志
tail -f logs/user_service.log

# Budget Service 日志
tail -f logs/budget_service.log
```

### 查看 gRPC 调用日志

```bash
# 查看用户画像查询
grep "获取用户画像" logs/user_service.log

# 查看预算检查
grep "检查预算" logs/budget_service.log

# 查看预算扣减
grep "扣减预算" logs/budget_service.log
```

## 🔧 开发指南

### 修改 Protobuf 定义

1. 编辑 proto 文件
```bash
vim proto/budget.proto
vim proto/user.proto
```

2. 重新生成代码
```bash
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       proto/budget.proto proto/user.proto
```

3. 重新编译服务
```bash
go build -o build/user_service grpc_server/user_service.go
go build -o build/budget_service grpc_server/budget_service.go
```

### 添加新的 gRPC 方法

1. 在 `.proto` 文件中定义新方法
2. 重新生成代码
3. 在服务实现中添加方法实现
4. 在客户端中调用新方法

## 🛡️ 服务降级

客户端已实现服务降级逻辑：

```go
// 如果 gRPC 服务不可用，自动使用模拟数据
if c.client == nil {
    // 返回模拟数据
    return &UserProfile{
        UserID: userID,
        Tags:   []string{"默认标签"},
    }, nil
}

// 正常调用 gRPC
resp, err := c.client.GetUserProfile(ctx, req)
```

### 测试降级逻辑

```bash
# 停止 User Service
pkill -f user_service

# 发送竞价请求（会使用降级数据）
curl -X POST http://localhost:8088/bid \
  -H "Content-Type: application/json" \
  -d '{"id":"test","imp":[{"id":"1","banner":{"w":300,"h":250}}],"user":{"id":"user_12345"}}'
```

## 📈 性能监控

### 查看服务状态

```bash
# 查看进程
ps aux | grep -E "(user_service|budget_service|dsp-system)"

# 查看端口占用
lsof -i:50051  # User Service
lsof -i:50052  # Budget Service
lsof -i:8088   # DSP HTTP
```

### gRPC 调用统计

```bash
# 统计用户画像查询次数
grep -c "获取用户画像" logs/user_service.log

# 统计预算检查次数
grep -c "检查预算" logs/budget_service.log

# 统计竞价请求次数
grep -c "POST.*bid" logs/dsp-system.log
```

## 🐛 故障排查

### 1. gRPC 服务无法启动

**问题**：端口被占用
```bash
# 查找占用端口的进程
lsof -i:50051
lsof -i:50052

# 杀死进程
kill -9 <PID>
```

### 2. gRPC 连接失败

**问题**：`connect: connection refused`

**解决方案**：
1. 确认 gRPC 服务已启动
2. 检查服务端口是否正确
3. 查看服务日志

```bash
tail -f logs/user_service.log
tail -f logs/budget_service.log
```

### 3. 编译错误

**问题**：`found packages budget and user in proto`

**解决方案**：
```bash
# 清理旧的 pb 文件
rm -f proto/*.pb.go

# 重新生成
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       proto/budget.proto proto/user.proto
```

## 📚 相关文档

- [Protobuf 官方文档](https://protobuf.dev/)
- [gRPC Go 文档](https://grpc.io/docs/languages/go/)
- [项目总体文档](PROJECT_SUMMARY.md)
- [日志系统文档](LOG_SETUP.md)

## ✅ 验证清单

- [x] User Service 编译成功
- [x] Budget Service 编译成功
- [x] gRPC 服务能正常启动
- [x] DSP HTTP Service 能连接到 gRPC 服务
- [x] 竞价请求能正常处理
- [x] gRPC 调用日志记录正常
- [x] 服务降级功能正常
- [x] 启动/停止脚本正常工作

## 🎯 下一步

1. **数据库集成**：将内存存储替换为 Redis/MySQL
2. **服务发现**：集成 Consul/Etcd 实现服务注册与发现
3. **负载均衡**：实现 gRPC 客户端负载均衡
4. **链路追踪**：集成 OpenTelemetry 实现分布式追踪
5. **监控告警**：集成 Prometheus + Grafana
6. **熔断限流**：实现 Circuit Breaker 和 Rate Limiter

## 🎊 总结

✅ **已完成的功能**：
- protobuf 定义完整
- gRPC 服务端实现完整
- gRPC 客户端集成完整
- 服务降级逻辑完善
- 启动脚本齐全
- 测试通过

🚀 **系统现在是真正的微服务架构，支持水平扩展和服务解耦！**


