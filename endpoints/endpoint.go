package endpoints

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"compress/gzip"
	"encoding/json"
	"net"
	"net/http"

	"github.com/alice-lg/birdwatcher/bird"
	"github.com/julienschmidt/httprouter"
)

type endpoint func(*http.Request, httprouter.Params, bool) (bird.Parsed, bool)

var Conf ServerConfig

func CheckAccess(req *http.Request) error {
	if len(Conf.AllowFrom) == 0 {
		return nil // AllowFrom ALL
	}

	ipStr, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		log.Println("Error parsing IP address:", err)
		return fmt.Errorf("error parsing source IP address")
	}
	clientIP := net.ParseIP(ipStr)
	if clientIP == nil {
		log.Println("Invalid IP address format:", ipStr)
		return fmt.Errorf("invalid source IP address format")
	}
	for _, allowed := range Conf.AllowFrom {
		if _, allowedNet, err := net.ParseCIDR(allowed); err == nil {
			if allowedNet.Contains(clientIP) {
				return nil
			}
		} else if allowedIP := net.ParseIP(allowed); allowedIP != nil {
			if allowedIP.Equal(clientIP) {
				return nil
			}
		} else {
			log.Printf("Invalid IP/CIDR format in configuration: %s\n", allowed);
		}
	}
	log.Println("Rejecting access from:", ipStr);
	return fmt.Errorf("%s is not allowed to access this service", ipStr);
}

func CheckUseCache(req *http.Request) bool {
	qs := req.URL.Query()

	if Conf.AllowUncached &&
		len(qs["uncached"]) == 1 && qs["uncached"][0] == "true" {
		return false
	}

	return true
}

func Endpoint(wrapped endpoint) httprouter.Handle {
	return func(w http.ResponseWriter,
		r *http.Request,
		ps httprouter.Params) {

		// Access Control
		if err := CheckAccess(r); err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		res := make(map[string]interface{})

		useCache := CheckUseCache(r)
		ret, from_cache := wrapped(r, ps, useCache)

		if reflect.DeepEqual(ret, bird.NilParse) {
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		if reflect.DeepEqual(ret, bird.BirdError) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-Type", "application/json")
			js, _ := json.Marshal(ret)
			w.Write(js)
			return
		}
		res["api"] = GetApiInfo(&ret, from_cache)

		for k, v := range ret {
			res[k] = v
		}

		w.Header().Set("Content-Type", "application/json")

		// Check if compression is supported
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			// Compress response
			w.Header().Set("Content-Encoding", "gzip")
			gz := gzip.NewWriter(w)
			defer gz.Close()
			json := json.NewEncoder(gz)
			json.Encode(res)
		} else {
			json := json.NewEncoder(w)
			json.Encode(res) // Fall back to uncompressed response
		}
	}
}

func Version(version string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(version))
	}
}
