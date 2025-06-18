package server

import (
	"fmt"
	"go-pixiv-proxy/conf"
	"net/http"
	"strings"

	"github.com/lfhy/log"
)

func Run() {
	InitClient()
	if conf.Domain != "" {
		indexHtml = strings.ReplaceAll(indexHtml, "{image-examples}", docExampleImg)
		indexHtml = strings.ReplaceAll(indexHtml, "http://example.com", conf.Domain)
	}
	http.HandleFunc("/", handlePixivProxy)
	log.Infof("started: %s:%s %s", conf.Host, conf.Port, conf.Domain)
	err := http.ListenAndServe(fmt.Sprintf("%s:%s", conf.Host, conf.Port), nil)
	if err != nil {
		log.Error("start failed: ", err)
	}
}
