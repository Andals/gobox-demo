package main

import (
	"flag"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strconv"
	"strings"

	"andals/gobox/http/controller"
	"andals/gobox/http/gracehttp"
	"andals/gobox/pidfile"

	"gdemo/conf"
	"gdemo/controller/api/index"
	"gdemo/errno"
	"gdemo/log"
)

func main() {
	var prjHome string

	flag.StringVar(&prjHome, "prjHome", "", "prjHome absolute path")
	flag.Parse()

	prjHome = strings.TrimRight(prjHome, "/")
	if prjHome == "" {
		fmt.Println("missing flag prjHome: ")
		flag.PrintDefaults()
		os.Exit(errno.E_SYS_INVALID_FLAG_PRJ_HOME)
	}

	e := conf.Init(prjHome)
	if e != nil {
		fmt.Println(e.Error())
		os.Exit(e.Errno())
	}

	e = log.Init(conf.ServerConf.LogRoot)
	if e != nil {
		fmt.Println(e.Error())
		os.Exit(e.Errno())
	}
	defer func() {
		log.Free()
	}()

	pf, err := pidfile.CreatePidFile(conf.ServerConf.ApiPidFile)
	if err != nil {
		fmt.Printf("create pid file %s failed, error: %s\n", conf.ServerConf.ApiPidFile, err.Error())
		os.Exit(errno.E_SYS_SAVE_PID_FILE_FAIL)
	}

	if conf.ServerConf.IsDev {
		go func() {
			http.ListenAndServe(":6060", nil)
		}()
	}

	cl := controller.NewController()
	index.RegAction(cl)

	err = gracehttp.ListenAndServe(conf.ServerConf.ApiGoHttpHost+":"+conf.ServerConf.ApiGoHttpPort, cl)
	if err != nil {
		fmt.Println("pid:" + strconv.Itoa(os.Getpid()) + ", err:" + err.Error())
	}

	if err := pidfile.ClearPidFile(pf); err != nil {
		fmt.Printf("clear pid file failed, error: %s\n", err.Error())
	}
}
