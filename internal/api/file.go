package api

import (
	"github.com/gin-gonic/gin"
	"os"
)

func init() {
	err := os.MkdirAll("uploads", os.ModePerm)
	if err != nil {
		panic(err)
	}
}

func UploadAdd(c *gin.Context) {
	f, err := c.FormFile("file")
	if err != nil {
		panic(err)
	}

	err = c.SaveUploadedFile(f, "uploads/"+f.Filename)
	if err != nil {
		panic(err)
	}

	c.JSON(0, f.Filename)
}

func UploadGet(c *gin.Context) {
	filePath := c.Query("path")
	c.File("uploads/" + filePath)
}
