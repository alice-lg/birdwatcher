package endpoints

import (
	"reflect"
	"net/http"

	"github.com/alice-lg/birdwatcher/bird"
	"github.com/julienschmidt/httprouter"
)

func Symbols(r *http.Request, ps httprouter.Params) (bird.Parsed, bool) {
	return bird.Symbols()
}

func SymbolTables(r *http.Request, ps httprouter.Params) (bird.Parsed, bool) {
	val, from_cache := bird.Symbols()
	if reflect.DeepEqual(val, bird.NilParse) || reflect.DeepEqual(val, bird.BirdError) {
		return val, from_cache
	}
	return bird.Parsed{"symbols": val["symbols"].(bird.Parsed)["routing table"]}, from_cache
}

func SymbolProtocols(r *http.Request, ps httprouter.Params) (bird.Parsed, bool) {
	val, from_cache := bird.Symbols()
	if reflect.DeepEqual(val, bird.NilParse) || reflect.DeepEqual(val, bird.BirdError) {
		return val, from_cache
	}
	return bird.Parsed{"symbols": val["symbols"].(bird.Parsed)["protocol"]}, from_cache
}
