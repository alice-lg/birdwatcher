package endpoints

import (
	"net/http"

	"github.com/alice-lg/birdwatcher/bird"
	"github.com/julienschmidt/httprouter"
)

func Symbols(r *http.Request, ps httprouter.Params, useCache bool) (bird.Parsed, bool) {
	return bird.Symbols(useCache)
}

func SymbolTables(r *http.Request, ps httprouter.Params, useCache bool) (bird.Parsed, bool) {
	val, from_cache := bird.Symbols(useCache)
	if bird.IsSpecial(val) {
		return val, from_cache
	}
	return bird.Parsed{"symbols": val["symbols"].(bird.Parsed)["routing table"]}, from_cache
}

func SymbolProtocols(r *http.Request, ps httprouter.Params, useCache bool) (bird.Parsed, bool) {
	val, from_cache := bird.Symbols(useCache)
	if bird.IsSpecial(val) {
		return val, from_cache
	}
	return bird.Parsed{"symbols": val["symbols"].(bird.Parsed)["protocol"]}, from_cache
}
