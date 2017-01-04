package athena

import (
	"fmt"
	"os"
	"testing"

	"family/conf"
)

func init() {
	prjHome := os.Getenv("GOPATH")

	e := conf.Init(prjHome)
	if e != nil {
		fmt.Println("Init error: " + e.Error())
	}
}

func TestNeedJumpToHttps(t *testing.T) {
	var userAgent string
	var isDev, httpsEnabled, isHttps, supportHttps bool

	userAgent = "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/44.0.2403.130 Safari/537.36"
	isDev = false
	httpsEnabled = true
	isHttps = false
	supportHttps = true
	got := needJumpToHttps(isDev, httpsEnabled, isHttps, supportHttps, userAgent)
	expect := true

	if got != expect {
		t.Errorf("got [%s] expected [%s]", got, expect)
	}

	isDev = true
	got = needJumpToHttps(isDev, httpsEnabled, isHttps, supportHttps, userAgent)
	expect = false
	if got != expect {
		t.Errorf("got [%s] expected [%s]", got, expect)
	}

	isDev = false
	httpsEnabled = false
	got = needJumpToHttps(isDev, httpsEnabled, isHttps, supportHttps, userAgent)
	expect = false
	if got != expect {
		t.Errorf("got [%s] expected [%s]", got, expect)
	}
}

func TestNeedJumpToHttp(t *testing.T) {
	var isDev, httpsEnabled, isHttps, supportHttps bool

	isDev = false
	httpsEnabled = true
	isHttps = true
	supportHttps = false
	got := needJumpToHttp(isDev, httpsEnabled, isHttps, supportHttps)
	expect := true

	if got != expect {
		t.Errorf("got [%s] expected [%s]", got, expect)
	}

	isDev = true
	got = needJumpToHttp(isDev, httpsEnabled, isHttps, supportHttps)
	expect = false
	if got != expect {
		t.Errorf("got [%s] expected [%s]", got, expect)
	}

	isDev = false
	httpsEnabled = false
	got = needJumpToHttp(isDev, httpsEnabled, isHttps, supportHttps)
	expect = true
	if got != expect {
		t.Errorf("got [%s] expected [%s]", got, expect)
	}
}

func TestNeedJumpToMobile(t *testing.T) {
	var requestFile, ndr, userAgent string

	requestFile = "index.html"
	userAgent = "Mozilla/5.0 (iPhone; CPU iPhone OS 8_0 like Mac OS X) AppleWebKit/600.1.3 (KHTML, like Gecko) Version/8.0 Mobile/12A4345d Safari/600.1.4"
	got := needJumpToMobile(requestFile, ndr, userAgent)
	expect := true
	if got != expect {
		t.Errorf("got [%s] expected [%s]", got, expect)
	}

	requestFile = "other.html"
	got = needJumpToMobile(requestFile, ndr, userAgent)
	expect = false
	if got != expect {
		t.Errorf("got [%s] expected [%s]", got, expect)
	}

	requestFile = "index.html"
	ndr = "random"
	got = needJumpToMobile(requestFile, ndr, userAgent)
	expect = false
	if got != expect {
		t.Errorf("got [%s] expected [%s]", got, expect)
	}
}
