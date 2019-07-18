package endpoints

import (
	"net/http"

	"github.com/alice-lg/birdwatcher/bird"
	"github.com/julienschmidt/httprouter"
)

func Protocols(r *http.Request, ps httprouter.Params, useCache bool) (bird.Parsed, bool) {
	return bird.Protocols(useCache)
}

func Bgp(r *http.Request, ps httprouter.Params, useCache bool) (bird.Parsed, bool) {
	return bird.ProtocolsBgp(useCache)
}

func ProtocolsShort(r *http.Request, ps httprouter.Params, useCache bool) (bird.Parsed, bool) {
	return bird.ProtocolsShort(useCache)
}
