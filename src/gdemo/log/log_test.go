package log

import (
	"fmt"
	"os"
	"testing"

	"gdemo/conf"
)

func init() {
	prjHome := os.Getenv("GOPATH")

	e := conf.Init(prjHome)
	if e != nil {
		fmt.Println("Init error: " + e.Error())
	}
}

func TestInit(t *testing.T) {
	got := Init(conf.ServerConf.LogRoot)

	if got != nil {
		t.Errorf("got [%s] expected nil", got)
	}
}
