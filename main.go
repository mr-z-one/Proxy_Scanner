package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"

	flag "github.com/spf13/pflag"

	"proxyScanner/utils"

	"github.com/lqqyt2423/go-mitmproxy/proxy"
	log "github.com/sirupsen/logrus"
)

//var titleRegexp = regexp.MustCompile("(<title>)(.*?)(</title>)")

// type ChangeHtml struct {
// 	proxy.BaseAddon
// }

// func (c *ChangeHtml) Response(f *proxy.Flow) {
// 	contentType := f.Response.Header.Get("Content-Type")
// 	if !strings.Contains(contentType, "text/html") {
// 		return
// 	}

//		// change html <title> end with: " - go-mitmproxy"
//		f.Response.ReplaceToDecodedBody()
//		f.Response.Body = titleRegexp.ReplaceAll(f.Response.Body, []byte("${1}go-mitmproxy${3}${2}"))
//		f.Response.Header.Set("Content-Length", strconv.Itoa(len(f.Response.Body)))
//	}
var up_stream_proxy string
var port int
var scope []string

func main() {
	parsFlag()

	opts := &proxy.Options{
		Addr:              ":" + strconv.Itoa(port),
		StreamLargeBodies: 1024 * 1024 * 10,
		SslInsecure:       true,
	}

	p, err := proxy.NewProxy(opts)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	p.SetUpstreamProxy(func(req *http.Request) (*url.URL, error) {
		proxyURL, err := url.Parse(up_stream_proxy)
		if err != nil {
			return nil, fmt.Errorf("invalid proxy URL: %w", err)
		}

		host := req.URL.Host

		var r *regexp.Regexp
		if len(scope) != 0 {
			r = utils.ScopeToDomain(scope)
		} else {
			r = regexp.MustCompile(".*")
		}

		if r.MatchString(host) {
			log.Printf("Using upstream proxy %s for %s", up_stream_proxy, host)
			return proxyURL, nil
		}

		// Returning nil with this error tells mitmproxy to not use upstream proxy
		return nil, errors.New("pass proxy")
	})

	p.Start()
	//log.Fatal(p.Start())
}

func parsFlag() {
	flag.StringVar(&up_stream_proxy, "proxy", "http://127.0.0.1:8080", "up stream proxy")
	flag.IntVar(&port, "port", 9080, "port of current proxy")
	flag.StringSliceVar(&scope, "scope", []string{}, "In scope domain")

	flag.Parse()
}
