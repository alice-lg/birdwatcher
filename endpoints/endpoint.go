package endpoints

import (
	"encoding/json"
	"net/http"

	"github.com/ecix/birdwatcher/bird"
	"github.com/julienschmidt/httprouter"
)

func Endpoint(wrapped func(httprouter.Params) (bird.Parsed, bool)) httprouter.Handle {
	return func(w http.ResponseWriter,
		r *http.Request,
		ps httprouter.Params) {
		res := make(map[string]interface{})

		ret, from_cache := wrapped(ps)
		res["api"] = GetApiInfo(from_cache)

		for k, v := range ret {
			res[k] = v
		}

		js, _ := json.Marshal(res)

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}
