package conf

import (
	"github.com/lfhy/flag"
	"github.com/lfhy/log"
)

var (
	Host           string
	Port           string
	Domain         string
	Cookies        string
	Proxy          string
	RemoteEndpoint string
	RemoteAk       string
	RemoteSk       string
	RemoteBucket   string
	RemoteDir      string
)

func init() {
	flag.StringFullVar(&Host, "h", "server", "host", "GPP_HOST", "0.0.0.0", "Host")
	flag.StringFullVar(&Port, "p", "server", "port", "GPP_PORT", "18090", "Port")
	flag.StringFullVar(&Domain, "d", "server", "domain", "GPP_DOMAIN", "", "Server Domain")
	flag.StringFullVar(&Cookies, "c", "server", "cookies", "GPP_COOKIES", "", "Cookie")
	flag.StringConfigVar(&Proxy, "proxy", "server", "proxy", "", "Proxy Support HTTP SOCKS5")
	flag.StringConfigVar(&RemoteEndpoint, "s3-endpoint", "s3", "endpoint", "https://fgws3-ocloud.ihep.ac.cn", "S3 Endpoint")
	flag.StringConfigVar(&RemoteAk, "s3-ak", "s3", "ak", "", "S3 Ak")
	flag.StringConfigVar(&RemoteSk, "s3-sk", "s3", "sk", "", "S3 Sk")
	flag.StringConfigVar(&RemoteBucket, "s3-bucket", "s3", "bucket", "", "S3 Bucket")
	flag.StringConfigVar(&RemoteDir, "s3-dir", "s3", "dir", "image/ai/pixiv", "S3 Dir")
	flag.Parse()
	log.SetLogLevel(log.LogLevelInfo)
}
