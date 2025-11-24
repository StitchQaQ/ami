package main

import (
	"context"
	"log"
	"net"
	"sync"
	"time"

	pb "dsp-system/proto"

	"google.golang.org/grpc"
)

// UserServer 用户服务实现
type UserServer struct {
	pb.UnimplementedUserServiceServer
	
	// 内存存储（实际项目应使用数据库）
	profiles map[string]*UserProfile
	mu       sync.RWMutex
}

// UserProfile 用户画像
type UserProfile struct {
	UserID     string
	Tags       []string
	Age        int32
	Gender     string
	Interests  []string
	City       string
	DeviceType string
}

// NewUserServer 创建用户服务
func NewUserServer() *UserServer {
	s := &UserServer{
		profiles: make(map[string]*UserProfile),
	}
	
	// 初始化一些测试数据
	s.initTestData()
	
	return s
}

// initTestData 初始化测试数据
func (s *UserServer) initTestData() {
	testProfiles := []UserProfile{
		{
			UserID:     "user_001",
			Tags:       []string{"男性", "25-34岁", "运动爱好者", "科技爱好者"},
			Age:        28,
			Gender:     "male",
			Interests:  []string{"sports", "technology", "gaming"},
			City:       "北京",
			DeviceType: "ios",
		},
		{
			UserID:     "user_002",
			Tags:       []string{"女性", "18-24岁", "购物达人", "美妆爱好者"},
			Age:        22,
			Gender:     "female",
			Interests:  []string{"shopping", "beauty", "fashion"},
			City:       "上海",
			DeviceType: "android",
		},
		{
			UserID:     "user_003",
			Tags:       []string{"男性", "35-44岁", "商务人士", "汽车爱好者"},
			Age:        38,
			Gender:     "male",
			Interests:  []string{"business", "cars", "finance"},
			City:       "深圳",
			DeviceType: "ios",
		},
		{
			UserID:     "user_12345",
			Tags:       []string{"科技爱好者", "程序员", "男性"},
			Age:        30,
			Gender:     "male",
			Interests:  []string{"technology", "programming", "reading"},
			City:       "杭州",
			DeviceType: "android",
		},
	}
	
	for _, profile := range testProfiles {
		s.profiles[profile.UserID] = &profile
	}
	
	log.Printf("初始化 %d 个测试用户画像", len(testProfiles))
}

// GetUserProfile 获取用户画像
func (s *UserServer) GetUserProfile(ctx context.Context, req *pb.GetUserProfileRequest) (*pb.GetUserProfileResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	log.Printf("获取用户画像: UserID=%s", req.UserId)
	
	profile, exists := s.profiles[req.UserId]
	if !exists {
		// 返回默认画像
		log.Printf("用户不存在，返回默认画像: UserID=%s", req.UserId)
		return &pb.GetUserProfileResponse{
			UserId:     req.UserId,
			Tags:       []string{"新用户"},
			Age:        0,
			Gender:     "unknown",
			Interests:  []string{},
			City:       "未知",
			DeviceType: "unknown",
		}, nil
	}
	
	return &pb.GetUserProfileResponse{
		UserId:     profile.UserID,
		Tags:       profile.Tags,
		Age:        profile.Age,
		Gender:     profile.Gender,
		Interests:  profile.Interests,
		City:       profile.City,
		DeviceType: profile.DeviceType,
	}, nil
}

// UpdateUserBehavior 更新用户行为
func (s *UserServer) UpdateUserBehavior(ctx context.Context, req *pb.UpdateUserBehaviorRequest) (*pb.UpdateUserBehaviorResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	log.Printf("更新用户行为: UserID=%s, Behavior=%s, AdID=%s", req.UserId, req.Behavior, req.AdId)
	
	profile, exists := s.profiles[req.UserId]
	if !exists {
		// 创建新用户画像
		profile = &UserProfile{
			UserID:     req.UserId,
			Tags:       []string{"新用户"},
			Interests:  []string{},
			DeviceType: "unknown",
		}
		s.profiles[req.UserId] = profile
		log.Printf("创建新用户画像: UserID=%s", req.UserId)
	}
	
	// 根据行为更新用户画像（简化逻辑）
	switch req.Behavior {
	case "click":
		// 点击行为可能增加相关兴趣标签
		log.Printf("用户点击广告: UserID=%s, AdID=%s", req.UserId, req.AdId)
	case "conversion":
		// 转化行为说明用户对该类型广告很感兴趣
		log.Printf("用户转化: UserID=%s, AdID=%s", req.UserId, req.AdId)
	case "view":
		// 浏览行为
		log.Printf("用户浏览广告: UserID=%s, AdID=%s", req.UserId, req.AdId)
	}
	
	return &pb.UpdateUserBehaviorResponse{
		Success: true,
		Message: "行为记录成功",
	}, nil
}

// BatchGetUserProfiles 批量获取用户画像
func (s *UserServer) BatchGetUserProfiles(ctx context.Context, req *pb.BatchGetUserProfilesRequest) (*pb.BatchGetUserProfilesResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	log.Printf("批量获取用户画像: Count=%d", len(req.UserIds))
	
	var profiles []*pb.GetUserProfileResponse
	
	for _, userID := range req.UserIds {
		profile, exists := s.profiles[userID]
		if !exists {
			// 返回默认画像
			profiles = append(profiles, &pb.GetUserProfileResponse{
				UserId:     userID,
				Tags:       []string{"新用户"},
				Age:        0,
				Gender:     "unknown",
				Interests:  []string{},
				City:       "未知",
				DeviceType: "unknown",
			})
			continue
		}
		
		profiles = append(profiles, &pb.GetUserProfileResponse{
			UserId:     profile.UserID,
			Tags:       profile.Tags,
			Age:        profile.Age,
			Gender:     profile.Gender,
			Interests:  profile.Interests,
			City:       profile.City,
			DeviceType: profile.DeviceType,
		})
	}
	
	return &pb.BatchGetUserProfilesResponse{
		Profiles: profiles,
	}, nil
}

func main() {
	// 创建 gRPC 服务器
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("监听失败: %v", err)
	}
	
	grpcServer := grpc.NewServer()
	userServer := NewUserServer()
	
	pb.RegisterUserServiceServer(grpcServer, userServer)
	
	log.Println("======================================")
	log.Println("User gRPC 服务启动成功")
	log.Println("监听地址: localhost:50051")
	log.Println("======================================")
	log.Printf("当前时间: %s", time.Now().Format("2006-01-02 15:04:05"))
	
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}

