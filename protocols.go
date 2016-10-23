package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func Protocols(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	res := make(map[string]interface{})

	res["api"] = GetApiInfo()

	lines, err := readLines(conf.Conf.FileName)
	if err != nil {
		return
	}

	pattern(conf.Res["getprotocol"], lines)

	js, _ := json.Marshal(res)

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
