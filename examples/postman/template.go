package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Hurricanezwf/toolbox/http"
)

/*
{
	"url": "http://localhost:14001/host",
	"method": "get",
	"header": {
		"content-type": "application/json"
	},
	"param": {
		"key1": "value1"
	},
	"body": {
		"code": 1,
		"msg": "invalid param"
	},
	"option": {
		"connect_timeout": 5,
		"readwrite_timeout": 10
	}
}
*/

type HTTPOption struct {
	ConnectTimeout   int `json:"connect_timeout"`
	ReadWriteTimeout int `json:"readwrite_timeout"`
}

type HTTP struct {
	URL    string            `json:"url"`
	Method string            `json:"method"`
	Header map[string]string `json:"header"`
	Param  map[string]string `json:"param"`
	Body   interface{}       `json:"body"`
	Option HTTPOption        `json:"option"`
}

func (h *HTTP) Load(path string) {
	if len(tpl) <= 0 {
		log.Printf("\033[31mMissing template file\033[0m\n")
		os.Exit(-1)
	}

	absPath, err := filepath.Abs(tpl)
	if err != nil {
		log.Printf("\033[31mFormat template to abs failed, %v\033[0m\n", err)
		os.Exit(-1)
	}

	data, err := ioutil.ReadFile(absPath)
	if err != nil {
		log.Printf("\033[31mRead template failed, %v\033[0m\n", err)
		os.Exit(-1)
	}

	if err = json.Unmarshal(data, &h); err != nil {
		log.Printf("\033[31mUnmarshal template failed, %v\033[0m\n", err)
		os.Exit(-1)
	}
}

func (h *HTTP) Do() {
	h.Method = strings.ToLower(h.Method)
	switch h.Method {
	case "head":
		h.head()
	case "get":
		h.get()
	case "post":
		h.post()
	case "put":
		h.put()
	case "delete":
		h.del()
	default:
		log.Printf("\033[31mUnsupport http method %s\033[0m\n", h.Method)
		os.Exit(-1)
	}
}

func (h HTTP) head() {
	option := &http.Options{
		Headers: h.Header,
		Params:  h.Param,
	}

	if h.Option.ConnectTimeout > 0 {
		option.ConnectTimeout = time.Duration(h.Option.ConnectTimeout) * time.Second
	}
	if h.Option.ReadWriteTimeout > 0 {
		option.RWTimeout = time.Duration(h.Option.ReadWriteTimeout) * time.Second
	}

	if err := http.Head(h.URL, option); err != nil {
		log.Printf("\033[31mHEAD Request err, %v\033[0m\n", err)
	} else {
		log.Printf("\033[32mHEAD Requst OK\033[0m")
	}
}

func (h HTTP) get() {
	option := &http.Options{
		Headers: h.Header,
		Params:  h.Param,
	}

	if h.Option.ConnectTimeout > 0 {
		option.ConnectTimeout = time.Duration(h.Option.ConnectTimeout) * time.Second
	}
	if h.Option.ReadWriteTimeout > 0 {
		option.RWTimeout = time.Duration(h.Option.ReadWriteTimeout) * time.Second
	}

	var resp interface{}
	if err := http.Get(h.URL, option, &resp); err != nil {
		log.Printf("\033[31mGET Request err, %v\033[0m\n", err)
	} else {
		log.Printf("\033[32mGET Requst OK, Resp: %+v\033[0m", resp)
	}
}

func (h HTTP) post() {
	option := &http.Options{
		Headers: h.Header,
		Params:  h.Param,
	}

	if h.Body != nil {
		option.Body = h.Body
	}
	if h.Option.ConnectTimeout > 0 {
		option.ConnectTimeout = time.Duration(h.Option.ConnectTimeout) * time.Second
	}
	if h.Option.ReadWriteTimeout > 0 {
		option.RWTimeout = time.Duration(h.Option.ReadWriteTimeout) * time.Second
	}

	var resp interface{}
	if err := http.Post(h.URL, option, &resp); err != nil {
		log.Printf("\033[31mPOST Request err, %v\033[0m\n", err)
	} else {
		log.Printf("\033[32mPOST Requst OK, Resp: %+v\033[0m", resp)
	}
}

func (h HTTP) put() {
	option := &http.Options{
		Headers: h.Header,
		Params:  h.Param,
	}

	if h.Body != nil {
		option.Body = h.Body
	}
	if h.Option.ConnectTimeout > 0 {
		option.ConnectTimeout = time.Duration(h.Option.ConnectTimeout) * time.Second
	}
	if h.Option.ReadWriteTimeout > 0 {
		option.RWTimeout = time.Duration(h.Option.ReadWriteTimeout) * time.Second
	}

	var resp interface{}
	if err := http.Put(h.URL, option, &resp); err != nil {
		log.Printf("\033[31mPut Request err, %v\033[0m\n", err)
	} else {
		log.Printf("\033[32mPut Requst OK, Resp: %+v\033[0m", resp)
	}
}

func (h HTTP) del() {
	option := &http.Options{
		Headers: h.Header,
		Params:  h.Param,
	}

	if h.Body != nil {
		option.Body = h.Body
	}
	if h.Option.ConnectTimeout > 0 {
		option.ConnectTimeout = time.Duration(h.Option.ConnectTimeout) * time.Second
	}
	if h.Option.ReadWriteTimeout > 0 {
		option.RWTimeout = time.Duration(h.Option.ReadWriteTimeout) * time.Second
	}

	var resp interface{}
	if err := http.Delete(h.URL, option, &resp); err != nil {
		log.Printf("\033[31mDelete Request err, %v\033[0m\n", err)
	} else {
		log.Printf("\033[32mDelete Requst OK, Resp: %+v\033[0m", resp)
	}
}
