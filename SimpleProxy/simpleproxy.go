package SimpleProxy

import (
	"awesomeProject/ConfigParser"
	"awesomeProject/SimpleServer"
	"crypto/tls"
	"fmt"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

const (
	configDir      = "Config"
	configLocation = "./Config/config.YAML"
	logsDir        = "Logs"
)

func init() {
	// Check if Config Dir exists
	if _, err := os.Stat(configLocation); os.IsNotExist(err) {
		if err = os.Mkdir(configDir, 0755); err != nil {
			panic("Error making config directory")
		}
		panic("Config Not found..Please your config.YAML in the config directory")
	}
	// Check if Logs dir exists and init Logs
	if _, err := os.Stat(logsDir); os.IsNotExist(err) {
		if err = os.Mkdir(logsDir, 0755); err != nil {
			log.Fatalf("Failed to create logs directory: %s", err)
		}
	}
}

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
func StartProxy() {
	fmt.Println("\n -> Logger Initialized in logs directory")
	log.Println("Starting Reverse Proxy Server")
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
	err = http.ListenAndServeTLS(":"+port, SimpleServer.CertLocation, SimpleServer.KeyLocation, nil)
	if err != nil {
		log.Panicln("Cannot start proxy:", err)
	}
}
