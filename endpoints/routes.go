package endpoints

import (
	"fmt"

	"github.com/ecix/birdwatcher/bird"
	"github.com/julienschmidt/httprouter"
)

func ProtoRoutes(ps httprouter.Params) (bird.Parsed, bool) {
	protocol, err := ValidateProtocolParam(ps.ByName("protocol"))
	if err != nil {
		return bird.Parsed{"error": fmt.Sprintf("%s", err)}, false
	}
	return bird.RoutesProto(protocol)
}

func RoutesFiltered(ps httprouter.Params) (bird.Parsed, bool) {
	protocol, err := ValidateProtocolParam(ps.ByName("protocol"))
	if err != nil {
		return bird.Parsed{"error": fmt.Sprintf("%s", err)}, false
	}
	return bird.RoutesFiltered(protocol)
}

func RoutesPrefixed(ps httprouter.Params) (bird.Parsed, bool) {
	prefix, err := ValidatePrefixParam(ps.ByName("prefix"))
	if err != nil {
		return bird.Parsed{"error": fmt.Sprintf("%s", err)}, false
	}
	return bird.RoutesPrefixed(prefix)
}

func TableRoutes(ps httprouter.Params) (bird.Parsed, bool) {
	return bird.RoutesTable(ps.ByName("table"))
}

func ProtoCount(ps httprouter.Params) (bird.Parsed, bool) {
	protocol, err := ValidateProtocolParam(ps.ByName("protocol"))
	if err != nil {
		return bird.Parsed{"error": fmt.Sprintf("%s", err)}, false
	}
	return bird.RoutesProtoCount(protocol)
}

func TableCount(ps httprouter.Params) (bird.Parsed, bool) {
	return bird.RoutesTable(ps.ByName("table"))
}

func RouteNet(ps httprouter.Params) (bird.Parsed, bool) {
	return bird.RoutesLookupTable(ps.ByName("net"), "master")
}

func RouteNetTable(ps httprouter.Params) (bird.Parsed, bool) {
	return bird.RoutesLookupTable(ps.ByName("net"), ps.ByName("table"))
}
