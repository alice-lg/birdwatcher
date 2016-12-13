package endpoints

import (
	"fmt"
	"log"
	"strings"

	"encoding/json"
	"net/http"

	"github.com/ecix/birdwatcher/bird"
	"github.com/julienschmidt/httprouter"
)

var Conf ServerConfig

func CheckAccess(req *http.Request) error {
	if len(Conf.AllowFrom) == 0 {
		return nil // AllowFrom ALL
	}

	// Extract IP
	tokens := strings.Split(req.RemoteAddr, ":")
	ip := strings.Join(tokens[:len(tokens)-1], ":")
	ip = strings.Replace(ip, "[", "", -1)
	ip = strings.Replace(ip, "]", "", -1)

	// Check Access
	for _, allowed := range Conf.AllowFrom {
		if ip == allowed {
			return nil
		}
	}

	// Log this request
	log.Println("Rejecting access from:", ip)

	return fmt.Errorf("%s is not allowed to access this service.", ip)
}

func Endpoint(wrapped func(httprouter.Params) (bird.Parsed, bool)) httprouter.Handle {
	return func(w http.ResponseWriter,
		r *http.Request,
		ps httprouter.Params) {

		// Access Control
		if err := CheckAccess(r); err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		res := make(map[string]interface{})

		ret, from_cache := wrapped(ps)
		if ret == nil {
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		res["api"] = GetApiInfo(from_cache)

		for k, v := range ret {
			res[k] = v
		}

		js, _ := json.Marshal(res)

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}
