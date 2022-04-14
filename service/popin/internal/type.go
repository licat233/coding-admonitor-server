package internal

//AccountConfig 账户信息
type AccountConfig struct {
	AccountID   string `json:"accountId"`
	AccountName string `json:"accountName"`
	Nickname    string `json:"nickname"`
	Status      bool   `json:"status"`
}

/** RespData 正常的响应 **/
type RespData struct {
	Code int    `json:"code"`
	Data Data   `json:"data"`
	Msg  string `json:"msg"`
}

//Detail 广告数据详细
type Detail struct {
	AccountCurrency string `json:"accountCurrency"`
	AccountID       string `json:"accountId"`
	AccountName     string `json:"accountName"`
	Cpc             string `json:"avgCpc"`
	Click           string `json:"click"`
	Cv              string `json:"conversion"`
	Cpa             string `json:"conversionPrice"`
	Cvr             string `json:"conversionRate"`
	Ctr             string `json:"ctr"`
	Currency        string `json:"currency"`
	Imp             string `json:"impression"`
	Cost            string `json:"spend"`
}

//Summary 数据归总
type Summary struct {
	AvgCpc                 string `json:"avgCpc"`
	Click                  string `json:"click"`
	Conversion             string `json:"conversion"`
	ConversionCost         string `json:"conversionCost"`
	ConversionPrice        string `json:"conversionPrice"`
	ConversionRate         string `json:"conversionRate"`
	Ctr                    string `json:"ctr"`
	Cv                     string `json:"cv"`
	CvAddPaymentInfo       string `json:"cvAddPaymentInfo"`
	CvAddToCar             string `json:"cvAddToCar"`
	CvAddToWishlist        string `json:"cvAddToWishlist"`
	CvAppInstall           string `json:"cvAppInstall"`
	CvCompleteRegistration string `json:"cvCompleteRegistration"`
	CvLead                 string `json:"cvLead"`
	CvOther                string `json:"cvOther"`
	CvPurchase             string `json:"cvPurchase"`
	CvStartCheckout        string `json:"cvStartCheckout"`
	CvViewContent          string `json:"cvViewContent"`
	Impression             string `json:"impression"`
	Mcv                    string `json:"mcv"`
	Search                 string `json:"search"`
	Spend                  string `json:"spend"`
}

//Data 数据载体
type Data struct {
	CurrentPage int      `json:"currentPage"`
	Detail      []Detail `json:"detail"`
	Summary     Summary  `json:"summary"`
	Total       int      `json:"total"`
}

/** 心跳的响应 **/
type HeartbeatResp struct {
	Code int `json:"code"`
	Data struct {
		AdAuditNum             string        `json:"adAuditNum"`
		Balance                string        `json:"balance"`
		CreditLineAccountAlarm string        `json:"creditLineAccountAlarm"`
		SspInfo                []interface{} `json:"sspInfo"`
		UserTags               []UserTags    `json:"userTags"`
	} `json:"data"`
	Msg string `json:"msg"`
}

type UserTags struct {
	Label string `json:"label"`
	Uv    int    `json:"uv"`
	Value string `json:"value"`
}
