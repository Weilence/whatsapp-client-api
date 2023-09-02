package controller

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/denisbrodbeck/machineid"
	"github.com/mattn/go-ieproxy"
	"github.com/weilence/whatsapp-client/config"
	"github.com/weilence/whatsapp-client/internal/pkg/whatsapp"
	"github.com/weilence/whatsapp-client/internal/utils"
)

type MachineInfoRes struct {
	MachineID string `json:"machineId"`
	Version   string `json:"version"`
}

func MachineInfo(_ *utils.HttpContext, _ *struct{}) (_ interface{}, err error) {
	machineID, err := machineid.ProtectedID("whatsapp-client")
	if err != nil {
		panic(err)
	}
	machineID = strings.ToUpper(machineID[:16])

	return MachineInfoRes{
		MachineID: machineID,
		Version:   version,
	}, nil
}

type SetProxyReq struct {
	Url string `json:"url"`
}

func SetProxy(_ *utils.HttpContext, req *SetProxyReq) (*struct{}, error) {
	proxyURL, err := url.Parse(req.Url)
	if err != nil {
		return nil, fmt.Errorf("proxy format err: %w", err)
	}

	config.Config.Proxy = req.Url
	config.Save()

	clients := whatsapp.GetClients()
	for _, v := range clients {
		v.SetProxy(http.ProxyURL(proxyURL))
	}

	return nil, nil
}

func TestProxy(_ *utils.HttpContext, _ *struct{}) (*struct{}, error) {
	proxy, err := getProxy()
	if err != nil {
		return nil, err
	}

	client := http.Client{
		Transport: &http.Transport{
			Proxy: proxy,
		},
	}

	_, err = client.Get("https://web.whatsapp.com")
	if err != nil {
		return nil, fmt.Errorf("proxy test err: %w", err)
	}

	return nil, nil
}

func getProxy() (func(*http.Request) (*url.URL, error), error) {
	var proxy func(*http.Request) (*url.URL, error)

	proxyConfig := config.Config.Proxy
	if proxyConfig == "" {
		proxy = ieproxy.GetProxyFunc()
	} else {
		proxyURL, err := url.Parse(proxyConfig)
		if err != nil {
			return nil, fmt.Errorf("proxy format err: %w", err)
		}
		proxy = http.ProxyURL(proxyURL)
	}

	return proxy, nil
}
