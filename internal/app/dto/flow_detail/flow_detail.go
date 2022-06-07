package flow_detail

type Result struct {
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	Success bool   `json:"success"`
	Data    struct {
		List      []Item `json:"list"`
		IsMore    bool   `json:"is_more"`
		PageSize  int    `json:"page_size"`
		StartLoad int    `json:"start_load"`
	} `json:"data"`
	ServerTime    int    `json:"server_time"`
	ExecTime      int    `json:"exec_time"`
	RequestId     string `json:"request_id"`
	PreviewConfig any    `json:"preview_config"`
}

type Item struct {
	Price       string `json:"price"`
	Name        string `json:"name"`
	Spec        string `json:"spec"`
	Sizes       []any  `json:"sizes"`
	Status      int    `json:"status"`
	Type        int    `json:"type"`
	Activity    []any  `json:"activity"`
	Oid         int    `json:"oid"`
	Id          string `json:"id"`
	OriginPrice string `json:"origin_price"`
	VipPrice    string `json:"vip_price"`
	StockNumber int    `json:"stock_number"`
}
