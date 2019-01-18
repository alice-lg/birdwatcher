package endpoints

import (
	"fmt"
	"net/http"

	"github.com/alice-lg/birdwatcher/bird"
	"github.com/julienschmidt/httprouter"
)

func ProtoRoutes(r *http.Request, ps httprouter.Params) (bird.Parsed, bool) {
	protocol, err := ValidateProtocolParam(ps.ByName("protocol"))
	if err != nil {
		return bird.Parsed{"error": fmt.Sprintf("%s", err)}, false
	}
	return bird.RoutesProto(protocol)
}

func RoutesFiltered(r *http.Request, ps httprouter.Params) (bird.Parsed, bool) {
	protocol, err := ValidateProtocolParam(ps.ByName("protocol"))
	if err != nil {
		return bird.Parsed{"error": fmt.Sprintf("%s", err)}, false
	}
	return bird.RoutesFiltered(protocol)
}

func RoutesNoExport(r *http.Request, ps httprouter.Params) (bird.Parsed, bool) {
	protocol, err := ValidateProtocolParam(ps.ByName("protocol"))
	if err != nil {
		return bird.Parsed{"error": fmt.Sprintf("%s", err)}, false
	}
	return bird.RoutesNoExport(protocol)
}

func RoutesPrefixed(r *http.Request, ps httprouter.Params) (bird.Parsed, bool) {
	qs := r.URL.Query()
	prefixl := qs["prefix"]
	if len(prefixl) != 1 {
		return bird.Parsed{"error": "need a prefix as single query parameter"}, false
	}

	prefix, err := ValidatePrefixParam(prefixl[0])
	if err != nil {
		return bird.Parsed{"error": fmt.Sprintf("%s", err)}, false
	}
	return bird.RoutesPrefixed(prefix)
}

func TableRoutes(r *http.Request, ps httprouter.Params) (bird.Parsed, bool) {
	return bird.RoutesTable(ps.ByName("table"))
}

func TableRoutesFiltered(r *http.Request, ps httprouter.Params) (bird.Parsed, bool) {
	return bird.RoutesTableFiltered(ps.ByName("table"))
}

func TableAndPeerRoutes(r *http.Request, ps httprouter.Params) (bird.Parsed, bool) {
	return bird.RoutesTableAndPeer(ps.ByName("table"), ps.ByName("peer"))
}

func ProtoCount(r *http.Request, ps httprouter.Params) (bird.Parsed, bool) {
	protocol, err := ValidateProtocolParam(ps.ByName("protocol"))
	if err != nil {
		return bird.Parsed{"error": fmt.Sprintf("%s", err)}, false
	}
	return bird.RoutesProtoCount(protocol)
}

func ProtoPrimaryCount(r *http.Request, ps httprouter.Params) (bird.Parsed, bool) {
	protocol, err := ValidateProtocolParam(ps.ByName("protocol"))
	if err != nil {
		return bird.Parsed{"error": fmt.Sprintf("%s", err)}, false
	}
	return bird.RoutesProtoPrimaryCount(protocol)
}

func TableCount(r *http.Request, ps httprouter.Params) (bird.Parsed, bool) {
	return bird.RoutesTableCount(ps.ByName("table"))
}

func RouteNet(r *http.Request, ps httprouter.Params) (bird.Parsed, bool) {
	return bird.RoutesLookupTable(ps.ByName("net"), "master")
}

func RouteNetTable(r *http.Request, ps httprouter.Params) (bird.Parsed, bool) {
	return bird.RoutesLookupTable(ps.ByName("net"), ps.ByName("table"))
}

func PipeRoutesFiltered(r *http.Request, ps httprouter.Params) (bird.Parsed, bool) {
	qs := r.URL.Query()
	table := qs["table"][0]
	pipe := qs["pipe"][0]
	return bird.PipeRoutesFiltered(pipe, table)
}

func PipeRoutesFilteredCount(r *http.Request, ps httprouter.Params) (bird.Parsed, bool) {
	qs := r.URL.Query()
	table := qs["table"][0]
	pipe := qs["pipe"][0]
	address := qs["address"][0]
	return bird.PipeRoutesFilteredCount(pipe, table, address)
}

func PeerRoutes(r *http.Request, ps httprouter.Params) (bird.Parsed, bool) {
	return bird.RoutesPeer(ps.ByName("peer"))
}
