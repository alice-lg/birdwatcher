package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

func Test(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprint(w, "Test\n")
}

func main() {
	r := httprouter.New()
	r.GET("/test", Test)

	log.Fatal(http.ListenAndServe(":8080", r))
}
