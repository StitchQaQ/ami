#!/bin/bash

# 启动所有服务的脚本

cd "$(dirname "$0")/.."

echo "========================================="
echo "启动 DSP 系统所有服务"
echo "========================================="

# 检查是否已有服务在运行
if pgrep -f "user_service" > /dev/null; then
    echo "⚠️  User Service 已在运行"
else
    echo "启动 User Service..."
    ./scripts/start_user_service.sh > logs/user_service.log 2>&1 &
    sleep 2
    echo "✓ User Service 启动成功 (端口: 50051)"
fi

if pgrep -f "budget_service" > /dev/null; then
    echo "⚠️  Budget Service 已在运行"
else
    echo "启动 Budget Service..."
    ./scripts/start_budget_service.sh > logs/budget_service.log 2>&1 &
    sleep 2
    echo "✓ Budget Service 启动成功 (端口: 50052)"
fi

echo ""
echo "等待 gRPC 服务启动..."
sleep 2

if pgrep -f "dsp-system" > /dev/null; then
    echo "⚠️  DSP HTTP Service 已在运行"
else
    echo "启动 DSP HTTP Service..."
    go run main.go > logs/dsp-http.log 2>&1 &
    sleep 2
    echo "✓ DSP HTTP Service 启动成功 (端口: 8088)"
fi

echo ""
echo "========================================="
echo "✅ 所有服务启动完成！"
echo "========================================="
echo ""
echo "服务列表:"
echo "  - User Service (gRPC):   localhost:50051"
echo "  - Budget Service (gRPC): localhost:50052"
echo "  - DSP Service (HTTP):    localhost:8088"
echo ""
echo "查看日志:"
echo "  tail -f logs/user_service.log"
echo "  tail -f logs/budget_service.log"
echo "  tail -f logs/dsp-system.log"
echo ""
echo "停止所有服务:"
echo "  ./scripts/stop_all.sh"
echo ""



