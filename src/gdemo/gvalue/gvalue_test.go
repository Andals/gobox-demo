package gvalue

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"family/conf"
	"family/log"
)

func init() {
	prjHome := os.Getenv("GOPATH")

	e := conf.Init(prjHome)
	if e != nil {
		fmt.Println("Init error: " + e.Error())
	}

	logRoot := conf.ServerConf.LogRoot()
	log.Init(logRoot)

	Init()
}

func TestGetHuajiaoLiveDataFromApi(t *testing.T) {
	got, err := GetHuajiaoLiveDataFromApi()
	expectSubstr := "<script>var HUAJIAO_LIVE_DATA = "

	if err != nil {
		t.Fatalf(err.Error())
	}

	gotStr := string(got)
	if gotStr == "" {
		t.Fatalf("getHuajiaoLiveDataFromApi fail got empty")
	}

	if !strings.Contains(gotStr, expectSubstr) {
		t.Errorf("getHuajiaoLiveDataFromApi fail got[%s]", got)
	}
}

func TestGetAthenaPageFromFile(t *testing.T) {
	var fileName string

	fileName = "index.html"
	got, err := GetAthenaPageFromFile(fileName)

	if err != nil {
		t.Fatalf(err.Error())
	}

	gotStr := string(got)
	if gotStr == "" {
		t.Fatalf("getAthenaPageFromFile fail got empty")
	}

	fileName = "test.html"
	got, _ = GetAthenaPageFromFile(fileName)
	gotStr = string(got)
	if gotStr != "" {
		t.Fatalf("getAthenaPageFromFile fail got [%s] expect empty", got)
	}
}
