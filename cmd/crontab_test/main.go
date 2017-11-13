package main

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
	"time"
	"toolbox/crontab"
)

func init() {
	http.HandleFunc("/add", Add)
	http.HandleFunc("/del", Del)
	http.HandleFunc("/list", List)
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	crontab.Open()

	http.ListenAndServe(":10000", nil)
}

func Add(w http.ResponseWriter, r *http.Request) {
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

	log.Printf("[/add] taskname:%s, spec: %s\n", taskName, spec)

	if !crontab.SpecValid(spec) {
		w.Write([]byte("Invalid spec"))
		return
	}

	t := crontab.NewTask(taskName, spec, Do, taskName)
	if err := crontab.Add(t); err != nil {
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

	affeted := crontab.Del(taskName)
	w.Write([]byte(fmt.Sprintf("Del success, %d affeted", affeted)))
}

func List(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(405)
		return
	}

	log.Printf("[/list]")

	res := crontab.List()
	str := strings.Join(res, "\n")
	w.Write([]byte(str))
}

func Do(param interface{}) error {
	fmt.Printf("[%s] haha_%v\n", time.Now().Format("15:04:05"), param)
	return nil
}
