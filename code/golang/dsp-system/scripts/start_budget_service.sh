#!/bin/bash

# Budget gRPC Service 启动脚本

cd "$(dirname "$0")/.."

echo "========================================="
echo "启动 Budget gRPC Service"
echo "========================================="

# 编译
go build -o build/budget_service grpc_server/budget_service.go

# 运行
./build/budget_service



