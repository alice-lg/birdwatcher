package main

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func Protocols(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	res := make(map[string]interface{})

	res["api"] = GetApiInfo()

	lines, err := readLines(conf.Conf.FileName)
	if err != nil {
		slog.Err("Couldn't find file: " + conf.Conf.FileName)
		return
	}

	res["protocols"] = pattern("getprotocol", lines)

	js, _ := json.Marshal(res)

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
