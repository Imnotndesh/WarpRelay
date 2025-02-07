package SimpleProxy

import (
	"awesomeProject/ConfigParser"
	"awesomeProject/SimpleServer"
	"crypto/tls"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

const (
	configLocation = "./Config/config.YAML"
)

var (
	proxies = make(map[string]*httputil.ReverseProxy)
)

func StartProxy() {
	var err error
	parser := ConfigParser.ConfigParser{
		ConfigLocation: configLocation,
	}
	err = parser.ParseConfig()
	if err != nil {
		log.Println("Cannot parse config:", err)
		return
	}
	// Port from YAML
	port := parser.GetProxyPort()
	// User endpoints from YAML
	endpoints := parser.GetEndpoints()

	// ReverseProxy for each entry
	for i := range endpoints {
		var backendURL *url.URL
		endpoint := &endpoints[i]
		backendURL, err = url.Parse(endpoint.BackendUrl)
		if err != nil {
			log.Fatalf("Cannot parse url for %s: %v. Check Config!", endpoint.ProxyEndpoint, err)
			return
		}
		endpoint.ParsedUrl = backendURL
		proxy := httputil.NewSingleHostReverseProxy(backendURL)
		proxy.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		proxies[endpoint.ProxyEndpoint] = proxy
		proxy.Director = func(req *http.Request) {
			req.URL.Scheme = endpoint.ParsedUrl.Scheme
			req.URL.Host = endpoint.ParsedUrl.Host
			req.Host = endpoint.ParsedUrl.Host

			req.URL.Path = strings.TrimPrefix(req.URL.Path, endpoint.ProxyEndpoint)
			req.URL.Path = "/" + strings.TrimPrefix(req.URL.Path, "/")
			log.Printf("Forwarding request to: %s%s", endpoint.ParsedUrl, req.URL.Path)
			if len(req.URL.Path) == 0 {
				req.URL.Path = "/"
			}
			// req.Header.Set("X-Forwarded-For", req.Header.Get("X-Forwarded-For")+","+req.RemoteAddr)
		}
		proxyPath := endpoint.ProxyEndpoint
		log.Printf("Proxy set for %s -> %s", proxyPath, endpoint.BackendUrl)
		http.HandleFunc(proxyPath+"/", func(w http.ResponseWriter, r *http.Request) {
			proxies[proxyPath].ServeHTTP(w, r)
		})
	}

	// Start reverse proxy server
	log.Println("Starting proxy on port ", port)
	err = http.ListenAndServeTLS(":"+port, SimpleServer.CertLocation, SimpleServer.KeyLocation, nil)
	if err != nil {
		log.Panicln("Cannot start proxy:", err)
	}
}
