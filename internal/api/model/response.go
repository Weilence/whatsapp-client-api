package model

type ResponseModel struct {
	Code    int         `json:"code"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type ResponseList struct {
	Total int64 `json:"total"`
	List  any   `json:"list"`
}
