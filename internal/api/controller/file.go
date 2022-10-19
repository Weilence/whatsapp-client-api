package controller

import (
	"github.com/weilence/whatsapp-client/internal/api"
	"log"
	"os"
)

func init() {
	err := os.MkdirAll("uploads", os.ModePerm)
	if err != nil {
		log.Panic(err)
	}
}

func UploadAdd(c *api.HttpContext, _ *struct{}) (interface{}, error) {
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

func UploadGet(c *api.HttpContext, req *uploadGetReq) (interface{}, error) {
	c.File("uploads/" + req.Path)
	return nil, nil
}
