package main

import (
	"flag"
	"go-pixiv-proxy/server"
)

func main() {
	flag.Parse()
	server.Run()
}
