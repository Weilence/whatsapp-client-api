package model

type ResponseModel struct {
	Code    int         `json:"code"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func (m ResponseModel) Error() string {
	return m.Message
}

func Error(code int, message string) *ResponseModel {
	return &ResponseModel{Code: code, Message: message}
}

type ResponseList struct {
	Total int64 `json:"total"`
	List  any   `json:"list"`
}
