package controller

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func init() {
	err := os.MkdirAll("uploads", os.ModePerm)
	if err != nil {
		log.Panic(err)
	}
}

func UploadAdd(c *gin.Context, _ *struct{}) (interface{}, error) {
	f, err := c.FormFile("file")
	if err != nil {
		return nil, err
	}

	err = c.SaveUploadedFile(f, "uploads/"+f.Filename)
	if err != nil {
		return nil, err
	}

	return f.Filename, nil
}

type uploadGetReq struct {
	Path string `form:"path"`
}

func UploadGet(c *gin.Context, req *uploadGetReq) (interface{}, error) {
	c.File("uploads/" + req.Path)
	return nil, nil
}
