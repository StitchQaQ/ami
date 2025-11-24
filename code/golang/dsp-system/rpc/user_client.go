package rpc

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "dsp-system/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// UserProfile 用户画像
type UserProfile struct {
	UserID   string
	Age      int
	Gender   string
	Tags     []string
	Interests []string
}

// UserClient 用户服务客户端
type UserClient struct {
	addr   string
	conn   *grpc.ClientConn
	client pb.UserServiceClient
}

// NewUserClient 创建用户服务客户端
func NewUserClient(addr string) *UserClient {
	return &UserClient{
		addr: addr,
	}
}

// Connect 连接到gRPC服务
func (c *UserClient) Connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		c.addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return fmt.Errorf("连接用户服务失败: %v", err)
	}

	c.conn = conn
	c.client = pb.NewUserServiceClient(conn)
	
	log.Printf("用户服务连接成功: %s", c.addr)
	return nil
}

// Close 关闭连接
func (c *UserClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// GetUserProfile 获取用户画像
func (c *UserClient) GetUserProfile(ctx context.Context, userID string) (*UserProfile, error) {
	if c.client == nil {
		// 降级逻辑：gRPC 未连接时使用模拟数据
		log.Printf("用户服务未连接，使用模拟数据: UserID=%s", userID)
		return &UserProfile{
			UserID:    userID,
			Age:       28,
			Gender:    "male",
			Tags:      []string{"男性", "25-34岁", "运动爱好者", "科技爱好者"},
			Interests: []string{"篮球", "编程", "旅游"},
		}, nil
	}

	req := &pb.GetUserProfileRequest{
		UserId: userID,
	}
	
	resp, err := c.client.GetUserProfile(ctx, req)
	if err != nil {
		log.Printf("获取用户画像失败: %v", err)
		return nil, err
	}
	
	profile := &UserProfile{
		UserID:    resp.UserId,
		Age:       int(resp.Age),
		Gender:    resp.Gender,
		Tags:      resp.Tags,
		Interests: resp.Interests,
	}
	
	log.Printf("获取用户画像成功: UserID=%s, Tags=%v", userID, profile.Tags)
	return profile, nil
}

// BatchGetUserProfiles 批量获取用户画像
func (c *UserClient) BatchGetUserProfiles(ctx context.Context, userIDs []string) (map[string]*UserProfile, error) {
	if c.client == nil {
		// 降级逻辑：逐个调用
		log.Printf("用户服务未连接，使用单个调用: Count=%d", len(userIDs))
		profiles := make(map[string]*UserProfile)
		for _, userID := range userIDs {
			profile, err := c.GetUserProfile(ctx, userID)
			if err != nil {
				continue
			}
			profiles[userID] = profile
		}
		return profiles, nil
	}

	req := &pb.BatchGetUserProfilesRequest{
		UserIds: userIDs,
	}
	
	resp, err := c.client.BatchGetUserProfiles(ctx, req)
	if err != nil {
		log.Printf("批量获取用户画像失败: %v", err)
		return nil, err
	}
	
	profiles := make(map[string]*UserProfile)
	for _, p := range resp.Profiles {
		profiles[p.UserId] = &UserProfile{
			UserID:    p.UserId,
			Age:       int(p.Age),
			Gender:    p.Gender,
			Tags:      p.Tags,
			Interests: p.Interests,
		}
	}
	
	log.Printf("批量获取用户画像成功: Count=%d", len(profiles))
	return profiles, nil
}

// UpdateUserBehavior 更新用户行为（点击、转化等）
func (c *UserClient) UpdateUserBehavior(ctx context.Context, userID string, behavior string, adID string) error {
	if c.client == nil {
		log.Printf("用户服务未连接，跳过行为更新: UserID=%s, Behavior=%s", userID, behavior)
		return nil
	}

	req := &pb.UpdateUserBehaviorRequest{
		UserId:    userID,
		Behavior:  behavior,
		AdId:      adID,
		Timestamp: time.Now().Unix(),
	}
	
	resp, err := c.client.UpdateUserBehavior(ctx, req)
	if err != nil {
		log.Printf("更新用户行为失败: %v", err)
		return err
	}
	
	if !resp.Success {
		return fmt.Errorf("更新失败: %s", resp.Message)
	}
	
	log.Printf("更新用户行为成功: UserID=%s, Behavior=%s, AdID=%s", userID, behavior, adID)
	return nil
}

