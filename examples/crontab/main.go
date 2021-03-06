package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Hurricanezwf/toolbox/crontab"
	"github.com/Hurricanezwf/toolbox/crontab/scheduler"
	"github.com/Hurricanezwf/toolbox/crontab/scheduler/types"
)

var cron = crontab.New()

func init() {
	http.HandleFunc("/add/basic", AddBasic)
	http.HandleFunc("/add/pro", AddPro)
	http.HandleFunc("/del", Del)
	http.HandleFunc("/list", List)
}

func main() {
	http.ListenAndServe(":10000", nil)
}

func AddBasic(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(405)
		return
	}

	var (
		taskName string
		spec     string
	)

	r.ParseForm()
	taskName = r.FormValue("taskname")
	spec = r.FormValue("spec")

	log.Printf("[/add/basic] taskname:%s, spec: %s\n", taskName, spec)

	sche, err := scheduler.GetParser(types.ParserBasic).Parse(spec)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Invalid spec"))
		return
	}

	t := crontab.NewTask(taskName, spec, sche, func() error {
		return Do(taskName)
	})
	if err := cron.Add(t); err != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte("Add success"))
}

func AddPro(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(405)
		return
	}

	var (
		taskName string
		spec     string
	)

	r.ParseForm()
	taskName = r.FormValue("taskname")
	spec = r.FormValue("spec")

	log.Printf("[/add/pro] taskname:%s, spec: %s\n", taskName, spec)

	sche, err := scheduler.GetParser(types.ParserPro).Parse(spec)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Invalid spec"))
		return
	}

	t := crontab.NewTask(taskName, spec, sche, func() error {
		return Do(taskName)
	})
	if err := cron.Add(t); err != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte("Add success"))
}

func Del(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(405)
		return
	}

	var taskName string

	r.ParseForm()
	taskName = r.FormValue("taskname")

	log.Printf("[/del] taskname: %s\n", taskName)

	affeted, err := cron.Del(taskName)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte(fmt.Sprintf("Del success, %d affeted", affeted)))
}

func List(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(405)
		return
	}

	log.Printf("[/list]")

	res, err := cron.List()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	str := strings.Join(res, "\n")
	w.Write([]byte(str))
}

func Do(param interface{}) error {
	fmt.Printf("[%s] haha_%v\n", time.Now().Format("15:04:05"), param)
	return nil
}
