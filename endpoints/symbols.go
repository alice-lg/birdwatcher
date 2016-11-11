package endpoints

import (
	"github.com/julienschmidt/httprouter"
	"github.com/mchackorg/birdwatcher/bird"
)

func Symbols(ps httprouter.Params) bird.Parsed {
	return bird.Symbols()
}

func SymbolTables(ps httprouter.Params) bird.Parsed {
	return bird.Parsed{"symbols": bird.Symbols()["routing table"]}
}

func SymbolProtocols(ps httprouter.Params) bird.Parsed {
	return bird.Parsed{"symbols": bird.Symbols()["protocols"]}
}
