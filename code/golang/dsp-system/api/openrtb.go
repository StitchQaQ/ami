package api

// OpenRTB 2.5 简化版协议定义

// BidRequest OpenRTB竞价请求
type BidRequest struct {
	ID      string   `json:"id"`                // 请求ID
	Imp     []Imp    `json:"imp"`               // 广告位列表
	Site    *Site    `json:"site,omitempty"`    // 网站信息
	App     *App     `json:"app,omitempty"`     // APP信息
	Device  *Device  `json:"device,omitempty"`  // 设备信息
	User    *User    `json:"user,omitempty"`    // 用户信息
	Test    int      `json:"test,omitempty"`    // 测试标志 0=正式 1=测试
	TMax    int      `json:"tmax,omitempty"`    // 超时时间(ms)
	WSeat   []string `json:"wseat,omitempty"`   // 白名单席位
	BSeat   []string `json:"bseat,omitempty"`   // 黑名单席位
	AllImps int      `json:"allimps,omitempty"` // 0=部分 1=全部
	Cur     []string `json:"cur,omitempty"`     // 货币类型
}

// Imp 广告位
type Imp struct {
	ID          string   `json:"id"`                    // 广告位ID
	Banner      *Banner  `json:"banner,omitempty"`      // Banner广告
	Video       *Video   `json:"video,omitempty"`       // 视频广告
	Native      *Native  `json:"native,omitempty"`      // 原生广告
	DisplayMgr  string   `json:"displaymanager,omitempty"` // 渲染器
	TagID       string   `json:"tagid,omitempty"`       // 广告位标识
	BidFloor    float64  `json:"bidfloor,omitempty"`    // 底价
	BidFloorCur string   `json:"bidfloorcur,omitempty"` // 底价货币
	Secure      *int     `json:"secure,omitempty"`      // HTTPS标志
}

// Banner Banner广告信息
type Banner struct {
	W      int      `json:"w,omitempty"`      // 宽度
	H      int      `json:"h,omitempty"`      // 高度
	WMax   int      `json:"wmax,omitempty"`   // 最大宽度
	WMin   int      `json:"wmin,omitempty"`   // 最小宽度
	HMax   int      `json:"hmax,omitempty"`   // 最大高度
	HMin   int      `json:"hmin,omitempty"`   // 最小高度
	ID     string   `json:"id,omitempty"`     // Banner ID
	Pos    int      `json:"pos,omitempty"`    // 广告位置
	BType  []int    `json:"btype,omitempty"`  // 屏蔽的Banner类型
	BAttr  []int    `json:"battr,omitempty"`  // 屏蔽的创意属性
	MIMEs  []string `json:"mimes,omitempty"`  // 支持的MIME类型
	API    []int    `json:"api,omitempty"`    // 支持的API框架
}

// Video 视频广告信息
type Video struct {
	MIMEs          []string `json:"mimes"`                    // 支持的MIME类型
	MinDuration    int      `json:"minduration,omitempty"`    // 最小时长(秒)
	MaxDuration    int      `json:"maxduration,omitempty"`    // 最大时长(秒)
	Protocols      []int    `json:"protocols,omitempty"`      // 支持的协议
	W              int      `json:"w,omitempty"`              // 宽度
	H              int      `json:"h,omitempty"`              // 高度
	StartDelay     int      `json:"startdelay,omitempty"`     // 开始延迟
	Linearity      int      `json:"linearity,omitempty"`      // 线性度
	Skip           int      `json:"skip,omitempty"`           // 可跳过
	SkipMin        int      `json:"skipmin,omitempty"`        // 最小跳过时间
	SkipAfter      int      `json:"skipafter,omitempty"`      // 跳过前必看时间
	Sequence       int      `json:"sequence,omitempty"`       // 序列号
	BAttr          []int    `json:"battr,omitempty"`          // 屏蔽的创意属性
	MaxExtended    int      `json:"maxextended,omitempty"`    // 最大扩展时长
	MinBitrate     int      `json:"minbitrate,omitempty"`     // 最小比特率
	MaxBitrate     int      `json:"maxbitrate,omitempty"`     // 最大比特率
	BoxingAllowed  int      `json:"boxingallowed,omitempty"`  // 允许装箱
	PlaybackMethod []int    `json:"playbackmethod,omitempty"` // 播放方式
	Delivery       []int    `json:"delivery,omitempty"`       // 交付方式
	Pos            int      `json:"pos,omitempty"`            // 广告位置
	API            []int    `json:"api,omitempty"`            // 支持的API框架
}

// Native 原生广告信息
type Native struct {
	Request string   `json:"request"`          // 原生广告请求
	Ver     string   `json:"ver,omitempty"`    // 原生广告版本
	API     []int    `json:"api,omitempty"`    // 支持的API框架
	BAttr   []int    `json:"battr,omitempty"`  // 屏蔽的创意属性
}

// Site 网站信息
type Site struct {
	ID     string `json:"id,omitempty"`     // 网站ID
	Name   string `json:"name,omitempty"`   // 网站名称
	Domain string `json:"domain,omitempty"` // 网站域名
	Cat    []string `json:"cat,omitempty"`  // IAB内容分类
	Page   string `json:"page,omitempty"`   // 当前页面URL
	Ref    string `json:"ref,omitempty"`    // 来源页面URL
}

// App APP信息
type App struct {
	ID     string   `json:"id,omitempty"`     // APP ID
	Name   string   `json:"name,omitempty"`   // APP名称
	Bundle string   `json:"bundle,omitempty"` // APP包名
	Domain string   `json:"domain,omitempty"` // APP域名
	Cat    []string `json:"cat,omitempty"`    // IAB内容分类
	Ver    string   `json:"ver,omitempty"`    // APP版本
}

