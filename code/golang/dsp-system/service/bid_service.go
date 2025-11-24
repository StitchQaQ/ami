package service

import (
	"context"
	"dsp-system/api"
	"dsp-system/repository"
	"dsp-system/rpc"
	"errors"
	"fmt"
	"log"
	"time"
)

// BidService 竞价服务
type BidService struct {
	adSelector    *AdSelector
	userClient    *rpc.UserClient
	budgetClient  *rpc.BudgetClient
	redisCache    *repository.RedisCache
	clickhouseRepo *repository.ClickHouseRepo
}

// NewBidService 创建竞价服务
func NewBidService(
	adSelector *AdSelector,
	userClient *rpc.UserClient,
	budgetClient *rpc.BudgetClient,
	redisCache *repository.RedisCache,
	clickhouseRepo *repository.ClickHouseRepo,
) *BidService {
	return &BidService{
		adSelector:    adSelector,
		userClient:    userClient,
		budgetClient:  budgetClient,
		redisCache:    redisCache,
		clickhouseRepo: clickhouseRepo,
	}
}

// ProcessBid 处理竞价请求
func (s *BidService) ProcessBid(ctx context.Context, req *api.BidRequest) (*api.BidResponse, error) {
	startTime := time.Now()

	// 1. 获取用户标签（并行调用）
	userProfileChan := make(chan *rpc.UserProfile, 1)
	go func() {
		profile, err := s.getUserProfile(ctx, req)
		if err != nil {
			log.Printf("获取用户画像失败: %v", err)
			userProfileChan <- nil
			return
		}
		userProfileChan <- profile
	}()

	// 2. 广告位信息解析
	if len(req.Imp) == 0 {
		return nil, errors.New("no impressions in request")
	}

	// 3. 获取用户画像结果
	userProfile := <-userProfileChan
	if userProfile == nil {
		userProfile = &rpc.UserProfile{
			UserID: s.extractUserID(req),
			Tags:   []string{},
		}
	}

	log.Printf("用户画像: UserID=%s, Tags=%v", userProfile.UserID, userProfile.Tags)

	// 4. 选择候选广告
	candidates := s.adSelector.SelectAds(ctx, req, userProfile)
	if len(candidates) == 0 {
		log.Printf("无匹配的广告")
		return nil, nil
	}

	log.Printf("候选广告数量: %d", len(candidates))

	// 6. 预算校验和出价计算
	var bids []api.Bid
	for _, candidate := range candidates {
		// 检查预算
		hasBudget, err := s.budgetClient.CheckBudget(ctx, candidate.CampaignID, candidate.BidPrice)
		if err != nil {
			log.Printf("预算检查失败: %v", err)
			continue
		}

		if !hasBudget {
			log.Printf("预算不足: CampaignID=%s", candidate.CampaignID)
			continue
		}

		// 构建竞价响应
		bid := api.Bid{
			ID:         fmt.Sprintf("bid_%s_%d", req.ID, time.Now().UnixNano()),
			ImpID:      candidate.ImpID,
			Price:      candidate.BidPrice,
			AdID:       candidate.AdID,
			AdM:        candidate.Creative,
			NURL:       fmt.Sprintf("http://dsp.example.com/win?price=${AUCTION_PRICE}"),
			BURL:       fmt.Sprintf("http://dsp.example.com/bill?bidid=%s", candidate.AdID),
			CampaignID: candidate.CampaignID,
			CreativeID: candidate.CreativeID,
			AdDomain:   []string{candidate.Domain},
			W:          candidate.Width,
			H:          candidate.Height,
		}

		bids = append(bids, bid)

		// 限制返回数量（避免超时）
		if len(bids) >= 3 {
			break
		}
	}

	if len(bids) == 0 {
		log.Printf("预算校验后无可用广告")
		return nil, nil
	}

	// 6. 记录竞价日志（异步）
	go s.logBidRequest(req, bids, time.Since(startTime))

	// 7. 构建响应
	response := &api.BidResponse{
		ID:  req.ID,
		Cur: "CNY",
		SeatBid: []api.SeatBid{
			{
				Bid:  bids,
				Seat: "dsp-seat-001",
			},
		},
	}

	log.Printf("竞价完成: Bids=%d, Duration=%dms", len(bids), time.Since(startTime).Milliseconds())

	return response, nil
}

// getUserProfile 获取用户画像
func (s *BidService) getUserProfile(ctx context.Context, req *api.BidRequest) (*rpc.UserProfile, error) {
	userID := s.extractUserID(req)
	if userID == "" {
		return nil, errors.New("no user id found")
	}

	// 先查缓存
	cacheKey := fmt.Sprintf("user_profile:%s", userID)
	if profile := s.redisCache.GetUserProfile(ctx, cacheKey); profile != nil {
		return profile, nil
	}

	// 调用RPC获取用户画像
	profile, err := s.userClient.GetUserProfile(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 写入缓存
	s.redisCache.SetUserProfile(ctx, cacheKey, profile, 5*time.Minute)

	return profile, nil
}

// extractUserID 提取用户ID
func (s *BidService) extractUserID(req *api.BidRequest) string {
	if req.User != nil && req.User.ID != "" {
		return req.User.ID
	}
	if req.Device != nil && req.Device.IFA != "" {
		return req.Device.IFA
	}
	if req.Device != nil && req.Device.DIDSHA1 != "" {
		return req.Device.DIDSHA1
	}
	return ""
}

// logBidRequest 记录竞价日志
func (s *BidService) logBidRequest(req *api.BidRequest, bids []api.Bid, duration time.Duration) {
	err := s.clickhouseRepo.LogBidRequest(context.Background(), req, bids, duration)
	if err != nil {
		log.Printf("记录竞价日志失败: %v", err)
	}
}


