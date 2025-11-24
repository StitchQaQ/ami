package main

import (
	"context"
	"log"
	"net"

	pb "github.com/dsp-system/protobuf-demo/pb/proto"
	"google.golang.org/grpc"
)

type userServer struct {
	pb.UnimplementedUserServiceServer
}

func (s *userServer) GetUser(ctx context.Context, req *pb.GetUserProfileRequest) (*pb.GetUserProfileResponse, error) {
	return &pb.GetUserProfileResponse{
		UserId:   req.UserId,
		UserName: req.UserName,
		UserAge:  req.UserAge,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()

	// 启动服务器
	log.Printf("gRPC服务启动, 监听地址：%s", lis.Addr().String())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
