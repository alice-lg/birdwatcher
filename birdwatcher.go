package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

func Status(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprint(w, "status\n")
}

func Routes(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprint(w, "routes\n")
}

func Protocols(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprint(w, "protocols\n")
}

func main() {
	r := httprouter.New()
	r.GET("/status", Status)
	r.GET("/routes", Routes)
	r.GET("/protocols", Protocols)

	log.Fatal(http.ListenAndServe(":8080", r))
}
