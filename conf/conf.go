package conf

import (
	"github.com/lfhy/flag"
	"github.com/lfhy/log"
)

var (
	Host    string
	Port    string
	Domain  string
	Cookies string
	Proxy   string
)

func init() {
	flag.StringFullVar(&Host, "h", "server", "host", "GPP_HOST", "0.0.0.0", "Host")
	flag.StringFullVar(&Port, "p", "server", "port", "GPP_PORT", "18090", "Port")
	flag.StringFullVar(&Domain, "d", "server", "domain", "GPP_DOMAIN", "", "Server Domain")
	flag.StringFullVar(&Cookies, "c", "server", "cookies", "GPP_COOKIES", "", "Cookie")
	flag.StringConfigVar(&Proxy, "p", "server", "proxy", "", "Proxy Support HTTP SOCKS5")
	flag.Parse()
	log.SetLogLevel(log.LogLevelInfo)
}
