package endpoints

import (
	"github.com/ecix/birdwatcher/bird"
	"github.com/julienschmidt/httprouter"
)

func ProtoRoutes(ps httprouter.Params) bird.Parsed {
	return bird.RoutesProto(ps.ByName("protocol"))
}

func TableRoutes(ps httprouter.Params) bird.Parsed {
	return bird.RoutesTable(ps.ByName("table"))
}

func ProtoCount(ps httprouter.Params) bird.Parsed {
	return bird.RoutesProtoCount(ps.ByName("protocol"))
}

func TableCount(ps httprouter.Params) bird.Parsed {
	return bird.RoutesTable(ps.ByName("table"))
}

func RouteNet(ps httprouter.Params) bird.Parsed {
	return bird.RoutesLookupTable(ps.ByName("net"), "master")
}

func RouteNetTable(ps httprouter.Params) bird.Parsed {
	return bird.RoutesLookupTable(ps.ByName("net"), ps.ByName("table"))
}
