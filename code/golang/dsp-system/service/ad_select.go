package service

import (
	"context"
	"dsp-system/api"
	"dsp-system/rpc"
	"log"
	"math/rand"
)

// AdCandidate 广告候选
type AdCandidate struct {
	AdID        string
	ImpID       string
	CampaignID  string
	CreativeID  string
	BidPrice    float64
	Creative    string
	Domain      string
	Width       int
	Height      int
	TargetTags  []string
	Score       float64
}

// AdSelector 广告选择服务
type AdSelector struct {
	// 这里可以注入广告库、规则引擎等
}

// NewAdSelector 创建广告选择服务
func NewAdSelector() *AdSelector {
	return &AdSelector{}
}

// SelectAds 选择匹配的广告
func (s *AdSelector) SelectAds(ctx context.Context, req *api.BidRequest, userProfile *rpc.UserProfile) []AdCandidate {
	var candidates []AdCandidate

	// 遍历每个广告位
	for _, imp := range req.Imp {
		// 1. 根据广告位类型筛选广告
		ads := s.getAdsByImp(&imp)

		// 2. 根据用户标签匹配广告
		matchedAds := s.matchAdsByUserTags(ads, userProfile)

		// 3. 计算广告得分
		for _, ad := range matchedAds {
			ad.Score = s.calculateAdScore(ad, userProfile, &imp)
			candidates = append(candidates, ad)
		}
	}

	// 4. 排序（按得分降序）
	candidates = s.sortAdsByScore(candidates)

	log.Printf("广告选择完成: Total=%d", len(candidates))

	return candidates
}

// getAdsByImp 根据广告位获取候选广告
func (s *AdSelector) getAdsByImp(imp *api.Imp) []AdCandidate {
	// 实际项目中，这里应该从数据库或缓存中查询
	// 这里为示例，返回模拟数据
	
	mockAds := []AdCandidate{
		{
			AdID:       "ad_001",
			ImpID:      imp.ID,
			CampaignID: "campaign_001",
			CreativeID: "creative_001",
			BidPrice:   5.50,
			Creative:   "<a href='http://example.com'><img src='http://cdn.example.com/ad1.jpg' /></a>",
			Domain:     "example.com",
			TargetTags: []string{"男性", "25-34岁", "运动爱好者"},
		},
		{
			AdID:       "ad_002",
			ImpID:      imp.ID,
			CampaignID: "campaign_002",
			CreativeID: "creative_002",
			BidPrice:   4.80,
			Creative:   "<a href='http://shop.com'><img src='http://cdn.shop.com/ad2.jpg' /></a>",
			Domain:     "shop.com",
			TargetTags: []string{"女性", "18-24岁", "购物达人"},
		},
		{
			AdID:       "ad_003",
			ImpID:      imp.ID,
			CampaignID: "campaign_003",
			CreativeID: "creative_003",
			BidPrice:   6.20,
			Creative:   "<a href='http://tech.com'><img src='http://cdn.tech.com/ad3.jpg' /></a>",
			Domain:     "tech.com",
			TargetTags: []string{"科技爱好者", "程序员"},
		},
	}

	// 根据广告位尺寸筛选
	if imp.Banner != nil {
		for i := range mockAds {
			mockAds[i].Width = imp.Banner.W
			mockAds[i].Height = imp.Banner.H
		}
	}

	return mockAds
}

// matchAdsByUserTags 根据用户标签匹配广告
func (s *AdSelector) matchAdsByUserTags(ads []AdCandidate, userProfile *rpc.UserProfile) []AdCandidate {
	if userProfile == nil || len(userProfile.Tags) == 0 {
		// 没有用户标签，返回所有广告
		return ads
	}

	var matched []AdCandidate

	// 创建用户标签映射
	userTagsMap := make(map[string]bool)
	for _, tag := range userProfile.Tags {
		userTagsMap[tag] = true
	}

	// 筛选匹配的广告
	for _, ad := range ads {
		// 计算标签匹配度
		matchCount := 0
		for _, targetTag := range ad.TargetTags {
			if userTagsMap[targetTag] {
				matchCount++
			}
		}

		// 至少匹配一个标签，或者没有定向标签
		if matchCount > 0 || len(ad.TargetTags) == 0 {
			matched = append(matched, ad)
		}
	}

	log.Printf("标签匹配: UserTags=%v, Matched=%d/%d", userProfile.Tags, len(matched), len(ads))

	return matched
}

// calculateAdScore 计算广告得分
func (s *AdSelector) calculateAdScore(ad AdCandidate, userProfile *rpc.UserProfile, imp *api.Imp) float64 {
	score := 0.0

	// 1. 出价权重（40%）
	score += ad.BidPrice * 0.4

	// 2. 标签匹配度（30%）
	if userProfile != nil && len(userProfile.Tags) > 0 {
		userTagsMap := make(map[string]bool)
		for _, tag := range userProfile.Tags {
			userTagsMap[tag] = true
		}

		matchCount := 0
		for _, targetTag := range ad.TargetTags {
			if userTagsMap[targetTag] {
				matchCount++
			}
		}

		if len(ad.TargetTags) > 0 {
			matchRate := float64(matchCount) / float64(len(ad.TargetTags))
			score += matchRate * 3.0 // 3.0是标签匹配的基础分
		}
	}

	// 3. 随机因子（10%）- 避免总是返回相同广告
	score += rand.Float64() * 1.0

	return score
}

// sortAdsByScore 按得分排序（降序）
func (s *AdSelector) sortAdsByScore(candidates []AdCandidate) []AdCandidate {
	// 简单的冒泡排序（实际项目中应使用更高效的排序算法）
	n := len(candidates)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if candidates[j].Score < candidates[j+1].Score {
				candidates[j], candidates[j+1] = candidates[j+1], candidates[j]
			}
		}
	}
	return candidates
}

