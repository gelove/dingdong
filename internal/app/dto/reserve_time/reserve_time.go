package reserve_time

type Result struct {
	Success    bool    `json:"success"`
	Code       int     `json:"code"`
	ServerTime int     `json:"server_time"`
	IsTrade    int     `json:"is_trade"`
	Msg        string  `json:"msg"`
	TradeTag   string  `json:"tradeTag"`
	Data       []*Item `json:"data"`
}

type Item struct {
	Closed              bool    `json:"closed"`
	IsNewRules          bool    `json:"is_new_rules"`
	DefaultSelect       bool    `json:"default_select"`
	AreaLevel           int     `json:"area_level"`
	StationID           string  `json:"station_id"`
	BusySoonArrivalText string  `json:"busy_soon_arrival_text"`
	EtaTraceID          string  `json:"eta_trace_id"`
	Times               []*Time `json:"time"`
	StationDelayText    any     `json:"station_delay_text"`
}

type Time struct {
	IsInvalid        bool       `json:"is_invalid"`
	DateStrTimestamp int        `json:"date_str_timestamp"`
	DateStr          string     `json:"date_str"`
	Day              string     `json:"day"`
	Times            []*GoTimes `json:"times"`
	InvalidPrompt    any        `json:"invalid_prompt"`
	TimeFullTextTip  any        `json:"time_full_text_tip"`
}

type GoTimes struct {
	FullFlag       bool   `json:"fullFlag"`
	ArrivalTime    bool   `json:"arrival_time"`
	Type           int    `json:"type"`
	DisableType    int    `json:"disableType"`
	StartTimestamp int64  `json:"start_timestamp"`
	EndTimestamp   int64  `json:"end_timestamp"`
	DisableMsg     string `json:"disableMsg"`
	StartTime      string `json:"start_time"`
	EndTime        string `json:"end_time"`
	TextMsg        string `json:"textMsg"`
	ArrivalTimeMsg string `json:"arrival_time_msg"`
	SelectMsg      string `json:"select_msg"`
}
