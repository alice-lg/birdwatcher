package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func Protocols(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprint(w, "protocols\n")

	lines, err := readLines(conf.Conf.FileName)
	if err != nil {
		return
	}

	pattern(conf.Res["getprotocol"], lines)
}
