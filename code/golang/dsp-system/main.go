package main

import (
	"context"
	"dsp-system/config"
	"dsp-system/handler"
	"dsp-system/logger"
	"dsp-system/repository"
	"dsp-system/rpc"
	"dsp-system/service"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. 加载配置
	cfg := config.LoadConfig()

	// 2. 初始化日志系统
	if err := logger.Init(&cfg.Log); err != nil {
		panic("日志系统初始化失败: " + err.Error())
	}
	logger.Info("配置加载完成")

	// 3. 初始化基础设施层
	redisCache := repository.NewRedisCache(&cfg.Redis)
	defer redisCache.Close()

	clickhouseRepo := repository.NewClickHouseRepo(&cfg.ClickHouse)
	defer clickhouseRepo.Close()

	// 4. 初始化RPC客户端
	userClient := rpc.NewUserClient(cfg.RPC.UserServiceAddr)
	budgetClient := rpc.NewBudgetClient(cfg.RPC.BudgetServiceAddr)

	// 连接gRPC服务（带重试和降级）
	if err := userClient.Connect(); err != nil {
		logger.Warnf("用户服务连接失败(将使用降级逻辑): %v", err)
	} else {
		defer userClient.Close()
		logger.Info("用户服务连接成功")
	}

	if err := budgetClient.Connect(); err != nil {
		logger.Warnf("预算服务连接失败(将使用降级逻辑): %v", err)
	} else {
		defer budgetClient.Close()
		logger.Info("预算服务连接成功")
	}

	logger.Info("RPC客户端初始化完成")

	// 5. 初始化服务层
	adSelector := service.NewAdSelector()
	bidService := service.NewBidService(
		adSelector,
		userClient,
		budgetClient,
		redisCache,
		clickhouseRepo,
	)

	// 6. 初始化Handler层
	rtbHandler := handler.NewRTBHandler(bidService)

	// 7. 配置Gin
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = logger.GetGinWriter()
	router := gin.New()
	router.Use(gin.LoggerWithWriter(logger.GetGinWriter()))
	router.Use(gin.Recovery())

	// 8. 注册路由
	// 健康检查
	router.GET("/health", rtbHandler.HealthCheck)
	router.GET("/stats", rtbHandler.Stats)

	// RTB竞价接口
	router.POST("/bid", rtbHandler.HandleBidRequest)

	// 竞价结果回调
	router.GET("/win", handleWinNotice)
	router.GET("/bill", handleBillNotice)

	// 9. 启动HTTP服务器
	srv := &http.Server{
		Addr:           ":" + cfg.Server.Port,
		Handler:        router,
		ReadTimeout:    100 * time.Millisecond, // RTB对响应时间要求很高
		WriteTimeout:   100 * time.Millisecond,
		MaxHeaderBytes: 1 << 20,
	}

	// 10. 优雅启动和关闭
	go func() {
		logger.Infof("DSP服务启动: http://localhost:%s", cfg.Server.Port)
		logger.Info("======================================")
		logger.Info("API文档:")
		logger.Info("  POST /bid        - OpenRTB竞价接口")
		logger.Info("  GET  /health     - 健康检查")
		logger.Info("  GET  /stats      - 统计信息")
		logger.Info("  GET  /win        - 赢标通知")
		logger.Info("  GET  /bill       - 计费通知")
		logger.Info("======================================")

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("服务器启动失败: %v", err)
		}
	}()

	// 11. 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("正在关闭服务器...")

	// 12. 优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Errorf("服务器关闭错误: %v", err)
	}

	logger.Info("服务器已关闭")
}

// handleWinNotice 处理赢标通知
func handleWinNotice(c *gin.Context) {
	price := c.Query("price")
	bidID := c.Query("bidid")

	logger.Infof("赢标通知: BidID=%s, Price=%s", bidID, price)

	// 实际项目中应该:
	// 1. 记录赢标日志到ClickHouse
	// 2. 扣减预算（调用预算服务）
	// 3. 更新统计数据

	c.String(http.StatusOK, "OK")
}

// handleBillNotice 处理计费通知
func handleBillNotice(c *gin.Context) {
	bidID := c.Query("bidid")

	logger.Infof("计费通知: BidID=%s", bidID)

	// 实际项目中应该:
	// 1. 记录曝光日志
	// 2. 触发广告展示监控

	c.String(http.StatusOK, "OK")
}
