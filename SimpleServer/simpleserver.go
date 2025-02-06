package SimpleServer

import (
	"log"
	"net/http"
)

var (
	killServer bool
)

func helloHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			sayHelloWithName(w, r)
		case http.MethodGet:
			sayHello(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
	}
}
func sayHello(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Hello User"))
	if err != nil {
		log.Panicln(err)
	}
}
func sayHelloWithName(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if _, err := w.Write([]byte("Hello " + name)); err != nil {
		log.Panicln(err)
		return
	}

}
func StartServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", helloHandler())
	for killServer != true {
		if err := http.ListenAndServe(":9080", mux); err != nil {
			log.Panicln("Cannot start server rn")
		}
	}
}
func StopServer() {
	killServer = false
}
