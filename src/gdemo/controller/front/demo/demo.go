package athena

import (
	"bytes"
	"html"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"andals/gobox/crypto"
	"andals/gobox/http/controller"
	glog "andals/gobox/log"

	"family/conf"
	"family/log"

	"family/controller/front/base"
	sathena "family/svc/athena"
	"family/svc/huid"
)

const (
	JUMP_URL_HTTP   = "http://hao.360.cn"
	JUMP_URL_HTTPS  = "https://hao.360.cn"
	JUMP_URL_MOBILE = "http://m.3600.com/?a0012"
)

const (
	INDEX_FILE_NAME = "index.html"
)

const (
	COOKIE_NAME_HUID      = "__huid"
	COOKIE_NAME_SID       = "__hsid"
	COOKIE_PATH           = "/"
	COOKIE_DOMAIN         = "360.cn"
	COOKIE_MAX_AGE        = 86400
	COOKIE_INDEX_JUMP     = "360WEBINDEXCK"
	COOKIE_LM_INDEX_JUMP  = "LM_WEBINDEX"
	COOKIE_SKIP_JUMP_LOOP = "skipnum"

	SID_LEN = 16
)

const (
	JUMP_URL_PERSONAL = "http://hao.360.cn/i/index.html"
	JUMP_URL_KUANPING = "http://hao.360.cn/kuanping.html"
)

var JUMP_WHITELIST = map[string]int{
	JUMP_URL_PERSONAL: 1,
	JUMP_URL_KUANPING: 1,
}

var serverDataReplaceOld = []byte("<!--{poseidon app=server_data part=data}-->")
var huajiaoLiveDataReplaceOld = []byte("<!--{poseidon app=huajiao_live_data part=data}-->")

func RegAction(cl *controller.Controller) {
	for requestUri, _ := range conf.PageConf.AthenaRedirectConf {
		cl.AddExactMatchAction(requestUri, AthenaRedirect)
	}

	for requestUri, _ := range conf.PageConf.AthenaPageConf {
		cl.AddExactMatchAction(requestUri, AthenaPage)
	}
}

func AthenaRedirect(context *controller.Context, args []string) {
	var redirectUri string
	athenaRedirectConf := conf.PageConf.AthenaRedirectConf

	redirectConfItem, ok := athenaRedirectConf[context.Req.URL.Path]
	if ok {
		redirectUri = redirectConfItem.RedirectUri
	} else {
		redirectUri = "/"
	}

	controller.Redirect302(redirectUri)
}

func AthenaPage(context *controller.Context, args []string) {
	athenaPageConf := conf.PageConf.AthenaPageConf
	requestUri := context.Req.URL.Path

	pageConfItem, ok := athenaPageConf[requestUri]
	if !ok {
		context.RespBody = []byte(nil)
		return
	}

	if pageConfItem.Alias != "" {
		requestUri = pageConfItem.Alias
		pageConfItem = athenaPageConf[pageConfItem.Alias]
	}

	var ip string
	rip, ok := context.TransData[sathena.TRANS_DATA_KEY_REAL_IP]
	if ok {
		realIp, ok := rip.(string)
		if ok {
			ip = realIp
		}
	}

	isHttps := base.IsHttps(context)
	processHttpJump(context, isHttps, pageConfItem.SupportHttps)
	processHttpsJump(context, isHttps, pageConfItem.SupportHttps, ip)

	requestFile := strings.TrimLeft(requestUri, "/")
	processMobileJump(context, requestFile)
	processPcJump(context, requestFile)

	athenaPage := getAthenaPage(context, requestFile, isHttps)

	if requestFile == INDEX_FILE_NAME {
		processUserData(context, ip)
		serverData := sathena.GetServerData(ip)
		huajiaoData := sathena.GetHuajiaoLiveData(context.TransData)

		athenaPage = bytes.Replace(athenaPage, serverDataReplaceOld, serverData, 1)
		athenaPage = bytes.Replace(athenaPage, huajiaoLiveDataReplaceOld, huajiaoData, 1)
	}

	context.RespBody = []byte(athenaPage)
}

func processHttpJump(context *controller.Context, isHttps, supportHttps bool) {
	if needJumpToHttp(conf.ServerConf.IsDev, conf.ServerConf.HttpsServiceEnabled, isHttps, supportHttps) {
		controller.Redirect302(JUMP_URL_HTTP + context.Req.RequestURI)
	}
}

func needJumpToHttp(isDev, httpsEnabled, isHttps, supportHttps bool) bool {
	if isDev {
		return false
	}

	if httpsEnabled {
		if isHttps && !supportHttps {
			return true
		}
	} else {
		if isHttps {
			return true
		}
	}

	return false
}

func processHttpsJump(context *controller.Context, isHttps, supportHttps bool, ip string) {
	if needJumpToHttps(conf.ServerConf.IsDev, conf.ServerConf.HttpsServiceEnabled, isHttps, supportHttps, context.Req.UserAgent()) {
		var logger glog.ILogger
		logger, e := glog.NewSimpleLogger(log.HttpsJumpLogWriter, glog.LEVEL_INFO, glog.NewWebFormater(context.Rid, []byte(ip)))
		if e != nil {
			logger = new(glog.NoopLogger)
		}

		logger.Info([]byte(context.Req.URL.Path))
		logger.Free()

		controller.Redirect302(JUMP_URL_HTTPS + context.Req.RequestURI)
	}
}

