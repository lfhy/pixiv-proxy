package server

import (
	"bytes"
	_ "embed"
	"go-pixiv-proxy/conf"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/lfhy/log"
	"github.com/tidwall/gjson"
)

var (
	headers = map[string]string{
		"Referer":    "https://www.pixiv.net",
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.113 Safari/537.36",
	}
	client = &http.Client{
		Transport: &http.Transport{
			Proxy: func(request *http.Request) (u *url.URL, e error) {
				return http.ProxyFromEnvironment(request)
			},
		},
	}
)

var (
	//go:embed index.html
	indexHtml     string
	directTypes   = []string{"img-original", "img-master", "c", "user-profile", "img-zip-ugoira"}
	imgTypes      = []string{"original", "regular", "small", "thumb", "mini"}
	docExampleImg = `![regular](http://example.com/98505703?t=regular)

![small](http://example.com/98505703?t=small)

![thumb](http://example.com/98505703?t=thumb)

![mini](http://example.com/98505703?t=mini)`

	//go:embed markdeep.js
	markDeepJS string
)

type Illust struct {
	origUrl string
	urls    map[string]gjson.Result
}

func handlePixivProxy(rw http.ResponseWriter, req *http.Request) {
	var err error
	var realUrl string
	c := &Context{
		rw:  rw,
		req: req,
	}
	path := req.URL.Path
	log.Info(req.Method, " ", req.URL.String())
	spl := strings.Split(path, "/")[1:]
	switch spl[0] {
	case "markdeep.js":
		MarkDeep(rw, req)
		return
	case "":
		c.String(200, indexHtml)
		return
	case "favicon.ico":
		// c.WriteHeader(404)
		http.Redirect(rw, req, "https://www.3000y.ac.cn/favicon.ico", http.StatusFound)
		return
	case "api":
		handleIllustInfo(c)
		return
	}
	imgType := req.URL.Query().Get("t")
	if imgType == "" {
		imgType = "original"
	}
	if !in(imgTypes, imgType) {
		c.String(400, "invalid query")
		return
	}
	if in(directTypes, spl[0]) {
		realUrl = "https://i.pximg.net" + path
	} else {
		if _, err = strconv.Atoi(spl[0]); err != nil {
			c.String(400, "invalid query")
			return
		}
		illust, err := getIllust(spl[0])
		if err != nil {
			c.String(400, "pixiv api error")
			return
		}
		if r, ok := illust.urls[imgType]; ok {
			realUrl = r.String()
		} else {
			c.String(400, "this image type not exists")
			return
		}
		if realUrl == "" {
			c.String(400, "this image needs login, set GPP_COOKIES env.")
		}
		if len(spl) > 1 {
			realUrl = strings.Replace(realUrl, "_p0", "_p"+spl[1], 1)
		}
	}
	info, err := url.Parse(path)
	if err != nil {
		c.String(500, err.Error())
		return
	}
	objectKey := info.Path
	t := info.Query().Get("t")
	if t != "" {
		objectKey = filepath.Join(t, objectKey)
	}
	if hasRemote {
		url, ok := HeadRemote(objectKey)
		if ok {
			http.Redirect(rw, req, url, http.StatusFound)
			return
		}
	}
	data := proxyHttpReq(c, realUrl, "fetch pixiv image error", hasRemote)
	if hasRemote && len(data) > 0 {
		go func() {
			PutDataToRemote(data, objectKey)
		}()
	}
}

func handleIllustInfo(c *Context) {
	params := strings.Split(c.req.URL.Path, "/")
	pid := params[len(params)-1]
	if _, err := strconv.Atoi(pid); err != nil {
		c.String(400, "pid invalid")
		return
	}
	proxyHttpReq(c, "https://www.pixiv.net/ajax/illust/"+pid, "pixiv api error", false)
}

func getIllust(pid string) (*Illust, error) {
	b, err := httpGetBytes("https://www.pixiv.net/ajax/illust/" + pid)
	if err != nil {
		return nil, err
	}
	g := gjson.ParseBytes(b)
	imgUrl := g.Get("body.urls.original").String()
	return &Illust{
		origUrl: imgUrl,
		urls:    g.Get("body.urls").Map(),
	}, nil
}

func in(orig []string, str string) bool {
	for _, b := range orig {
		if b == str {
			return true
		}
	}
	return false
}

type Context struct {
	rw  http.ResponseWriter
	req *http.Request
}

func (c *Context) write(b []byte, status int) {
	c.rw.WriteHeader(status)
	_, err := c.rw.Write(b)
	if err != nil {
		log.Error(err)
	}
}

func (c *Context) String(status int, s string) {
	c.write([]byte(s), status)
}

func (c *Context) WriteHeader(statusCode int) {
	c.rw.WriteHeader(statusCode)
}

func proxyHttpReq(c *Context, url string, errMsg string, buffer bool) []byte {
	resp, err := httpGet(url)
	if err != nil {
		c.String(500, errMsg)
		return nil
	}
	defer resp.Body.Close()
	copyHeader(c.rw.Header(), resp.Header)
	resp.Header.Del("Cookie")
	resp.Header.Del("Set-Cookie")
	buff := bytes.NewBuffer(nil)
	var writer io.Writer
	if buffer {
		writer = io.MultiWriter(c.rw, buff)
	} else {
		writer = c.rw
	}
	_, err = io.Copy(writer, resp.Body)
	if err != nil && buffer {
		return nil
	}
	return buff.Bytes()
}

func httpGet(u string) (*http.Response, error) {
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	req.Header.Set("Cookie", conf.Cookies)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func httpGetReadCloser(u string) (io.ReadCloser, error) {
	resp, err := httpGet(u)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func httpGetBytes(url string) ([]byte, error) {
	body, err := httpGetReadCloser(url)
	if err != nil {
		return nil, err
	}
	defer body.Close()
	b, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}
