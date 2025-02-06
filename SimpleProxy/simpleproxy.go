package SimpleProxy

import (
	"io"
	"log"
	"net/http"
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
	http.HandleFunc("/", sampleRevProxy)
	err := http.ListenAndServe(":7080", nil)
	if err != nil {
		return
	}
}
