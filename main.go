package main

import (
	"net/http"
	"net/url"
	"regexp"

	"github.com/lqqyt2423/go-mitmproxy/proxy"
	log "github.com/sirupsen/logrus"
)

var titleRegexp = regexp.MustCompile("(<title>)(.*?)(</title>)")

// type ChangeHtml struct {
// 	proxy.BaseAddon
// }

// func (c *ChangeHtml) Response(f *proxy.Flow) {
// 	contentType := f.Response.Header.Get("Content-Type")
// 	if !strings.Contains(contentType, "text/html") {
// 		return
// 	}

// 	// change html <title> end with: " - go-mitmproxy"
// 	f.Response.ReplaceToDecodedBody()
// 	f.Response.Body = titleRegexp.ReplaceAll(f.Response.Body, []byte("${1}go-mitmproxy${3}${2}"))
// 	f.Response.Header.Set("Content-Length", strconv.Itoa(len(f.Response.Body)))
// }

func main() {
	opts := &proxy.Options{
		Addr:              ":9080",
		StreamLargeBodies: 1024 * 1024 * 10,
		SslInsecure:       true,
	}

	p, err := proxy.NewProxy(opts)
	if err != nil {
		log.Fatal(err)
	}
	p.SetUpstreamProxy(func(req *http.Request) (*url.URL, error) {
		url, err := url.Parse("http://127.0.0.1:8080")
		host := req.URL.Host

		r, _ := regexp.Compile("sku.ac.ir")
		if r.MatchString(host) {

			return url, err
		}
		return nil, nil
	})

	p.Start()
	//log.Fatal(p.Start())
}
