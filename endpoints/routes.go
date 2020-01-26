package endpoints

import (
	"fmt"
	"net/http"

	"github.com/alice-lg/birdwatcher/bird"
	"github.com/julienschmidt/httprouter"
)

func ProtoRoutes(r *http.Request, ps httprouter.Params, useCache bool) (bird.Parsed, bool) {
	protocol, err := ValidateProtocolParam(ps.ByName("protocol"))
	if err != nil {
		return bird.Parsed{"error": fmt.Sprintf("%s", err)}, false
	}

	return bird.RoutesProto(useCache, protocol)
}

func RoutesFiltered(r *http.Request, ps httprouter.Params, useCache bool) (bird.Parsed, bool) {
	protocol, err := ValidateProtocolParam(ps.ByName("protocol"))
	if err != nil {
		return bird.Parsed{"error": fmt.Sprintf("%s", err)}, false
	}

	return bird.RoutesFiltered(useCache, protocol)
}

func RoutesExport(r *http.Request, ps httprouter.Params, useCache bool) (bird.Parsed, bool) {
	protocol, err := ValidateProtocolParam(ps.ByName("protocol"))
	if err != nil {
		return bird.Parsed{"error": fmt.Sprintf("%s", err)}, false
	}

	return bird.RoutesExport(useCache, protocol)
}

func RoutesNoExport(r *http.Request, ps httprouter.Params, useCache bool) (bird.Parsed, bool) {
	protocol, err := ValidateProtocolParam(ps.ByName("protocol"))
	if err != nil {
		return bird.Parsed{"error": fmt.Sprintf("%s", err)}, false
	}

	return bird.RoutesNoExport(useCache, protocol)
}

func RoutesPrefixed(r *http.Request, ps httprouter.Params, useCache bool) (bird.Parsed, bool) {
	qs := r.URL.Query()
	prefixl := qs["prefix"]
	if len(prefixl) != 1 {
		return bird.Parsed{"error": "need a prefix as single query parameter"}, false
	}

	prefix, err := ValidatePrefixParam(prefixl[0])
	if err != nil {
		return bird.Parsed{"error": fmt.Sprintf("%s", err)}, false
	}

	return bird.RoutesPrefixed(useCache, prefix)
}

func TableRoutes(r *http.Request, ps httprouter.Params, useCache bool) (bird.Parsed, bool) {
	table, err := ValidateProtocolParam(ps.ByName("table"))
	if err != nil {
		return bird.Parsed{"error": fmt.Sprintf("%s", err)}, false
	}

	return bird.RoutesTable(useCache, table)
}

func TableRoutesFiltered(r *http.Request, ps httprouter.Params, useCache bool) (bird.Parsed, bool) {
	table, err := ValidateProtocolParam(ps.ByName("table"))
	if err != nil {
		return bird.Parsed{"error": fmt.Sprintf("%s", err)}, false
	}

	return bird.RoutesTableFiltered(useCache, table)
}

func TableAndPeerRoutes(r *http.Request, ps httprouter.Params, useCache bool) (bird.Parsed, bool) {
	table, err := ValidateProtocolParam(ps.ByName("table"))
	if err != nil {
		return bird.Parsed{"error": fmt.Sprintf("%s", err)}, false
	}

	peer, err := ValidatePrefixParam(ps.ByName("peer"))
	if err != nil {
		return bird.Parsed{"error": fmt.Sprintf("%s", err)}, false
	}

	return bird.RoutesTableAndPeer(useCache, table, peer)
}

func ProtoCount(r *http.Request, ps httprouter.Params, useCache bool) (bird.Parsed, bool) {
	protocol, err := ValidateProtocolParam(ps.ByName("protocol"))
	if err != nil {
		return bird.Parsed{"error": fmt.Sprintf("%s", err)}, false
	}

	return bird.RoutesProtoCount(useCache, protocol)
}