func needJumpToHttps(isDev, httpsEnabled, isHttps, supportHttps bool, userAgent string) bool {
	if !isDev && httpsEnabled && !isHttps && supportHttps && isUaSupportHttps(userAgent) {
		return true
	}
	return false
}

func processMobileJump(context *controller.Context, requestFile string) {
	ndr := html.EscapeString(context.Req.FormValue("ndr"))
	if needJumpToMobile(requestFile, ndr, context.Req.UserAgent()) {
		controller.Redirect302(JUMP_URL_MOBILE)
	}
}

func needJumpToMobile(requestFile, ndr, userAgent string) bool {
	if requestFile != INDEX_FILE_NAME {
		return false
	}

	if ndr != "" || userAgent == "" {
		return false
	}

	reg := regexp.MustCompile(`(Android.+Mobile|iPhone)`)
	if reg.Match([]byte(userAgent)) {
		return true
	}
	return false
}

func processPcJump(context *controller.Context, requestFile string) {
	if requestFile != INDEX_FILE_NAME {
		return
	}

	src := html.EscapeString(context.Req.FormValue("src"))
	lmCookie, lmErr := context.Req.Cookie(COOKIE_LM_INDEX_JUMP)

	if src == "lm" && lmErr != nil {
		c := newCookie(COOKIE_INDEX_JUMP, "")
		http.SetCookie(context.RespWriter, c)
		return
	}

	indexCookie, indexErr := context.Req.Cookie(COOKIE_INDEX_JUMP)
	var jumpUrl string
	if lmErr == nil {
		jumpUrl = lmCookie.Value
	}
	if indexErr == nil && indexCookie.Value != "" {
		jumpUrl = indexCookie.Value
	}
	if jumpUrl == "" {
		return
	}

	_, ok := JUMP_WHITELIST[jumpUrl]
	if !ok {
		indexC := newCookie(COOKIE_INDEX_JUMP, "")
		http.SetCookie(context.RespWriter, indexC)

		lmC := newCookie(COOKIE_LM_INDEX_JUMP, "")
		http.SetCookie(context.RespWriter, lmC)
		return
	}

	queryStr := getReqQueryStr(context)
	if queryStr != "" {
		jumpUrl += "?" + queryStr
	}

	skipNumCookie, _ := context.Req.Cookie(COOKIE_SKIP_JUMP_LOOP)
	curUrl := context.Req.Host + context.Req.RequestURI
	skipNum, _ := strconv.Atoi(skipNumCookie.Value)
	if skipNum < 3 && !strings.Contains(jumpUrl, curUrl) {
		skipNum++
		c := &http.Cookie{
			Name:   COOKIE_SKIP_JUMP_LOOP,
			Value:  strconv.Itoa(skipNum),
			Path:   COOKIE_PATH,
			Domain: COOKIE_DOMAIN,
			MaxAge: int(time.Now().Unix() + 1),
		}
		http.SetCookie(context.RespWriter, c)

		controller.Redirect302(jumpUrl)
	}
}

func getAthenaPage(context *controller.Context, requestFile string, isHttps bool) []byte {
	var athenaPage []byte
	if conf.ServerConf.IsDev {
		url := "http://" + context.Req.Host + "/athena/" + requestFile
		queryStr := getReqQueryStr(context)
		if queryStr != "" {
			url += "?" + queryStr
		}
		athenaPage = sathena.GetDevAthenaPage(url, isHttps, conf.ServerConf.HttpsServiceEnabled, conf.ServerConf.QhbackServiceEnabled())
	} else {
		athenaPage = sathena.GetOnlineAthenaPage(requestFile, isHttps, conf.ServerConf.HttpsServiceEnabled, conf.ServerConf.QhbackServiceEnabled(), context.TransData)
	}

	return athenaPage
}

func getReqQueryStr(context *controller.Context) string {
	context.Req.ParseForm()
	form := context.Req.Form
	queryStr := form.Encode()

	return queryStr
}

func processUserData(context *controller.Context, ip string) {
	c, e := context.Req.Cookie(COOKIE_NAME_HUID)
	if e != nil || c.Value == "" {
		c = newCookie(COOKIE_NAME_HUID, huid.GenHuid(time.Now(), ip, context.Req.UserAgent(), huid.HUID_VERSION, huid.HUID_PRODUCER))
		http.SetCookie(context.RespWriter, c)
	}

	sid := crypto.Md5(context.Rid)[0:SID_LEN]
	c = newCookie(COOKIE_NAME_SID, string(sid))
	http.SetCookie(context.RespWriter, c)

	l, ok := context.TransData[sathena.TRANS_DATA_KEY_USER_LOG]
	if ok {
		logger, ok := l.(glog.ILogger)
		if ok {
			logger.Info(sid)
		}
	}
}

func newCookie(name, value string) *http.Cookie {
	return &http.Cookie{
		Name:   name,
		Value:  value,
		Path:   COOKIE_PATH,
		Domain: COOKIE_DOMAIN,
		MaxAge: COOKIE_MAX_AGE,
	}
}

func isUaSupportHttps(userAgent string) bool {
	if strings.Contains(userAgent, "MISE6.0") {
		return false
	}
	return true
}
