package main

import (
	"fmt"
	"net/http"
	"net/url"
	"proxyScanner/Model"
	"proxyScanner/Server"
	"proxyScanner/config"
	"proxyScanner/dataType"
	"proxyScanner/db"
	"proxyScanner/utils"
	"regexp"
	"strconv"
	"strings"

	flag "github.com/spf13/pflag"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/lqqyt2423/go-mitmproxy/proxy"
	log "github.com/sirupsen/logrus"
)

//var titleRegexp = regexp.MustCompile("(<title>)(.*?)(</title>)")

type seeHttp struct {
	proxy.BaseAddon
}

func (c *seeHttp) Requestheaders(f *proxy.Flow) {
	if !passProxy(scope, outScope, f.Request.Raw().Host, f.Request.Raw().URL.Path) {
		return
	}

	path := f.Request.URL.Path
	if strings.Contains(path, ".js") {

		modified := f.Request.Header.Get("If-Modified-Since")
		if modified != "" {
			fmt.Printf("Modified since %s\n", modified)
			f.Request.Header.Del("If-Modified-Since")
		}

		modified = f.Request.Header.Get("If-None-Match")
		if modified != "" {
			fmt.Printf("If-None-Match %s\n", modified)
			f.Request.Header.Del("If-None-Match")
		}
	}
}

func (c *seeHttp) Request(f *proxy.Flow) {
	if !passProxy(scope, outScope, f.Request.Raw().Host, f.Request.Raw().URL.Path) {

		//f.Request.URL.Host = "sentry.namava"
		//f.Request.URL.Path = "/redeacated.com"
		return
	}

}

func (c *seeHttp) Response(f *proxy.Flow) {
	if !passProxy(scope, outScope, f.Request.Raw().Host, f.Request.Raw().URL.Path) {
		return
	}

	parse_error := f.Request.Raw().ParseForm()
	params := dataType.JSONMap{}
	if parse_error == nil {

		params = utils.ValuesToJSONMap(f.Request.Raw().Form)
	} else {
		params = nil
	}

	HandleNullUnicode(f)

	data := Model.HttpRequest{

		Path:   f.Request.URL.RequestURI(),
		Method: f.Request.Method,
		Host:   f.Request.Raw().Host,

		RequestHeaders: utils.HeaderToJSONMap(f.Request.Header),
		RequestBodyRaw: string(f.Request.Body),

		Parameters: params,

		StatusCode:      f.Response.StatusCode,
		ResponseHeaders: utils.HeaderToJSONMap(f.Response.Header),
		ResponseBodyRaw: string(f.Response.Body),
	}
	database.Clauses(clause.OnConflict{DoNothing: true}).Create(&data)
	//
	//fmt.Println("--------------------------REQHEADER-----------------------------------")
	//
	//fmt.Println("Request : ", f.Request.Method, f.Request.URL, f.Request.Raw().Form, " ", f.Request.Header)
	//fmt.Println("--------------------------REQBody-----------------------------------")
	//fmt.Println(string(f.Request.Body))
	//fmt.Println("--------------------------RESHEADER-----------------------------------")
	//fmt.Println("Response : ", f.Response.StatusCode, " ", f.Response.Header)
	//fmt.Println("--------------------------RESbody-----------------------------------")
	//fmt.Println("Response : ", string(f.Response.Body))
	//fmt.Println("-------------------------------------------------------------")
}

func HandleNullUnicode(f *proxy.Flow) {
	if strings.Contains(string(f.Response.Body), "\\u0000") {

		rb := string(f.Response.Body)
		n_rb := strings.Replace(rb, "\\u0000", "\\u1111", -1)
		f.Response.Body = []byte(n_rb)
	}
	if strings.Contains(string(f.Request.Body), "\\u0000") {

		rb := string(f.Request.Body)
		n_rb := strings.Replace(rb, "\\u0000", "\\u1111", -1)
		f.Request.Body = []byte(n_rb)
	}
}

var database *gorm.DB
var up_stream_proxy string
var port int
var scope []string
var outScope []string

var username string
var password string
var databaseName string

var jsonConfig string

var incScopeRegex, outScopeRegex *regexp.Regexp

