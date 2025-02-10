package SimpleProxy

import (
	"awesomeProject/ConfigParser"
	"crypto/tls"
	"fmt"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path/filepath"
	"strings"
)

const (
	ConfigDir      = "Config"
	ConfigLocation = "./Config/config.YAML"
	logsDir        = "Logs"
	CertLocation   = "./Certs/server.crt"
	KeyLocation    = "./Certs/server.key"
)

var (
	proxies = make(map[string]*httputil.ReverseProxy)
)

func SetupLogger() {
	log.SetOutput(&lumberjack.Logger{
		Filename:   filepath.Join(logsDir, "proxy.log"),
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   true,
	})
}
func StartProxy(customConfigPath string) {
	fmt.Println("\n -> Logger Initialized in logs directory")
	log.Println("Starting Reverse Proxy Server")
	var (
		err                           error
		configPath, certPath, keyPath string
	)
	if customConfigPath != "" {
		configPath = customConfigPath
	} else {
		configPath = ConfigLocation
	}
	parser := ConfigParser.ConfigParser{
		ConfigLocation: configPath,
	}
	err = parser.ParseConfig()
	if err != nil {
		log.Println("Cannot parse config:", err)
		return
	}
	// Port from YAML
	port := parser.GetProxyPort()
	certPath, keyPath = parser.GetCertInfo()
	if certPath == "" || keyPath == "" {
		certPath = CertLocation
		keyPath = KeyLocation
	}
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

		proxy.ModifyResponse = func(resp *http.Response) error {
			if resp.StatusCode >= 300 && resp.StatusCode < 400 {
				location, locErr := resp.Location()
				if locErr == nil {
					originalLocation := location.String()
					newLocation := endpoint.ProxyEndpoint + location.Path
					log.Printf("Rewriting redirect: %s -> %s", originalLocation, newLocation)
					resp.Header.Set("Location", newLocation)
				}
			}
			return nil
		}

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
	fmt.Println(" -> Running Server at port :" + port)
	// Start reverse proxy server
	err = http.ListenAndServeTLS(":"+port, certPath, keyPath, nil)
	if err != nil {
		log.Panicln("Cannot start proxy:", err)
	}
}
