package conf

import (
	"fmt"
	"os"
	"os/user"
	"testing"
)

var prjHomeTest string

func init() {
	prjHomeTest = os.Getenv("GOPATH")

	e := Init(prjHomeTest)
	if e != nil {
		fmt.Println("Init error: " + e.Error())
	}
}

func TestPrjHome(t *testing.T) {
	got := ServerConf.PrjHome
	expect := prjHomeTest

	if got != expect {
		t.Errorf("got [%s] expected [%s]", got, expect)
	}
}

func TestHostname(t *testing.T) {
	got := ServerConf.Hostname
	expect, _ := os.Hostname()

	if got != expect {
		t.Errorf("got [%s] expected [%s]", got, expect)
	}
}

func TestUsername(t *testing.T) {
	got := ServerConf.Username
	curUser, _ := user.Current()
	expect := curUser.Username

	if got != expect {
		t.Errorf("got [%s] expected [%s]", got, expect)
	}
}

func TestPrjName(t *testing.T) {
	got := ServerConf.PrjName
	expect := "gdemo"

	if got != expect {
		t.Errorf("got [%s] expected [%s]", got, expect)
	}
}

func TestDataRoot(t *testing.T) {
	got := ServerConf.DataRoot
	expect := prjHomeTest + "/data"

	if got != expect {
		t.Errorf("got [%s] expected [%s]", got, expect)
	}
}

func TestLogRoot(t *testing.T) {
	got := ServerConf.LogRoot
	expect := prjHomeTest + "/logs"

	if got != expect {
		t.Errorf("got [%s] expected [%s]", got, expect)
	}
}

func TestTmpRoot(t *testing.T) {
	got := ServerConf.TmpRoot
	expect := prjHomeTest + "/tmp"

	if got != expect {
		t.Errorf("got [%s] expected [%s]", got, expect)
	}
}

func TestIsDev(t *testing.T) {
	got := ServerConf.IsDev

	if got != true && got != false {
		t.Errorf("got [%s] not bool", got)
	}
}

func TestApiGoHttpHost(t *testing.T) {
	got := ServerConf.ApiGoHttpHost

	if got != "127.0.0.1" {
		t.Errorf("got ApiGoHttpHost error")
	}
}

func TestApiGoHttpPort(t *testing.T) {
	got := ServerConf.ApiGoHttpPort

	if got == "" {
		t.Errorf("got ApiGoHttpPort empty")
	}
}

func TestApiPidFile(t *testing.T) {
	got := ServerConf.ApiPidFile
	expect := prjHomeTest + "/tmp/api.pid"

	if got != expect {
		t.Errorf("got [%s] expected [%s]", got, expect)
	}
}
