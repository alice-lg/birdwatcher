package endpoints

import (
	"github.com/julienschmidt/httprouter"
	"github.com/mchackorg/birdwatcher/bird"
)

func Protocols(ps httprouter.Params) bird.Parsed {
	return bird.Protocols()
}

func Bgp(ps httprouter.Params) bird.Parsed {
	return bird.ProtocolsBgp()
}
