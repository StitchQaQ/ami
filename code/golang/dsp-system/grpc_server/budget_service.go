package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	pb "dsp-system/proto"

	"google.golang.org/grpc"
)

// BudgetServer 预算服务实现
type BudgetServer struct {
	pb.UnimplementedBudgetServiceServer
	
	// 内存存储（实际项目应使用数据库）
	budgets map[string]*BudgetInfo
	mu      sync.RWMutex
}

// BudgetInfo 预算信息
type BudgetInfo struct {
	CampaignID      string
	TotalBudget     float64
	RemainingBudget float64
	DailyBudget     float64
	DailySpent      float64
	Status          string
}

// NewBudgetServer 创建预算服务
func NewBudgetServer() *BudgetServer {
	s := &BudgetServer{
		budgets: make(map[string]*BudgetInfo),
	}
	
	// 初始化一些测试数据
	s.initTestData()
	
	return s
}

// initTestData 初始化测试数据
func (s *BudgetServer) initTestData() {
	testCampaigns := []BudgetInfo{
		{
			CampaignID:      "campaign_001",
			TotalBudget:     10000.0,
			RemainingBudget: 8500.0,
			DailyBudget:     1000.0,
			DailySpent:      200.0,
			Status:          "active",
		},
		{
			CampaignID:      "campaign_002",
			TotalBudget:     5000.0,
			RemainingBudget: 3200.0,
			DailyBudget:     500.0,
			DailySpent:      150.0,
			Status:          "active",
		},
		{
			CampaignID:      "campaign_003",
			TotalBudget:     20000.0,
			RemainingBudget: 18500.0,
			DailyBudget:     2000.0,
			DailySpent:      500.0,
			Status:          "active",
		},
	}
	
	for _, campaign := range testCampaigns {
		s.budgets[campaign.CampaignID] = &campaign
	}
	
	log.Printf("初始化 %d 个测试活动预算", len(testCampaigns))
}

// CheckBudget 检查预算
func (s *BudgetServer) CheckBudget(ctx context.Context, req *pb.CheckBudgetRequest) (*pb.CheckBudgetResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	log.Printf("检查预算: CampaignID=%s, Amount=%.2f", req.CampaignId, req.Amount)
	
	budget, exists := s.budgets[req.CampaignId]
	if !exists {
		return &pb.CheckBudgetResponse{
			HasBudget: false,
			Remaining: 0,
			Message:   "活动不存在",
		}, nil
	}
	
	if budget.Status != "active" {
		return &pb.CheckBudgetResponse{
			HasBudget: false,
			Remaining: budget.RemainingBudget,
			Message:   fmt.Sprintf("活动状态异常: %s", budget.Status),
		}, nil
	}
	
	// 检查总预算和日预算
	hasBudget := budget.RemainingBudget >= req.Amount && 
	             (budget.DailyBudget - budget.DailySpent) >= req.Amount
	
	message := "预算充足"
	if !hasBudget {
		if budget.RemainingBudget < req.Amount {
			message = "总预算不足"
		} else {
			message = "日预算不足"
		}
	}
	
	return &pb.CheckBudgetResponse{
		HasBudget: hasBudget,
		Remaining: budget.RemainingBudget,
		Message:   message,
	}, nil
}

// DeductBudget 扣减预算
func (s *BudgetServer) DeductBudget(ctx context.Context, req *pb.DeductBudgetRequest) (*pb.DeductBudgetResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	log.Printf("扣减预算: CampaignID=%s, Amount=%.2f, BidID=%s", req.CampaignId, req.Amount, req.BidId)
	
	budget, exists := s.budgets[req.CampaignId]
	if !exists {
		return &pb.DeductBudgetResponse{
			Success:   false,
			Remaining: 0,
			Message:   "活动不存在",
		}, nil
	}
	
	if budget.RemainingBudget < req.Amount {
		return &pb.DeductBudgetResponse{
			Success:   false,
			Remaining: budget.RemainingBudget,
			Message:   "预算不足",
		}, nil
	}
	
	// 扣减预算
	budget.RemainingBudget -= req.Amount
	budget.DailySpent += req.Amount
	
	log.Printf("预算扣减成功: CampaignID=%s, Remaining=%.2f", req.CampaignId, budget.RemainingBudget)
	
	return &pb.DeductBudgetResponse{
		Success:   true,
		Remaining: budget.RemainingBudget,
		Message:   "扣减成功",
	}, nil
}

// GetBudgetInfo 获取预算信息
func (s *BudgetServer) GetBudgetInfo(ctx context.Context, req *pb.GetBudgetInfoRequest) (*pb.GetBudgetInfoResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	log.Printf("获取预算信息: CampaignID=%s", req.CampaignId)
	
	budget, exists := s.budgets[req.CampaignId]
	if !exists {
		return nil, fmt.Errorf("活动不存在: %s", req.CampaignId)
	}
	
	return &pb.GetBudgetInfoResponse{
		CampaignId:      budget.CampaignID,
		TotalBudget:     budget.TotalBudget,
		RemainingBudget: budget.RemainingBudget,
		DailyBudget:     budget.DailyBudget,
		DailySpent:      budget.DailySpent,
		Status:          budget.Status,
	}, nil
}

// RefundBudget 退还预算
func (s *BudgetServer) RefundBudget(ctx context.Context, req *pb.RefundBudgetRequest) (*pb.RefundBudgetResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	log.Printf("退还预算: CampaignID=%s, Amount=%.2f, Reason=%s", req.CampaignId, req.Amount, req.Reason)
	
	budget, exists := s.budgets[req.CampaignId]
	if !exists {
		return &pb.RefundBudgetResponse{
			Success:   false,
			Remaining: 0,
			Message:   "活动不存在",
		}, nil
	}
	
	// 退还预算
	budget.RemainingBudget += req.Amount
	budget.DailySpent -= req.Amount
	
	log.Printf("预算退还成功: CampaignID=%s, Remaining=%.2f", req.CampaignId, budget.RemainingBudget)
	
	return &pb.RefundBudgetResponse{
		Success:   true,
		Remaining: budget.RemainingBudget,
		Message:   "退还成功",
	}, nil
}

func main() {
	// 创建 gRPC 服务器
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("监听失败: %v", err)
	}
	
	grpcServer := grpc.NewServer()
	budgetServer := NewBudgetServer()
	
	pb.RegisterBudgetServiceServer(grpcServer, budgetServer)
	
	log.Println("======================================")
	log.Println("Budget gRPC 服务启动成功")
	log.Println("监听地址: localhost:50052")
	log.Println("======================================")
	
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}


