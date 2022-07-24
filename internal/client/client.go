package client

import (
	"log"
	"runtime"

	"github.com/spf13/viper"
	"github.com/zserge/lorca"
)

func OpenBrowser() lorca.UI {
	args := []string{}
	if runtime.GOOS == "linux" {
		args = append(args, "--class=Lorca")
	}
	url := "http://" + viper.GetString("web.host") + ":" + viper.GetString("web.port")
	ui, err := lorca.New(url, "", 1366, 768, args...)
	if err != nil {
		log.Fatal(err)
	}

	return ui
}