// Device 设备信息
type Device struct {
	UA            string  `json:"ua,omitempty"`            // User Agent
	Geo           *Geo    `json:"geo,omitempty"`           // 地理位置
	DNT           *int    `json:"dnt,omitempty"`           // Do Not Track
	Lmt           *int    `json:"lmt,omitempty"`           // Limit Ad Tracking
	IP            string  `json:"ip,omitempty"`            // IPv4地址
	IPv6          string  `json:"ipv6,omitempty"`          // IPv6地址
	DeviceType    int     `json:"devicetype,omitempty"`    // 设备类型
	Make          string  `json:"make,omitempty"`          // 设备制造商
	Model         string  `json:"model,omitempty"`         // 设备型号
	OS            string  `json:"os,omitempty"`            // 操作系统
	OSV           string  `json:"osv,omitempty"`           // 操作系统版本
	HWV           string  `json:"hwv,omitempty"`           // 硬件版本
	W             int     `json:"w,omitempty"`             // 屏幕宽度
	H             int     `json:"h,omitempty"`             // 屏幕高度
	PPI           int     `json:"ppi,omitempty"`           // 屏幕像素密度
	PxRatio       float64 `json:"pxratio,omitempty"`       // 物理像素比
	JS            int     `json:"js,omitempty"`            // 支持JavaScript
	Language      string  `json:"language,omitempty"`      // 浏览器语言
	Carrier       string  `json:"carrier,omitempty"`       // 运营商
	ConnectionType int    `json:"connectiontype,omitempty"` // 网络连接类型
	IFA           string  `json:"ifa,omitempty"`           // 广告ID (IDFA/AAID)
	DIDSHA1       string  `json:"didsha1,omitempty"`       // 设备ID SHA1
	DIDMD5        string  `json:"didmd5,omitempty"`        // 设备ID MD5
	DPIDSHA1      string  `json:"dpidsha1,omitempty"`      // 平台设备ID SHA1
	DPIDMD5       string  `json:"dpidmd5,omitempty"`       // 平台设备ID MD5
}

// Geo 地理位置信息
type Geo struct {
	Lat     float64 `json:"lat,omitempty"`     // 纬度
	Lon     float64 `json:"lon,omitempty"`     // 经度
	Type    int     `json:"type,omitempty"`    // 位置类型
	Country string  `json:"country,omitempty"` // 国家
	Region  string  `json:"region,omitempty"`  // 地区
	City    string  `json:"city,omitempty"`    // 城市
	ZIP     string  `json:"zip,omitempty"`     // 邮编
}

// User 用户信息
type User struct {
	ID         string `json:"id,omitempty"`         // 用户ID
	BuyerUID   string `json:"buyeruid,omitempty"`   // 买方用户ID
	YOB        int    `json:"yob,omitempty"`        // 出生年份
	Gender     string `json:"gender,omitempty"`     // 性别
	Keywords   string `json:"keywords,omitempty"`   // 关键词
	CustomData string `json:"customdata,omitempty"` // 自定义数据
	Geo        *Geo   `json:"geo,omitempty"`        // 地理位置
	Data       []Data `json:"data,omitempty"`       // 第三方数据
}

// Data 第三方数据
type Data struct {
	ID      string    `json:"id,omitempty"`   // 数据提供商ID
	Name    string    `json:"name,omitempty"` // 数据提供商名称
	Segment []Segment `json:"segment,omitempty"` // 数据段
}

// Segment 数据段
type Segment struct {
	ID    string `json:"id,omitempty"`    // 段ID
	Name  string `json:"name,omitempty"`  // 段名称
	Value string `json:"value,omitempty"` // 段值
}

// BidResponse OpenRTB竞价响应
type BidResponse struct {
	ID         string      `json:"id"`                   // 请求ID（必须与请求ID一致）
	SeatBid    []SeatBid   `json:"seatbid,omitempty"`    // 席位竞价列表
	BidID      string      `json:"bidid,omitempty"`      // 竞价ID
	Cur        string      `json:"cur,omitempty"`        // 货币类型
	CustomData string      `json:"customdata,omitempty"` // 自定义数据
	NBR        int         `json:"nbr,omitempty"`        // 不竞价原因码
}

// SeatBid 席位竞价
type SeatBid struct {
	Bid   []Bid  `json:"bid"`             // 竞价列表
	Seat  string `json:"seat,omitempty"`  // 席位ID
	Group int    `json:"group,omitempty"` // 0=独立竞价 1=组合竞价
}

// Bid 竞价
type Bid struct {
	ID      string   `json:"id"`                // 竞价ID
	ImpID   string   `json:"impid"`             // 对应的广告位ID
	Price   float64  `json:"price"`             // 竞价价格
	AdID    string   `json:"adid,omitempty"`    // 广告ID
	NURL    string   `json:"nurl,omitempty"`    // 竞价成功通知URL
	BURL    string   `json:"burl,omitempty"`    // 计费通知URL
	LURL    string   `json:"lurl,omitempty"`    // 竞价失败通知URL
	AdM     string   `json:"adm,omitempty"`     // 广告素材标记
	AdDomain []string `json:"adomain,omitempty"` // 广告主域名
	Bundle   string   `json:"bundle,omitempty"`  // APP包名
	IURL    string   `json:"iurl,omitempty"`    // 广告图片URL
	CampaignID string `json:"cid,omitempty"`    // 活动ID
	CreativeID string `json:"crid,omitempty"`   // 创意ID
	Cat     []string `json:"cat,omitempty"`     // 创意分类
	Attr    []int    `json:"attr,omitempty"`    // 创意属性
	W       int      `json:"w,omitempty"`       // 宽度
	H       int      `json:"h,omitempty"`       // 高度
}

