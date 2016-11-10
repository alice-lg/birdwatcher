package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
  "github.com/mchackorg/birdwatcher/bird"
)

func Status(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	res := make(map[string]interface{})

	res["api"] = GetApiInfo()

  _ = bird.Status()

	js, _ := json.Marshal(res)

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
