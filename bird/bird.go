package bird

import (
	"bytes"
	"io"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"os/exec"
)

var ClientConf BirdConfig
var StatusConf StatusConfig
var IPVersion = "4"
var RateLimitConf struct {
	sync.RWMutex
	Conf RateLimitConfig
}

type Cache struct {
	sync.RWMutex
	m map[string]Parsed
}

var ParsedCache = Cache{m: make(map[string]Parsed)}
var MetaCache = Cache{m: make(map[string]Parsed)}

var NilParse Parsed = (Parsed)(nil)
var BirdError Parsed = Parsed{"error": "bird unreachable"}

var RunQueue sync.Map

func IsSpecial(ret Parsed) bool {
	return reflect.DeepEqual(ret, NilParse) || reflect.DeepEqual(ret, BirdError)
}

// Determines the key in the cache, where the result of specific functions are stored.
// Eliminates the need to know what command was executed by that function.
func GetCacheKey(fname string, fargs ...interface{}) string {
	key := strings.ToLower(fname)

	for _, arg := range fargs {
		switch arg.(type) {
		case string:
			key += "_" + strings.ToLower(arg.(string))
		}
	}

	return key
}

func (c *Cache) Store(key string, val Parsed) {
	var ttl int = 5
	if ClientConf.CacheTtl > 0 {
		ttl = ClientConf.CacheTtl
	}
	cachedAt := time.Now().UTC()
	cacheTtl := cachedAt.Add(time.Duration(ttl) * time.Minute)

	c.Lock()
	// This is not a really ... clean way of doing this.
	val["ttl"] = cacheTtl
	val["cached_at"] = cachedAt

	c.m[key] = val
	c.Unlock()
}

func (c *Cache) Get(key string) (Parsed, bool) {
	c.RLock()
	val, ok := c.m[key]
	c.RUnlock()
	if !ok {
		return NilParse, false
	}

	ttl, correct := val["ttl"].(time.Time)
	if !correct || ttl.Before(time.Now()) {
		return NilParse, false
	}

	return val, ok
}

func Run(args string) (io.Reader, error) {
	args = "-r " + "show " + args // enforce birdc in restricted mode with "-r" argument
	argsList := strings.Split(args, " ")

	out, err := exec.Command(ClientConf.BirdCmd, argsList...).Output()
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(out), nil
}

func InstallRateLimitReset() {
	go func() {
		c := time.Tick(time.Second)

		for _ = range c {
			RateLimitConf.Lock()
			RateLimitConf.Conf.Reqs = RateLimitConf.Conf.Max
			RateLimitConf.Unlock()
		}
	}()
}

func checkRateLimit() bool {
	RateLimitConf.RLock()
	check := !RateLimitConf.Conf.Enabled
	RateLimitConf.RUnlock()
	if check {
		return true
	}

	RateLimitConf.RLock()
	check = RateLimitConf.Conf.Reqs < 1
	RateLimitConf.RUnlock()
	if check {
		return false
	}

	RateLimitConf.Lock()
	RateLimitConf.Conf.Reqs -= 1
	RateLimitConf.Unlock()

	return true
}

func RunAndParse(key string, cmd string, parser func(io.Reader) Parsed, updateCache func(*Parsed)) (Parsed, bool) {
	if val, ok := ParsedCache.Get(cmd); ok {
		return val, true
	}

	var wg sync.WaitGroup
	wg.Add(1)
	if queueGroup, queueLoaded := RunQueue.LoadOrStore(cmd, &wg); queueLoaded {
		(*queueGroup.(*sync.WaitGroup)).Wait()

		if val, ok := ParsedCache.Get(cmd); ok {
			return val, true
		} else {
			// TODO BirdError should also be signaled somehow
			return NilParse, false
		}
	}

	if !checkRateLimit() {
		wg.Done()
		RunQueue.Delete(cmd)
		return NilParse, false
	}

	out, err := Run(cmd)
	if err != nil {
		// ignore errors for now
		wg.Done()
		RunQueue.Delete(cmd)
		return BirdError, false
	}

	parsed := parser(out)

	if updateCache != nil {
		updateCache(&parsed)
	}

	ParsedCache.Store(cmd, parsed)

	wg.Done()

	RunQueue.Delete(cmd)

	return parsed, false
}

