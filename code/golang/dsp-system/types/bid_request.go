package types

// 这个包用于存放共享类型，避免循环导入

// BidRequestContext 竞价请求上下文（简化版）
type BidRequestContext struct {
	ID      string
	Imps    []ImpContext
	UserID  string
	Device  *DeviceContext
}

// ImpContext 广告位上下文
type ImpContext struct {
	ID       string
	BannerW  int
	BannerH  int
	BidFloor float64
}

// DeviceContext 设备上下文
type DeviceContext struct {
	IP         string
	DeviceType int
	OS         string
}

