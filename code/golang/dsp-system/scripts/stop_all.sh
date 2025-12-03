#!/bin/bash

# 停止所有服务的脚本

echo "========================================="
echo "停止 DSP 系统所有服务"
echo "========================================="

# 停止 User Service
if pgrep -f "user_service" > /dev/null; then
    echo "停止 User Service..."
    pkill -f "user_service"
    echo "✓ User Service 已停止"
else
    echo "User Service 未在运行"
fi

# 停止 Budget Service
if pgrep -f "budget_service" > /dev/null; then
    echo "停止 Budget Service..."
    pkill -f "budget_service"
    echo "✓ Budget Service 已停止"
else
    echo "Budget Service 未在运行"
fi

# 停止 DSP HTTP Service
if pgrep -f "dsp-system" > /dev/null || pgrep -f "go run main.go" > /dev/null; then
    echo "停止 DSP HTTP Service..."
    pkill -f "go run main.go"
    pkill -f "dsp-system"
    echo "✓ DSP HTTP Service 已停止"
else
    echo "DSP HTTP Service 未在运行"
fi

echo ""
echo "========================================="
echo "✅ 所有服务已停止"
echo "========================================="



