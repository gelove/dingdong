package dto

type Result struct {
	Success bool        `json:"success"`
	Code    int         `json:"code"`
	Msg     string      `json:"msg"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ErrorMessage struct {
	Message string `json:"message"`
}
