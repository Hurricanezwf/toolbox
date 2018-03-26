package main

import (
	"flag"

	"github.com/Hurricanezwf/toolbox/http"
)

var tpl string

func init() {
	flag.StringVar(&tpl, "f", "", "-f Template file to use")
	flag.Parse()
}

func main() {
	http.EnableHTTPDebug = true

	var h HTTP
	h.Load(tpl)
	h.Do()
}
