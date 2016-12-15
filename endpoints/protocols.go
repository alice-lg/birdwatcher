package endpoints

import (
	"net/http"

	"github.com/ecix/birdwatcher/bird"
	"github.com/julienschmidt/httprouter"
)

func Protocols(r *http.Request, ps httprouter.Params) (bird.Parsed, bool) {
	return bird.Protocols()
}

func Bgp(r *http.Request, ps httprouter.Params) (bird.Parsed, bool) {
	return bird.ProtocolsBgp()
}
