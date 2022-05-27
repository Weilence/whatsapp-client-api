package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
	"whatsapp-client/config"
	"whatsapp-client/internal/model"
	"whatsapp-client/internal/router"
	"whatsapp-client/pkg/whatsapp"
)

func main() {
	config.Init()
	model.Init()
	whatsapp.Init(model.SqlDB())

	r := router.Setup()

	host := viper.GetString("web.host")
	port := viper.GetString("web.port")
	if viper.GetBool("debug") {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	server := http.Server{
		Addr:    host + ":" + port,
		Handler: r,
	}

	log.Println("Server is running on " + host + ":" + port)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server listen err:%s", err)
		}
	}()

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
