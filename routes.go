package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/mchackorg/birdwatcher/bird"
	"net/http"
)

func ProtoRoutes(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	res := make(map[string]interface{})

	res["api"] = GetApiInfo()
	res["routes"] = bird.RoutesProto(ps.ByName("protocol"))["routes"]

	js, _ := json.Marshal(res)

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func TableRoutes(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	res := make(map[string]interface{})

	res["api"] = GetApiInfo()
	res["routes"] = bird.RoutesTable(ps.ByName("table"))["routes"]

	js, _ := json.Marshal(res)

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func ProtoCount(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	res := make(map[string]interface{})

	res["api"] = GetApiInfo()
	res["count"] = bird.RoutesProtoCount(ps.ByName("protocol"))

	js, _ := json.Marshal(res)

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func TableCount(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	res := make(map[string]interface{})

	res["api"] = GetApiInfo()
	res["count"] = bird.RoutesTable(ps.ByName("table"))

	js, _ := json.Marshal(res)

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func RouteNet(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	res := make(map[string]interface{})

	res["api"] = GetApiInfo()
	res["routes"] = bird.RoutesLookupTable(ps.ByName("net"), "master")

	js, _ := json.Marshal(res)

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func RouteNetTable(w http.ResponseWriter,
	r *http.Request,
	ps httprouter.Params) {
	res := make(map[string]interface{})

	res["api"] = GetApiInfo()
	res["routes"] = bird.RoutesLookupTable(ps.ByName("net"),
		ps.ByName("table"))

	js, _ := json.Marshal(res)

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
