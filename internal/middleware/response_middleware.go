package middleware

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
)

type responseWriter struct {
	gin.ResponseWriter
	Code int
}

type responseModel struct {
	Code    int         `json:"code"`
	Message string      `json:"message,omitempty"`
	Result  interface{} `json:"result,omitempty"`
}

func (r *responseWriter) WriteHeader(statusCode int) {
	r.Code = statusCode
	r.ResponseWriter.WriteHeader(200)
}

func (r *responseWriter) Write(b []byte) (n int, err error) {
	header := r.ResponseWriter.Header()
	if header.Get("Content-Type") == "application/json; charset=utf-8" {
		var v interface{}
		err = json.Unmarshal(b, &v)
		if err != nil {
			panic(err)
		}

		if r.Code == 0 {
			b, err = json.Marshal(responseModel{
				Result: v,
			})
		} else {
			b, err = json.Marshal(responseModel{
				Code:    r.Code,
				Message: v.(string),
			})
		}
		if err != nil {
			panic(err)
		}
	}

	return r.ResponseWriter.Write(b)
}

func NewResponse() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer = &responseWriter{
			ResponseWriter: c.Writer,
		}
		c.Next()
	}
}
