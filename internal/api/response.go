package api

import "github.com/gin-gonic/gin"

type ResponseModel struct {
	Code    int         `json:"code"`
	Message string      `json:"message,omitempty"`
	Result  interface{} `json:"result,omitempty"`
}

func Ok(g *gin.Context, v interface{}) {
	g.JSON(200, &ResponseModel{
		Code:   0,
		Result: v,
	})
}

func BadRequest(g *gin.Context, err error) {
	g.AbortWithStatusJSON(200, &ResponseModel{
		Code:    400,
		Message: err.Error(),
	})
}
