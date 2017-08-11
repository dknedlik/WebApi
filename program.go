package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func main() {
	// redirect every http request to https
	go http.ListenAndServe(":3000", http.HandlerFunc(redirect))

	router := mux.NewRouter()

	router.Handle("/", http.FileServer(http.Dir("./views")))
	// We will setup our server so we can serve static assest like images, css from the /static/{file} route
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	router.Handle("/records", http.HandlerFunc(notImplemented)).Methods("Get")
	router.Handle("/records/{id}", http.HandlerFunc(notImplemented)).Methods("Get")
	router.Handle("/records", http.HandlerFunc(notImplemented)).Methods("Post")
	err := http.ListenAndServeTLS(":4443", "server.crt", "server.key", router)

	if err != nil {
		log.Println(err.Error())
	}
}

func notImplemented(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Not Implemented"))
}

func redirect(w http.ResponseWriter, req *http.Request) {
	h := strings.Split(req.Host, ":")
	host := h[0]
	target := "https://" + host + ":4443" + req.URL.Path
	if len(req.URL.RawQuery) > 0 {
		target += "?" + req.URL.RawQuery
	}

	log.Printf("redirect to: %s", target)
	http.Redirect(w, req, target, http.StatusTemporaryRedirect)

}
