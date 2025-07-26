package server

import (
	"fmt"
	"net/http"
)

func MarkDeep(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Length", fmt.Sprint(len(markDeepJS)))
	w.Write([]byte(markDeepJS))
}