func Status() (Parsed, bool) {
	updateParsedCache := func(p *Parsed) {
		status := (*p)["status"].(Parsed)

		// Last Reconfig Timestamp source:
		var lastReconfig string
		switch StatusConf.ReconfigTimestampSource {
		case "bird":
			lastReconfig = status["last_reconfig"].(string)
			break
		case "config_modified":
			lastReconfig = lastReconfigTimestampFromFileStat(
				ClientConf.ConfigFilename,
			)
		case "config_regex":
			lastReconfig = lastReconfigTimestampFromFileContent(
				ClientConf.ConfigFilename,
				StatusConf.ReconfigTimestampMatch,
			)
		}

		status["last_reconfig"] = lastReconfig

		// Filter fields
		for _, field := range StatusConf.FilterFields {
			status[field] = nil
		}
	}

	birdStatus, from_cache := RunAndParse(GetCacheKey("Status"), "status", parseStatus, updateParsedCache)
	return birdStatus, from_cache
}

func Protocols() (Parsed, bool) {
	createMetaCache := func(p *Parsed) {
		metaProtocol := Parsed{"protocols": Parsed{"bird_protocol": Parsed{}}}

		for key, _ := range (*p)["protocols"].(Parsed) {
			parsed := (*p)["protocols"].(Parsed)[key].(Parsed)

			protocol := parsed["protocol"].(string)

			birdProtocol := parsed["bird_protocol"].(string)
			// Check if the structure for the current birdProtocol already exists inside the metaProtocol cache, if not create it (BGP|Pipe|etc)
			if _, ok := metaProtocol["protocols"].(Parsed)["bird_protocol"].(Parsed)[birdProtocol]; !ok {
				metaProtocol["protocols"].(Parsed)["bird_protocol"].(Parsed)[birdProtocol] = Parsed{}
			}
			metaProtocol["protocols"].(Parsed)["bird_protocol"].(Parsed)[birdProtocol].(Parsed)[protocol] = &parsed
		}

		MetaCache.Store(GetCacheKey("Protocols"), metaProtocol)
	}

	res, from_cache := RunAndParse(GetCacheKey("Protocols"), "protocols all", parseProtocols, createMetaCache)
	return res, from_cache
}

func ProtocolsBgp() (Parsed, bool) {
	protocols, from_cache := Protocols()
	if IsSpecial(protocols) {
		return protocols, from_cache
	}

	protocolsMeta, _ := MetaCache.Get(GetCacheKey("Protocols"))
	metaProtocol := protocolsMeta["protocols"].(Parsed)

	bgpProtocols := Parsed{}

	for key, protocol := range metaProtocol["bird_protocol"].(Parsed)["BGP"].(Parsed) {
		bgpProtocols[key] = *(protocol.(*Parsed))
	}

	return Parsed{"protocols": bgpProtocols,
		"ttl":       protocols["ttl"],
		"cached_at": protocols["cached_at"]}, from_cache
}

func Symbols() (Parsed, bool) {
	return RunAndParse(GetCacheKey("Symbols"), "symbols", parseSymbols, nil)
}

func RoutesPrefixed(prefix string) (Parsed, bool) {
	cmd := routeQueryForChannel("route " + prefix + " all")
	return RunAndParse(GetCacheKey("RoutesPrefixed", prefix), cmd, parseRoutes, nil)
}

func RoutesProto(protocol string) (Parsed, bool) {
	cmd := routeQueryForChannel("route all protocol " + protocol)
	return RunAndParse(GetCacheKey("RoutesProto", protocol), cmd, parseRoutes, nil)
}

func RoutesProtoCount(protocol string) (Parsed, bool) {
	cmd := routeQueryForChannel("route protocol "+protocol) + " count"
	return RunAndParse(GetCacheKey("RoutesProtoCount", protocol), cmd, parseRoutesCount, nil)
}

func RoutesProtoPrimaryCount(protocol string) (Parsed, bool) {
	cmd := routeQueryForChannel("route primary protocol "+protocol) + " count"
	return RunAndParse(GetCacheKey("RoutesProtoPrimaryCount", protocol), cmd, parseRoutesCount, nil)
}

func PipeRoutesFilteredCount(pipe string, table string, neighborAddress string) (Parsed, bool) {
	cmd := "route table " + table + " noexport " + pipe + " where from=" + neighborAddress + " count"
	return RunAndParse(GetCacheKey("PipeRoutesFilteredCount", table, pipe, neighborAddress), cmd, parseRoutesCount, nil)
}

func PipeRoutesFiltered(pipe string, table string) (Parsed, bool) {
	cmd := routeQueryForChannel("route table '" + table + "' noexport '" + pipe + "' all")
	return RunAndParse(GetCacheKey("PipeRoutesFiltered", table, pipe), cmd, parseRoutes, nil)
}

func RoutesFiltered(protocol string) (Parsed, bool) {
	cmd := routeQueryForChannel("route all filtered protocol " + protocol)
	return RunAndParse(GetCacheKey("RoutesFiltered", protocol), cmd, parseRoutes, nil)
}

func RoutesExport(protocol string) (Parsed, bool) {
	cmd := routeQueryForChannel("route all export " + protocol)
	return RunAndParse(GetCacheKey("RoutesExport", protocol), cmd, parseRoutes, nil)
}

