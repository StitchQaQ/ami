#!/bin/bash

# User gRPC Service 启动脚本

cd "$(dirname "$0")/.."

echo "========================================="
echo "启动 User gRPC Service"
echo "========================================="

# 编译
go build -o build/user_service grpc_server/user_service.go

# 运行
./build/user_service



