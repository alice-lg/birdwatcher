package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func Routes(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprint(w, "routes\n")
}
