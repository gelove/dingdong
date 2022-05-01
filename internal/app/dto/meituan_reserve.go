package dto

type MeiTuanReserveResult struct {
	Code int                `json:"code"`
	Data MeiTuanReserveData `json:"data"`
}

type MeiTuanReserveData struct {
	Msg       string `json:"msg"`
	Type      int    `json:"type"`
	BackColor string `json:"backColor"`
	FontColor string `json:"fontColor"`
	CycleType int    `json:"cycleType"`
}
