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

// BudgetInfo 预算信息
type BudgetInfo struct {
	CampaignID      string
	TotalBudget     float64
	UsedBudget      float64
	RemainingBudget float64
	DailyLimit      float64
	Status          string
}

// BudgetClient 预算服务客户端
type BudgetClient struct {
	addr   string
	conn   *grpc.ClientConn
	client pb.BudgetServiceClient
}

// NewBudgetClient 创建预算服务客户端
func NewBudgetClient(addr string) *BudgetClient {
	return &BudgetClient{
		addr: addr,
	}
}

// Connect 连接到gRPC服务
func (c *BudgetClient) Connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		c.addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return fmt.Errorf("连接预算服务失败: %v", err)
	}

	c.conn = conn
	c.client = pb.NewBudgetServiceClient(conn)

	log.Printf("预算服务连接成功: %s", c.addr)
	return nil
}

// Close 关闭连接
func (c *BudgetClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// CheckBudget 检查预算是否充足
func (c *BudgetClient) CheckBudget(ctx context.Context, campaignID string, bidPrice float64) (bool, error) {
	if c.client == nil {
		// 降级逻辑：gRPC 未连接时使用模拟数据
		log.Printf("预算服务未连接，使用模拟数据: CampaignID=%s, BidPrice=%.2f", campaignID, bidPrice)
		return bidPrice <= 10.0, nil
	}

	req := &pb.CheckBudgetRequest{
		CampaignId: campaignID,
		Amount:     bidPrice,
	}

	resp, err := c.client.CheckBudget(ctx, req)
	if err != nil {
		log.Printf("检查预算失败: %v", err)
		return false, err
	}

	log.Printf("检查预算: CampaignID=%s, HasBudget=%v, Remaining=%.2f",
		campaignID, resp.HasBudget, resp.Remaining)

	return resp.HasBudget, nil
}

// DeductBudget 扣减预算（竞价成功后）
func (c *BudgetClient) DeductBudget(ctx context.Context, campaignID string, amount float64) error {
	if c.client == nil {
		log.Printf("预算服务未连接，跳过扣减: CampaignID=%s, Amount=%.2f", campaignID, amount)
		return nil
	}

	req := &pb.DeductBudgetRequest{
		CampaignId: campaignID,
		Amount:     amount,
		BidId:      fmt.Sprintf("bid_%d", time.Now().UnixNano()),
	}

	resp, err := c.client.DeductBudget(ctx, req)
	if err != nil {
		log.Printf("扣减预算失败: %v", err)
		return err
	}

	if !resp.Success {
		return fmt.Errorf("扣减失败: %s", resp.Message)
	}

	log.Printf("扣减预算成功: CampaignID=%s, Remaining=%.2f", campaignID, resp.Remaining)
	return nil
}

// GetBudgetInfo 获取预算信息
func (c *BudgetClient) GetBudgetInfo(ctx context.Context, campaignID string) (*BudgetInfo, error) {
	if c.client == nil {
		log.Printf("预算服务未连接，返回模拟数据: CampaignID=%s", campaignID)
		return &BudgetInfo{
			CampaignID:      campaignID,
			TotalBudget:     10000.0,
			UsedBudget:      3500.0,
			RemainingBudget: 6500.0,
			DailyLimit:      1000.0,
			Status:          "active",
		}, nil
	}

	req := &pb.GetBudgetInfoRequest{
		CampaignId: campaignID,
	}

	resp, err := c.client.GetBudgetInfo(ctx, req)
	if err != nil {
		log.Printf("获取预算信息失败: %v", err)
		return nil, err
	}

	info := &BudgetInfo{
		CampaignID:      resp.CampaignId,
		TotalBudget:     resp.TotalBudget,
		UsedBudget:      resp.TotalBudget - resp.RemainingBudget,
		RemainingBudget: resp.RemainingBudget,
		DailyLimit:      resp.DailyBudget,
		Status:          resp.Status,
	}

	log.Printf("获取预算信息: CampaignID=%s, Remaining=%.2f", campaignID, info.RemainingBudget)
	return info, nil
}

// RefundBudget 退还预算（竞价失败时）
func (c *BudgetClient) RefundBudget(ctx context.Context, campaignID string, amount float64) error {
	if c.client == nil {
		log.Printf("预算服务未连接，跳过退还: CampaignID=%s, Amount=%.2f", campaignID, amount)
		return nil
	}

	req := &pb.RefundBudgetRequest{
		CampaignId: campaignID,
		Amount:     amount,
		BidId:      fmt.Sprintf("bid_%d", time.Now().UnixNano()),
		Reason:     "bid_failed",
	}

	resp, err := c.client.RefundBudget(ctx, req)
	if err != nil {
		log.Printf("退还预算失败: %v", err)
		return err
	}

	if !resp.Success {
		return fmt.Errorf("退还失败: %s", resp.Message)
	}

	log.Printf("退还预算成功: CampaignID=%s, Remaining=%.2f", campaignID, resp.Remaining)
	return nil
}

// BatchCheckBudget 批量检查预算
func (c *BudgetClient) BatchCheckBudget(ctx context.Context, requests map[string]float64) (map[string]bool, error) {
	// requests: campaignID -> bidPrice
	results := make(map[string]bool)

	for campaignID, bidPrice := range requests {
		hasBudget, err := c.CheckBudget(ctx, campaignID, bidPrice)
		if err != nil {
			log.Printf("检查预算失败: CampaignID=%s, Error=%v", campaignID, err)
			results[campaignID] = false
			continue
		}
		results[campaignID] = hasBudget
	}

	return results, nil
}
