package main

import (
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

func main() {
	r := httprouter.New()
	r.GET("/status", Status)
	r.GET("/routes", Routes)
	r.GET("/protocols", Protocols)

	log.Fatal(http.ListenAndServe(":8080", r))
}
