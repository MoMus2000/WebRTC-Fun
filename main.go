package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	PORT = ":3000"
)

var (
	Offer_URL     string
	Answer_URL    string
	ConnectionMap map[string]string = make(map[string]string)
)

func ServeOfferJs(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	http.ServeFile(w, r, "./static/js/offer.js")
}

func ServeAnswerJs(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	http.ServeFile(w, r, "./static/js/answer.js")
}

func ServeOfferPage(w http.ResponseWriter, r *http.Request) {
	Offer_URL = r.URL.String()
	w.WriteHeader(200)
	http.ServeFile(w, r, "./static/offer.html")
}

func ServeAnswerPage(w http.ResponseWriter, r *http.Request) {
	Answer_URL = r.URL.String()
	fmt.Println(Answer_URL)
	w.WriteHeader(200)
	http.ServeFile(w, r, "./static/answer.html")
}

func main() {
	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/", ServeOfferPage)
	r.HandleFunc("/js/offer", ServeOfferJs)
	r.HandleFunc("/js/answer", ServeAnswerJs)
	r.HandleFunc("/answer", ServeAnswerPage)

	fmt.Println("Listening on PORT ", PORT)
	http.ListenAndServe(PORT, r)
}
