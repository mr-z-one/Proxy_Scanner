package main

import (
	"fmt"
	"net/http"
	"net/url"
	"proxyScanner/config"
	"proxyScanner/db"
	"proxyScanner/utils"
	"regexp"
	"strconv"

	flag "github.com/spf13/pflag"

	"github.com/lqqyt2423/go-mitmproxy/proxy"
	log "github.com/sirupsen/logrus"
)

//var titleRegexp = regexp.MustCompile("(<title>)(.*?)(</title>)")

type seeHttp struct {
	proxy.BaseAddon
}

func (c *seeHttp) Response(f *proxy.Flow) {
	if !passProxy(scope, outScope, f.Request.Raw().Host) {
		return
	}
	fmt.Println("--------------------------REQHEADER-----------------------------------")
	fmt.Println("Request : ", f.Request.Header)
	fmt.Println("--------------------------REQBody-----------------------------------")
	fmt.Println(string(f.Request.Body))
	fmt.Println("--------------------------RESHEADER-----------------------------------")
	fmt.Println("Response : ", f.Response.Header)
	fmt.Println("--------------------------RESbody-----------------------------------")
	fmt.Println("Response : ", string(f.Response.Body))
	fmt.Println("-------------------------------------------------------------")
}

var up_stream_proxy string
var port int
var scope []string
var outScope []string

var username string
var password string
var databaseName string

var jsonConfig string

func main() {
	parsFlag()
	fetchConfigFile()
	showInputData()
	opts := &proxy.Options{
		Addr:              ":" + strconv.Itoa(port),
		StreamLargeBodies: 1024 * 1024 * 10,
		SslInsecure:       true,
		Debug:             0,
	}

	p, err := proxy.NewProxy(opts)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	db, err := db.GetDatabaseConnection(username, password, databaseName)
	if err != nil {
		fmt.Println("[-] failed to connect to database :", err)
		return
	} else {
		fmt.Println("[+] connect to database with username : ", username, "db : ", databaseName)
	}
	fmt.Println(db)

	p.SetUpstreamProxy(func(req *http.Request) (*url.URL, error) {
		proxyURL, err := url.Parse(up_stream_proxy)
		if err != nil {
			return nil, fmt.Errorf("[-] invalid proxy URL: %w", err)
		}

		host := req.URL.Host

		if passProxy(scope, outScope, host) {
			return proxyURL, nil
		} else {
			return nil, nil
		}

	})
	p.AddAddon(&seeHttp{})
	p.Start()
	//log.Fatal(p.Start())
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

func passProxy(inScope []string, outScope []string, host string) bool {
	var r *regexp.Regexp
	r = utils.ScopeToDomainRegex(outScope)

	if len(outScope) > 0 && r.MatchString(host) {
		// Returning nil with this error tells mitmproxy to not use upstream proxy
		return false
	}

	if len(scope) != 0 {
		r = utils.ScopeToDomainRegex(scope)
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