func RoutesNoExport(protocol string) (Parsed, bool) {
	// In case we have a multi table setup, we have to query
	// the pipe protocol.
	if ParserConf.PerPeerTables &&
		strings.HasPrefix(protocol, ParserConf.PeerProtocolPrefix) {

		// Replace prefix
		protocol = ParserConf.PipeProtocolPrefix +
			protocol[len(ParserConf.PeerProtocolPrefix):]
	}

	cmd := routeQueryForChannel("route all noexport " + protocol)
	return RunAndParse(GetCacheKey("RoutesNoExport", protocol), cmd, parseRoutes, nil)
}

func RoutesExportCount(protocol string) (Parsed, bool) {
	cmd := routeQueryForChannel("route export "+protocol) + " count"
	return RunAndParse(GetCacheKey("RoutesExportCount", protocol), cmd, parseRoutesCount, nil)
}

func RoutesTable(table string) (Parsed, bool) {
	return RunAndParse(GetCacheKey("RoutesTable", table), "route table "+table+" all", parseRoutes, nil)
}

func RoutesTableCount(table string) (Parsed, bool) {
	return RunAndParse(GetCacheKey("RoutesTableCount", table), "route table "+table+" count", parseRoutesCount, nil)
}

func RoutesLookupTable(net string, table string) (Parsed, bool) {
	return RunAndParse(GetCacheKey("RoutesLookupTable", net, table), "route for "+net+" table "+table+" all", parseRoutes, nil)
}

func RoutesLookupProtocol(net string, protocol string) (Parsed, bool) {
	return RunAndParse(GetCacheKey("RoutesLookupProtocol", net, protocol), "route for "+net+" protocol "+protocol+" all", parseRoutes, nil)
}

func RoutesPeer(peer string) (Parsed, bool) {
	cmd := routeQueryForChannel("route export " + peer)
	return RunAndParse(GetCacheKey("RoutesPeer", peer), cmd, parseRoutes, nil)
}

func RoutesDump() (Parsed, bool) {
	// TODO insert hook to update the cache with the route count information
	if ParserConf.PerPeerTables {
		return RoutesDumpPerPeerTable()
	}

	return RoutesDumpSingleTable()
}

func RoutesDumpSingleTable() (Parsed, bool) {
	importedRes, cached := RunAndParse(GetCacheKey("RoutesDumpSingleTable", "imported"), routeQueryForChannel("route all"), parseRoutes, nil)
	if IsSpecial(importedRes) {
		return importedRes, cached
	}
	filteredRes, cached := RunAndParse(GetCacheKey("RoutesDumpSingleTable", "filtered"), routeQueryForChannel("route all filtered"), parseRoutes, nil)
	if IsSpecial(filteredRes) {
		return filteredRes, cached
	}

	imported := importedRes["routes"]
	filtered := filteredRes["routes"]

	result := Parsed{
		"imported": imported,
		"filtered": filtered,
	}

	return result, cached
}

func RoutesDumpPerPeerTable() (Parsed, bool) {
	importedRes, cached := RunAndParse(GetCacheKey("RoutesDumpPerPeerTable", "imported"), routeQueryForChannel("route all"), parseRoutes, nil)
	if IsSpecial(importedRes) {
		return importedRes, cached
	}
	imported := importedRes["routes"]
	filtered := []Parsed{}

	// Get protocols with filtered routes
	protocolsRes, cached := ProtocolsBgp()
	if IsSpecial(protocolsRes) {
		return protocolsRes, cached
	}
	protocols := protocolsRes["protocols"].(Parsed)

	for protocol, details := range protocols {
		details := details.(Parsed)

		counters, ok := details["routes"].(Parsed)
		if !ok {
			continue
		}
		filterCount := counters["filtered"]
		if filterCount == 0 {
			continue // nothing to do here.
		}
		// Lookup filtered routes
		pfilteredRes, _ := RoutesFiltered(protocol)
		pfiltered, ok := pfilteredRes["routes"].([]Parsed)
		if !ok {
			continue // something went wrong...
		}

		filtered = append(filtered, pfiltered...)
	}

	result := Parsed{
		"imported": imported,
		"filtered": filtered,
	}

	return result, cached
}

func routeQueryForChannel(cmd string) string {
	status, _ := Status()
	if IsSpecial(status) {
		return cmd
	}

	birdStatus, ok := status["status"].(Parsed)
	if !ok {
		return cmd
	}

	version, ok := birdStatus["version"].(string)
	if !ok {
		return cmd
	}

	v, err := strconv.Atoi(string(version[0]))
	if err != nil || v <= 2 {
		return cmd
	}

	return cmd + " where net.type = NET_IP" + IPVersion
}
