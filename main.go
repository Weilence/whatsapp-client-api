package main

import (
	"flag"
	"fmt"

	"github.com/weilence/whatsapp-client/config"
	"github.com/weilence/whatsapp-client/internal/api/router"
)

var (
	version   string
	commitID  string
	buildTime string
)

func main() {
	flag.Parse()

	if *config.ShowV {
		fmt.Println("Version:", version)
		fmt.Println("Commit ID:", commitID)
		fmt.Println("Build Time:", buildTime)
		return
	}

	router.RunServer()
}
