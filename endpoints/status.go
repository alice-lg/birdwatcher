package endpoints

import (
	"github.com/julienschmidt/httprouter"
	"github.com/mchackorg/birdwatcher/bird"
)

func Status(ps httprouter.Params) bird.Parsed {
	return bird.Status()
}