func main() {

	//rrr := regexp.MustCompile(`(sentry\.namava|/api/v1\.0/medias/\d+/play-info)`)
	//fmt.Println(rrr.MatchString("/api/v1.0/medias/259743/play-info"))
	//return

	parsFlag()
	fetchConfigFile()
	showInputData()
	CreateRegex()
	opts := &proxy.Options{
		Addr:              ":" + strconv.Itoa(port),
		StreamLargeBodies: 1024 * 1024 * 5,
		SslInsecure:       true,
		Debug:             0,
	}

	p, err := proxy.NewProxy(opts)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	database, err = db.GetDatabaseConnection(username, password, databaseName)
	if err != nil {
		fmt.Println("[-] failed to connect to database :", err)
		return
	} else {
		fmt.Println("[+] connect to database with username : ", username, "db : ", databaseName)
	}

	database.AutoMigrate(&Model.HttpRequest{})
	fmt.Println("[+] table created..")
	go Server.StartServer(2031)
	//p.SetShouldInterceptRule(func(req *http.Request) bool {
	//	path := req.URL.Path
	//	host := req.URL.Host
	//
	//	if req.Method != "CONNECT" {
	//		fmt.Println("no CONNECT")
	//	}
	//	if strings.Contains(host, "namava") {
	//		fmt.Println("okkkkkk")
	//	}
	//
	//	if strings.Contains(path, "/api/v1.0/medias/") {
	//		fmt.Println("NAokkkkkk")
	//	}
	//
	//	return passProxy(scope, outScope, host, path)
	//})

	p.SetUpstreamProxy(func(req *http.Request) (*url.URL, error) {
		proxyURL, err := url.Parse(up_stream_proxy)
		if err != nil {
			return nil, fmt.Errorf("[-] invalid proxy URL: %w", err)
		}

		host := req.URL.Host

		if passProxy(scope, outScope, host, "") {
			return proxyURL, nil
		} else {
			return nil, nil
		}

	})

	p.AddAddon(&seeHttp{})

	p.Start()
	//log.Fatal(p.Start())
}

func CreateRegex() {
	incScopeRegex = utils.ScopeToDomainRegex(scope)
	outScopeRegex = utils.ScopeToDomainRegex(outScope)
}

func fetchConfigFile() {

	if jsonConfig == "" {
		return
	}
	c := config.ReadeConfig(jsonConfig)
	if c == nil {
		return
	}
	if username == "" {
		username = c.Username
	}
	if password == "" {
		password = c.Password
	}
	if databaseName == "" {
		databaseName = c.DbName
	}
	if port == 0 {
		port = c.Port
	}

	if up_stream_proxy == "" {
		up_stream_proxy = c.Proxy
	}

	if len(scope) == 0 {
		scope = c.Scope
	}
	if len(outScope) == 0 {
		outScope = c.OutScope
	}
}

func passProxy(inScope []string, outScope []string, host string, path string) bool {
	var r *regexp.Regexp

	r = outScopeRegex

	if len(outScope) > 0 && (r.MatchString(host) || r.MatchString(path)) {
		// Returning nil with this error tells mitmproxy to not use upstream proxy
		return false
	}

	if len(scope) != 0 {
		r = incScopeRegex
	} else {
		return true
	}

	if r.MatchString(host) {
		log.Printf("Using upstream proxy %s for %s", up_stream_proxy, host)
		return true
	}

	return false
}

func showInputData() {
	fmt.Println("[+] username:", username)
	fmt.Println("[+] dbName:", databaseName)
	fmt.Println("[+] In Scope Domain", scope)
	fmt.Println("[+] out Scope Domain", outScope)

}

func parsFlag() {
	flag.StringVar(&up_stream_proxy, "proxy", "", "up stream proxy")
	flag.IntVar(&port, "port", 0, "port of current proxy")
	flag.StringSliceVar(&scope, "scope", []string{}, "In scope domain")
	flag.StringSliceVar(&outScope, "out-scope", []string{}, "out scope domain")

	flag.StringVar(&username, "username", "", "username of pg database")
	flag.StringVar(&password, "password", "", "password of pg database")
	flag.StringVar(&databaseName, "db", "", "db name")

	flag.StringVar(&jsonConfig, "config", "", "json config file")
	flag.Parse()

}
