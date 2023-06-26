package controller

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/weilence/whatsapp-client/internal/api"
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
		return nil, fmt.Errorf("get form file: %w", err)
	}

	src, err := f.Open()
	if err != nil {
		return nil, fmt.Errorf("open upload file: %w", err)
	}
	defer src.Close()

	dst, err := os.Open("uploads/" + f.Filename)
	if err != nil {
		return nil, fmt.Errorf("open create file: %w", err)
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return nil, err
	}

	return f.Filename, nil
}

type uploadGetReq struct {
	Path string `query:"path"`
}

func UploadGet(c *api.HttpContext, req *uploadGetReq) (interface{}, error) {
	c.File("uploads/" + req.Path)
	return nil, nil
}
