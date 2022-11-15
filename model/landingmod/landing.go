package landingmod

type SetttingResponse struct {
	// Entrys 跳转入口
	Entrys []Entry `json:"entry"`
}

type Entry struct {
	// Name 入口名称, eg: 我是商户
	Name string `json:"name"`
	// URL 跳转链接, eg: www.baidu.com
	URL string `json:"url"`
	// BackgroundRPGA 背景色, eg: #F78870
	BackgroundRPGA string `json:"backgroundRPGA"`
	// Show ture: 显示，否则不显示
	Show bool `json:"show"`
}
