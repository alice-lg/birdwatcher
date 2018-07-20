package endpoints

import (
	"net/http"

	"github.com/alice-lg/birdwatcher/bird"
	"github.com/julienschmidt/httprouter"
)

func Symbols(r *http.Request, ps httprouter.Params) (bird.Parsed, bool) {
	return bird.Symbols()
}

func SymbolTables(r *http.Request, ps httprouter.Params) (bird.Parsed, bool) {
	val, from_cache := bird.Symbols()
	return bird.Parsed{"symbols": val["routing table"]}, from_cache
}

func SymbolProtocols(r *http.Request, ps httprouter.Params) (bird.Parsed, bool) {
	val, from_cache := bird.Symbols()
	return bird.Parsed{"symbols": val["protocols"]}, from_cache
}
