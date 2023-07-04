package controller

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/denisbrodbeck/machineid"
	"github.com/mattn/go-ieproxy"
	"github.com/weilence/whatsapp-client/config"
	"github.com/weilence/whatsapp-client/internal/api"
	"github.com/weilence/whatsapp-client/internal/pkg/whatsapp"
)

type MachineInfoRes struct {
	MachineID string `json:"machineId"`
	Version   string `json:"version"`
}

func MachineInfo(_ *api.HttpContext, _ *struct{}) (_ interface{}, err error) {
	machineID, err := machineid.ProtectedID("whatsapp-client")
	if err != nil {
		log.Fatal(err)
	}
	machineID = strings.ToUpper(machineID[:16])

	return MachineInfoRes{
		MachineID: machineID,
		Version:   version,
	}, nil
}

type SetProxyReq struct {
	Proxy string `json:"proxy"`
}

func SetProxy(_ *api.HttpContext, req *SetProxyReq) (*struct{}, error) {
	proxyURL, err := url.Parse(req.Proxy)
	if err != nil {
		return nil, fmt.Errorf("proxy format err: %w", err)
	}

	config.Config.Proxy = req.Proxy
	config.Save()

	clients := whatsapp.GetClients()
	for _, v := range clients {
		v.SetProxy(http.ProxyURL(proxyURL))
	}

	return nil, nil
}

func TestProxy(_ *api.HttpContext, _ *struct{}) (*struct{}, error) {
	proxyConfig := config.Config.Proxy
	var proxy func(*http.Request) (*url.URL, error)
	if proxyConfig == "" {
		proxy = ieproxy.GetProxyFunc()
		log.Printf("test system proxy")
	} else {
		proxyURL, err := url.Parse(proxyConfig)
		if err != nil {
			return nil, fmt.Errorf("proxy format err: %w", err)
		}
		proxy = http.ProxyURL(proxyURL)
		log.Printf("test custom proxy: %s", proxyConfig)
	}

	client := http.Client{
		Transport: &http.Transport{
			Proxy: proxy,
		},
	}

	_, err := client.Get("https://web.whatsapp.com")
	if err != nil {
		return nil, fmt.Errorf("proxy test err: %w", err)
	}

	return nil, nil
}
