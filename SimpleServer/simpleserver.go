package SimpleServer

import (
	"fmt"
	"log"
	"net/http"
)

// Cert Locations for HTTPS
const (
	CertLocation = "./Certs/server.crt"
	KeyLocation  = "./Certs/server.key"
)

var (
	killServer bool
)

func helloHandler(message string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			sayHelloWithName(w, r)
		case http.MethodGet:
			sayHello(w, r, message)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
	}
}
func secureHelloHandler(message string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			sayHelloSecurely(w, r, message)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
	}
}
func sayHello(w http.ResponseWriter, r *http.Request, message string) {
	if message == "" {
		message = "Hello World"
	}
	_, err := w.Write([]byte(message))
	if err != nil {
		log.Println(err)
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
func sayHelloSecurely(w http.ResponseWriter, r *http.Request, message string) {
	if message == "" {
		message = "Secure Hello World"
	}
	_, err := fmt.Fprintf(w, message)
	if err != nil {
		log.Panicln(err)
		return
	}
}
func StartServer(port, message string, secure ...bool) {
	mux := http.NewServeMux()
	runInSecureMode := secure[0]
	if runInSecureMode != false {
		// HTTPS simple server
		mux.HandleFunc("/secure-hello", secureHelloHandler(message))
		for killServer != true {
			log.Println("Starting secure server")
			if err := http.ListenAndServeTLS(port, CertLocation, KeyLocation, mux); err != nil {
				log.Panicln("Cannot start secure server: ", err)
			}
		}
	} else {
		// HTTP simple server
		mux.HandleFunc("/hello", helloHandler(message))
		for killServer != true {
			log.Println("Starting simpleServer...")
			if err := http.ListenAndServe(port, mux); err != nil {
				log.Panicln("Cannot start server rn")
			}
		}
	}
}
func StopServer() {
	killServer = false
}
