package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/gorilla/mux"
)

var PROD = false

func ServeIndex(w http.ResponseWriter, r *http.Request) {
	if PROD {
		http.ServeFile(w, r, "./web/index.html")
		return
	}
	http.ServeFile(w, r, "./web/index_local.html")
	return
}

func main() {
	flag.BoolVar(&PROD, "PROD", false, "SET TO PRODUCTION MODE ...")
	flag.Parse()

	r := mux.NewRouter()

	fmt.Println(fmt.Sprintf("PROD mode set to %t .. \nIP: %s:6969", PROD, GetOutboundIP().String()))

	r.HandleFunc("/", ServeIndex)

	http.ListenAndServe(":6969", r)
}

func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
