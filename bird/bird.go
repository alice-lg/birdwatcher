package bird

import (
	"os/exec"
	"strings"
	"sync"
	"time"
)

var ClientConf BirdConfig
var StatusConf StatusConfig
var RateLimitConf struct {
	sync.RWMutex
	Conf RateLimitConfig
}

var Cache = struct {
	sync.RWMutex
	m map[string]Parsed
}{m: make(map[string]Parsed)}

func fromCache(key string) (Parsed, bool) {
	Cache.RLock()
	val, ok := Cache.m[key]
	Cache.RUnlock()
	if !ok {
		return nil, false
	}

	ttl, correct := val["ttl"].(time.Time)
	if !correct || ttl.Before(time.Now()) {
		return nil, false
	}

	return val, ok
}

func toCache(key string, val Parsed) {
	val["ttl"] = time.Now().Add(5 * time.Minute)
	Cache.Lock()
	Cache.m[key] = val
	Cache.Unlock()
}

func Run(args string) ([]byte, error) {
	args = "show " + args
	argsList := strings.Split(args, " ")
	return exec.Command(ClientConf.BirdCmd, argsList...).Output()
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

func RunAndParse(cmd string, parser func([]byte) Parsed) (Parsed, bool) {
	if val, ok := fromCache(cmd); ok {
		return val, true
	}

	if !checkRateLimit() {
		return nil, false
	}

	out, err := Run(cmd)

	if err != nil {
		// ignore errors for now
		return Parsed{}, false
	}

	parsed := parser(out)
	toCache(cmd, parsed)
	return parsed, false
}

func Status() (Parsed, bool) {
	birdStatus, ok := RunAndParse("status", parseStatus)
	status := birdStatus["status"].(Parsed)

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

	birdStatus["status"] = status

	return birdStatus, ok
}

func Protocols() (Parsed, bool) {
	return RunAndParse("protocols all", parseProtocols)
}

func ProtocolsBgp() (Parsed, bool) {
	p, from_cache := Protocols()
	protocols := p["protocols"].([]string)

	bgpProto := Parsed{}

	for _, v := range protocols {
		if strings.Contains(v, " BGP ") {
			key := strings.Split(v, " ")[0]
			bgpProto[key] = parseBgp(v)
		}
	}

	return Parsed{"protocols": bgpProto}, from_cache
}

func Symbols() (Parsed, bool) {
	return RunAndParse("symbols", parseSymbols)
}

func RoutesPrefixed(prefix string) (Parsed, bool) {
	return RunAndParse("route all '"+prefix+"'", parseRoutes)
}

func RoutesProto(protocol string) (Parsed, bool) {
	return RunAndParse("route protocol '"+protocol+"' all",
		parseRoutes)
}

func RoutesProtoCount(protocol string) (Parsed, bool) {
	return RunAndParse("route protocol '"+protocol+"' count",
		parseRoutesCount)
}

func RoutesFiltered(protocol string) (Parsed, bool) {
	return RunAndParse("route protocol '"+protocol+"' filtered", parseRoutes)
}

func RoutesExport(protocol string) (Parsed, bool) {
	return RunAndParse("route export '"+protocol+"' all",
		parseRoutes)
}

func RoutesExportCount(protocol string) (Parsed, bool) {
	return RunAndParse("route export '"+protocol+"' count",
		parseRoutesCount)
}

func RoutesTable(table string) (Parsed, bool) {
	return RunAndParse("route table '"+table+"' all",
		parseRoutes)
}

func RoutesTableCount(table string) (Parsed, bool) {
	return RunAndParse("route table '"+table+"' count",
		parseRoutesCount)
}

func RoutesLookupTable(net string, table string) (Parsed, bool) {
	return RunAndParse("route for '"+net+"' table '"+table+"' all",
		parseRoutes)
}

func RoutesLookupProtocol(net string, protocol string) (Parsed, bool) {
	return RunAndParse("route for '"+net+"' protocol '"+protocol+"' all",
		parseRoutes)
}
