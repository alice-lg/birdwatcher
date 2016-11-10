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
