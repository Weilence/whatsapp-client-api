package model

type ResponseModel struct {
	Code    int         `json:"code"`
	Message string      `json:"message,omitempty"`
	Result  interface{} `json:"result,omitempty"`
}

type ResponseList struct {
	Total int64 `json:"total"`
	List  any   `json:"list"`
}
