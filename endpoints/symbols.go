package endpoints

import (
	"github.com/ecix/birdwatcher/bird"
	"github.com/julienschmidt/httprouter"
)

func Symbols(ps httprouter.Params) (bird.Parsed, bool) {
	return bird.Symbols()
}

func SymbolTables(ps httprouter.Params) (bird.Parsed, bool) {
  val, from_cache := bird.Symbols()
	return bird.Parsed{"symbols": val["routing table"]}, from_cache
}

func SymbolProtocols(ps httprouter.Params) (bird.Parsed, bool) {
  val, from_cache := bird.Symbols()
	return bird.Parsed{"symbols": val["protocols"]}, from_cache
}
