package endpoints

import (
	"net/http"

	"github.com/ecix/birdwatcher/bird"
	"github.com/julienschmidt/httprouter"
)

func Status(r *http.Request, ps httprouter.Params) (bird.Parsed, bool) {
	return bird.Status()
}
