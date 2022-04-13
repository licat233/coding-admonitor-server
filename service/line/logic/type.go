package logic

type RespData struct {
	Datas  []AccountCostData `json:"datas"`
	Paging struct {
		Page          int      `json:"page"`
		Size          int      `json:"size"`
		TotalElements int      `json:"total_elements"`
		Sorts         []string `json:"sorts"`
	} `json:"paging"`
	ReportCommand struct {
		Since string `json:"since"`
		Until string `json:"until"`
	} `json:"reportCommand"`
}

type AccountCostData struct {
	Id               string `json:"id"`
	AccountInfo      `json:"account"`
	ConfiguredStatus string `json:"configuredStatus"`
	DeliveryStatus   string `json:"deliveryStatus"`
	Statistics       `json:"statistics"`
}

type AccountInfo struct {
	Id               string `json:"id"`
	Name             string `json:"name"`
	ProductName      string `json:"productName"`
	Timezone         string `json:"timezone"`
	ConfiguredStatus string `json:"configuredStatus"`
	CorporateName    string `json:"corporateName"`
	Country          string `json:"country"`
	Currency         string `json:"currency"`
	CurrencySymbol   string `json:"currencySymbol"`
	DeliveryStatus   string `json:"deliveryStatus"`
}

type Statistics struct {
	Imp            uint    `json:"imp"`
	Click          uint    `json:"click"`
	Cv             uint    `json:"cv"`
	Ctr            float64 `json:"ctr"`
	Cvr            float64 `json:"cvr"`
	Cpm            float64 `json:"cpm"`
	Cpc            float64 `json:"cpc"`
	Cpa            float64 `json:"cpa"`
	Cost           float64 `json:"cost"`
	Currency       string  `json:"currency"`
	CurrencySymbol string  `json:"currencySymbol"`
}

type AccountConfig struct {
	Account  string `json:"account"`
	Nickname string `json:"nickname"`
	Status   bool   `json:"status"`
}
