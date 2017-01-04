package main

import (
	"flag"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strings"

	"andals/gobox/http/controller"
	"andals/gobox/http/gracehttp"
	"andals/gobox/pidfile"

	"gdemo/conf"
	"gdemo/controller/front"
	"gdemo/errno"
	"gdemo/gvalue"

	"andals/gobox/ipquery"
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

	e = log.Init(conf.ServerConf.LogRoot())
	if e != nil {
		fmt.Println(e.Error())
		os.Exit(e.Errno())
	}
	defer func() {
		log.Free()
	}()

	e = gvalue.Init()
	if e != nil {
		fmt.Println(e.Error())
		os.Exit(e.Errno())
	}

	err := ipquery.Load(conf.ServerConf.IpDataFile)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(errno.E_SYS_LOAD_IPQUERY_DATA_FAIL)
	}

	pf, err := pidfile.CreatePidFile(conf.ServerConf.FrontPidFile())
	if err != nil {
		fmt.Printf("create pid file %s failed, error: %s\n", conf.ServerConf.FrontPidFile(), err.Error())
		os.Exit(errno.E_SYS_SAVE_PID_FILE_FAIL)
	}

	if conf.ServerConf.IsDev {
		go func() {
			http.ListenAndServe(":6060", nil)
		}()
	}

	cl := controller.NewController()
	front.RegAction(cl)

	err = gracehttp.ListenAndServe(conf.ServerConf.FrontGolangHost+":"+conf.ServerConf.FrontGolangPort, cl)
	if err != nil {
		fmt.Println(err.Error())
	}

	if err := pidfile.ClearPidFile(pf); err != nil {
		fmt.Printf("clear pid file %s failed, error: %s\n", pf.Path, err.Error())
	}
}
