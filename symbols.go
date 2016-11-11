package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/mchackorg/birdwatcher/bird"
	"net/http"
)

func Symbols(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	res := make(map[string]interface{})

	res["api"] = GetApiInfo()

	res["symbols"] = bird.Symbols()

	js, _ := json.Marshal(res)

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func SymbolTables(w http.ResponseWriter,
	r *http.Request,
	ps httprouter.Params) {
	res := make(map[string]interface{})

	res["api"] = GetApiInfo()

	res["symbols"] = bird.Symbols()["routing table"]

	js, _ := json.Marshal(res)

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func SymbolProtocols(w http.ResponseWriter,
	r *http.Request,
	ps httprouter.Params) {
	res := make(map[string]interface{})

	res["api"] = GetApiInfo()

	res["symbols"] = bird.Symbols()["protocol"]

	js, _ := json.Marshal(res)

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
