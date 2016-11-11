package endpoints

import (
	"github.com/ecix/birdwatcher/bird"
	"github.com/julienschmidt/httprouter"
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
