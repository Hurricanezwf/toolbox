// http包简化封装了http的几种操作，仅适用于返回结构为json的情况
package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/astaxie/beego/httplib"
)

type Options struct {
	Headers        map[string]string
	Body           interface{} // Body表示struct结构,上层无须Marshal
	Params         map[string]string
	ConnectTimeout time.Duration
	RWTimeout      time.Duration
}

func Head(url string, option *Options) error {
	return send(httplib.Head(url), option, nil)
}

func Get(url string, option *Options, result interface{}) error {
	return send(httplib.Get(url), option, result)
}

func Post(url string, option *Options, result interface{}) error {
	return send(httplib.Post(url), option, result)
}

func Put(url string, option *Options, result interface{}) error {
	return send(httplib.Put(url), option, result)
}

func Delete(url string, option *Options, result interface{}) error {
	return send(httplib.Delete(url), option, result)
}

func send(req *httplib.BeegoHTTPRequest, option *Options, result interface{}) error {
	if req == nil {
		return fmt.Errorf("Http request instance is nil")
	}

	var (
		err error
		ok  bool
		b   []byte
	)

	if option != nil {
		// set headers
		for k, v := range option.Headers {
			req.Header(k, v)
		}

		// set params
		for k, v := range option.Params {
			req.Param(k, v)
		}

		// set rwtimeout
		if option.ConnectTimeout < time.Second || option.RWTimeout < time.Second {
			option.ConnectTimeout = 10 * time.Second
			option.RWTimeout = 10 * time.Second
		}
		req.SetTimeout(option.ConnectTimeout, option.RWTimeout)

		// set body
		if option.Body != nil {
			if b, ok = option.Body.([]byte); !ok {
				b, err = json.Marshal(option.Body)
				if err != nil {
					return fmt.Errorf("Marshal body failed, %v", err)
				}
			}
			if len(b) > 0 {
				req.Body(b)
			}
		}
	}

	resp, err := req.Response()
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("StatusCode(%d) != 200", resp.StatusCode)
	}

	if result != nil {
		if b, err = ioutil.ReadAll(resp.Body); err != nil {
			return fmt.Errorf("Read http body failed, %v", err)
		}
		if err = json.Unmarshal(b, result); err != nil {
			return fmt.Errorf("Bad response format, %v", err)
		}
	}
	return nil
}
