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
	got := PrjHome
	expect := prjHomeTest

	if got != expect {
		t.Errorf("got [%s] expected [%s]", got, expect)
	}
}

func TestHostname(t *testing.T) {
	got := Hostname
	expect, _ := os.Hostname()

	if got != expect {
		t.Errorf("got [%s] expected [%s]", got, expect)
	}
}

func TestUsername(t *testing.T) {
	got := Username
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

func TestFrontGoHttpPort(t *testing.T) {
	got := ServerConf.FrontGoHttpPort

	if got == "" {
		t.Errorf("got FrontGoHttpPort empty")
	}
}

func TestFrontPidFile(t *testing.T) {
	got := ServerConf.FrontPidFile
	expect := prjHomeTest + "/tmp/front.pid"

	if got != expect {
		t.Errorf("got [%s] expected [%s]", got, expect)
	}
}