func ProtoPrimaryCount(r *http.Request, ps httprouter.Params, useCache bool) (bird.Parsed, bool) {
	protocol, err := ValidateProtocolParam(ps.ByName("protocol"))
	if err != nil {
		return bird.Parsed{"error": fmt.Sprintf("%s", err)}, false
	}
	return bird.RoutesProtoPrimaryCount(useCache, protocol)
}

func TableCount(r *http.Request, ps httprouter.Params, useCache bool) (bird.Parsed, bool) {
	table, err := ValidateProtocolParam(ps.ByName("table"))
	if err != nil {
		return bird.Parsed{"error": fmt.Sprintf("%s", err)}, false
	}

	return bird.RoutesTableCount(useCache, table)
}

func RouteNet(r *http.Request, ps httprouter.Params, useCache bool) (bird.Parsed, bool) {
	net, err := ValidatePrefixParam(ps.ByName("net"))
	if err != nil {
		return bird.Parsed{"error": fmt.Sprintf("%s", err)}, false
	}

	return bird.RoutesLookupTable(useCache, net, "master")
}

func RouteNetTable(r *http.Request, ps httprouter.Params, useCache bool) (bird.Parsed, bool) {
	net, err := ValidatePrefixParam(ps.ByName("net"))
	if err != nil {
		return bird.Parsed{"error": fmt.Sprintf("%s", err)}, false
	}

	table, err := ValidateProtocolParam(ps.ByName("table"))
	if err != nil {
		return bird.Parsed{"error": fmt.Sprintf("%s", err)}, false
	}

	return bird.RoutesLookupTable(useCache, net, table)
}

func PipeRoutesFiltered(r *http.Request, ps httprouter.Params, useCache bool) (bird.Parsed, bool) {
	qs := r.URL.Query()

	if len(qs["table"]) != 1 {
		return bird.Parsed{"error": "need a table as single query parameter"}, false
	}
	table, err := ValidateProtocolParam(qs["table"][0])
	if err != nil {
		return bird.Parsed{"error": fmt.Sprintf("%s", err)}, false
	}

	if len(qs["pipe"]) != 1 {
		return bird.Parsed{"error": "need a pipe as single query parameter"}, false
	}
	pipe, err := ValidateProtocolParam(qs["pipe"][0])
	if err != nil {
		return bird.Parsed{"error": fmt.Sprintf("%s", err)}, false
	}

	return bird.PipeRoutesFiltered(useCache, pipe, table)
}

func PipeRoutesFilteredCount(r *http.Request, ps httprouter.Params, useCache bool) (bird.Parsed, bool) {
	qs := r.URL.Query()

	if len(qs["table"]) != 1 {
		return bird.Parsed{"error": "need a table as single query parameter"}, false
	}
	table, err := ValidateProtocolParam(qs["table"][0])
	if err != nil {
		return bird.Parsed{"error": fmt.Sprintf("%s", err)}, false
	}

	if len(qs["pipe"]) != 1 {
		return bird.Parsed{"error": "need a pipe as single query parameter"}, false
	}
	pipe, err := ValidateProtocolParam(qs["pipe"][0])
	if err != nil {
		return bird.Parsed{"error": fmt.Sprintf("%s", err)}, false
	}

	if len(qs["address"]) != 1 {
		return bird.Parsed{"error": "need a address as single query parameter"}, false
	}
	address, err := ValidatePrefixParam(qs["address"][0])
	if err != nil {
		return bird.Parsed{"error": fmt.Sprintf("%s", err)}, false
	}

	return bird.PipeRoutesFilteredCount(useCache, pipe, table, address)
}

func PeerRoutes(r *http.Request, ps httprouter.Params, useCache bool) (bird.Parsed, bool) {
	peer, err := ValidatePrefixParam(ps.ByName("peer"))
	if err != nil {
		return bird.Parsed{"error": fmt.Sprintf("%s", err)}, false
	}

	return bird.RoutesPeer(useCache, peer)
}
