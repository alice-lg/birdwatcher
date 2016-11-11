package endpoints

import (
	"github.com/ecix/birdwatcher/bird"
	"github.com/julienschmidt/httprouter"
)

func Protocols(ps httprouter.Params) (bird.Parsed, bool) {
	return bird.Protocols()
}

func Bgp(ps httprouter.Params) (bird.Parsed, bool) {
	return bird.ProtocolsBgp()
}
