package endpoints

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/mchackorg/birdwatcher/bird"
)

func Endpoint(wrapped func(httprouter.Params) (bird.Parsed)) httprouter.Handle {
  return func(w http.ResponseWriter,
              r *http.Request,
              ps httprouter.Params) {
	  res := make(map[string]interface{})

	  res["api"] = GetApiInfo()

    ret := wrapped(ps)

    for k, v := range ret {
	    res[k] = v
    }

	  js, _ := json.Marshal(res)

	  w.Header().Set("Content-Type", "application/json")
	  w.Write(js)
  }
}
