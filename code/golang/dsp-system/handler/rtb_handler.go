package handler

import (
	"dsp-system/api"
	"dsp-system/service"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// RTBHandler RTB请求处理器
type RTBHandler struct {
	bidService *service.BidService
}

// NewRTBHandler 创建RTB处理器
func NewRTBHandler(bidService *service.BidService) *RTBHandler {
	return &RTBHandler{
		bidService: bidService,
	}
}

// HandleBidRequest 处理竞价请求
// POST /bid
func (h *RTBHandler) HandleBidRequest(c *gin.Context) {
	startTime := time.Now()
	
	// 1. 解析请求
	var bidRequest api.BidRequest
	if err := c.ShouldBindJSON(&bidRequest); err != nil {
		log.Printf("解析竞价请求失败: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	log.Printf("收到竞价请求: ID=%s, Imps=%d", bidRequest.ID, len(bidRequest.Imp))

	// 2. 快速校验
	if bidRequest.ID == "" || len(bidRequest.Imp) == 0 {
		log.Printf("竞价请求参数不完整: %+v", bidRequest)
		h.sendNoBid(c, bidRequest.ID, 2) // NBR=2: 技术错误
		return
	}

	// 3. 调用竞价服务
	bidResponse, err := h.bidService.ProcessBid(c.Request.Context(), &bidRequest)
	if err != nil {
		log.Printf("竞价处理失败: %v", err)
		h.sendNoBid(c, bidRequest.ID, 3) // NBR=3: 无效请求
		return
	}

	// 4. 返回竞价响应
	if bidResponse == nil || len(bidResponse.SeatBid) == 0 {
		log.Printf("无可用广告返回")
		h.sendNoBid(c, bidRequest.ID, 1) // NBR=1: 技术原因
		return
	}

	// 5. 设置响应头
	c.Header("X-Processing-Time-Ms", fmt.Sprintf("%d", time.Since(startTime).Milliseconds()))
	
	log.Printf("竞价成功: RequestID=%s, Bids=%d, Time=%dms", 
		bidRequest.ID, 
		len(bidResponse.SeatBid[0].Bid),
		time.Since(startTime).Milliseconds())

	c.JSON(http.StatusOK, bidResponse)
}

// sendNoBid 发送无竞价响应
func (h *RTBHandler) sendNoBid(c *gin.Context, requestID string, nbr int) {
	c.JSON(http.StatusNoContent, api.BidResponse{
		ID:  requestID,
		NBR: nbr, // No-Bid Reason
	})
}

// HealthCheck 健康检查
// GET /health
func (h *RTBHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"time":   time.Now().Unix(),
	})
}

// Stats 统计信息
// GET /stats
func (h *RTBHandler) Stats(c *gin.Context) {
	// 这里可以返回DSP系统的统计信息
	// 例如：QPS、成功率、平均响应时间等
	c.JSON(http.StatusOK, gin.H{
		"qps":              0,
		"bid_rate":         0.0,
		"win_rate":         0.0,
		"avg_response_ms":  0,
		"total_requests":   0,
		"total_bids":       0,
		"total_wins":       0,
	})
}

