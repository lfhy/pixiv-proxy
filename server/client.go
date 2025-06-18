package server

import (
	"go-pixiv-proxy/conf"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/proxy"
)

var dialer proxy.Dialer

func GetDialer() proxy.Dialer {
	if dialer != nil {
		return dialer
	}
	var proxyDialer proxy.Dialer = proxy.Direct
	for _, prxoyUrl := range strings.Split(conf.Proxy, ",") {
		urlInfo, err := url.Parse(prxoyUrl)
		if err != nil {
			return proxyDialer
		}
		var auth *proxy.Auth = nil
		if urlInfo.User != nil {
			pwd, _ := urlInfo.User.Password()
			auth = &proxy.Auth{
				User:     urlInfo.User.Username(),
				Password: pwd,
			}
		}

		dialer, err := proxy.SOCKS5("tcp", urlInfo.Host, auth, proxyDialer)
		if err == nil {
			proxyDialer = dialer
		}
	}
	dialer = proxyDialer
	return dialer
}

func InitClient() {
	if conf.Proxy != "" {
		if strings.HasPrefix(conf.Proxy, "socks5://") {
			dialer := GetDialer()
			client = &http.Client{
				Transport: &http.Transport{
					Dial: dialer.Dial,
				},
			}
		} else if strings.HasPrefix(conf.Proxy, "http://") {
			proxyUrl, err := url.Parse(conf.Proxy)
			if err == nil {
				client = &http.Client{
					Transport: &http.Transport{
						Proxy: http.ProxyURL(proxyUrl),
					},
				}
			}
		}
	}
}
