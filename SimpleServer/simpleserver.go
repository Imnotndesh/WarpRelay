package SimpleServer

import (
	"log"
	"net/http"
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
func StartServer(port, message string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", helloHandler(message))
	for killServer != true {
		log.Println("Starting simpleServer...")
		if err := http.ListenAndServe(port, mux); err != nil {
			log.Panicln("Cannot start server rn")
		}
	}
}
func StopServer() {
	killServer = false
}
