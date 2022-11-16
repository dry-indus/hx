package searchmod

import "hx/global"

type Store struct {
	// StoreName 店铺名
	StoreName string `json:"storeName"`
	// Name 商户头像
	Prtrait string `json:"prtrait"`
	// Category 1:餐饮,2:服饰
	Category global.MerchantCategory `json:"category" enums:"1,2"`
	// Star 用户搜藏量
	Star int `json:"star"`
	// URL 跳转链接 eg: https://www.baidu.com, www.baidu.com
	URL string
}
