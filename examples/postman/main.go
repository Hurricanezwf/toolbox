package main

import "flag"

var tpl string

func init() {
	flag.StringVar(&tpl, "f", "", "-f Template file to use")
	flag.Parse()
}

func main() {
	var h HTTP
	h.Load(tpl)
	h.Do()
}
