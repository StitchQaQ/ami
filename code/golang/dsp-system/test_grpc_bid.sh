#!/bin/bash

# 测试 DSP 系统（使用真实 gRPC）

echo "========================================="
echo "测试 DSP 系统 (使用真实 gRPC)"
echo "========================================="
echo ""

# 1. 健康检查
echo "1. 健康检查"
curl -s http://localhost:8088/health
echo ""
echo ""

# 2. 发送竞价请求（会调用真实的 gRPC 服务）
echo "2. 发送竞价请求"
curl -X POST http://localhost:8088/bid \
  -H "Content-Type: application/json" \
  -d '{
    "id": "test_request_001",
    "imp": [
      {
        "id": "1",
        "banner": {
          "w": 300,
          "h": 250
        },
        "bidfloor": 3.0
      }
    ],
    "user": {
      "id": "user_12345"
    },
    "device": {
      "ua": "Mozilla/5.0...",
      "ip": "192.168.1.1"
    }
  }' | python3 -m json.tool 2>/dev/null || cat

echo ""
echo ""

# 3. 查看日志
echo "3. 最近的日志记录："
echo ""
echo "=== DSP HTTP 日志 ==="
tail -5 logs/dsp-http.log
echo ""
echo "=== User Service 日志 ==="
tail -5 logs/user_service.log
echo ""
echo "=== Budget Service 日志 ==="
tail -5 logs/budget_service.log

echo ""
echo "========================================="
echo "✅ 测试完成"
echo "========================================="



