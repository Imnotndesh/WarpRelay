package SimpleProxy

import (
	"awesomeProject/ConfigParser"
	"io"
	"log"
	"net/http"
)

const (
	configLocation = "./Config/config.YAML"
)

func sampleRevProxy(w http.ResponseWriter, r *http.Request) {
	backend := "http://0.0.0.0:9080" + r.URL.Path
	newRequest, err := http.NewRequest(r.Method, backend, r.Body)
	if err != nil {
		log.Panicln("Cannot create request:", err)
		return
	}
	client := &http.Client{}
	backendResp, err := client.Do(newRequest)
	if err != nil {
		log.Panicln("Cannot make request:", err)
		return
	}
	defer backendResp.Body.Close()
	for key, values := range backendResp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	w.WriteHeader(backendResp.StatusCode)
	_, err = io.Copy(w, backendResp.Body)
	if err != nil {
		return
	}
}
func StartProxy() {
	parser := ConfigParser.ConfigParser{
		ConfigLocation: configLocation,
	}
	err := parser.ParseConfig()
	if err != nil {
		log.Println("Cannot parse config:", err)
		return
	}
	port := parser.GetProxyPort()
	http.HandleFunc("/", sampleRevProxy)
	log.Println("Starting proxy on port ", port)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Panicln("Cannot start proxy:", err)
	}
}
