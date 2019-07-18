package bird

import (
	"bytes"
	"io"
	"log"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"os/exec"
)

type Cache interface {
	Set(key string, val Parsed, ttl int) error
	Get(key string) (Parsed, error)
	Expire() int
}

var ClientConf BirdConfig
var StatusConf StatusConfig
var IPVersion = "4"
var cache Cache // stores parsed birdc output
var CacheConf CacheConfig
var RateLimitConf struct {
	sync.RWMutex
	Conf RateLimitConfig
}
var RunQueue sync.Map // queue birdc commands before execution

var NilParse Parsed = (Parsed)(nil) // special Parsed values
var BirdError Parsed = Parsed{"error": "bird unreachable"}

func IsSpecial(ret Parsed) bool { // test for special Parsed values
	return reflect.DeepEqual(ret, NilParse) || reflect.DeepEqual(ret, BirdError)
}

// intitialize the Cache once during setup with either a MemoryCache or
// RedisCache implementation.
// TODO implement singleton pattern
func InitializeCache() {
	var err error
	if CacheConf.UseRedis {
		cache, err = NewRedisCache(CacheConf)
		if err != nil {
			log.Println("Could not initialize redis cache, falling back to memory cache:", err)
		}
	} else { // initialize the MemoryCache
		cache, err = NewMemoryCache()
		if err != nil {
			log.Fatal("Could not initialize MemoryCache:", err)
		}
	}
}

func ExpireCache() int {
	return cache.Expire()
}

/* Convenience method to make new entries in the cache.
 * Abstracts over the specific caching implementation and the ability to set
 * individual TTL values for entries. Always use the default TTL value from the
 * config.
 */
func toCache(key string, val Parsed) bool {
	var ttl int
	if ClientConf.CacheTtl > 0 {
		ttl = ClientConf.CacheTtl
	} else {
		ttl = 5 // five minutes
	}

	if err := cache.Set(key, val, ttl); err == nil {
		return true
	} else {
		log.Println(err)
		return false
	}
}

/* Convenience method to retrieve entries from the cache.
 * Abstracts over the specific caching implementations.
 * If err returned by cache.Get(key) is set, the value from the cache is not
 * used. There is either a fault e.g. missing entry or the ttl is expired.
 * Handling of specific error conditions e.g. ttl expired but entry present is
 * possible but currently not implemented.
 */
func fromCache(key string) (Parsed, bool) {
	val, err := cache.Get(key)
	if err == nil {
		return val, true
	} else {
		return val, false
	}
	//DEBUG log.Println(err)

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

func Run(args string) (io.Reader, error) {
	args = "-r " + "show " + args // enforce birdc in restricted mode with "-r" argument
	argsList := strings.Split(args, " ")

	// Allow for arguments in the config
	cmdArgs := strings.Split(ClientConf.BirdCmd, " ")
	birdc := cmdArgs[0]
	cmdArgs = cmdArgs[1:]

	cmd := []string{}
	cmd = append(cmd, cmdArgs...)
	cmd = append(cmd, argsList...)

	out, err := exec.Command(birdc, cmd...).Output()
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

func RunAndParse(useCache bool, key string, cmd string, parser func(io.Reader) Parsed, updateCache func(*Parsed)) (Parsed, bool) {
	var wg sync.WaitGroup

	if useCache {
		if val, ok := fromCache(cmd); ok {
			return val, true
		}
	}

	wg.Add(1)
	if queueGroup, queueLoaded := RunQueue.LoadOrStore(cmd, &wg); queueLoaded {
		(*queueGroup.(*sync.WaitGroup)).Wait()

		if val, ok := fromCache(cmd); ok {
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

	toCache(cmd, parsed)

	wg.Done()
	RunQueue.Delete(cmd)

	return parsed, false
}

func Status(useCache bool) (Parsed, bool) {
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

	birdStatus, from_cache := RunAndParse(useCache, GetCacheKey("Status"), "status", parseStatus, updateParsedCache)
	return birdStatus, from_cache
}

func ProtocolsShort(useCache bool) (Parsed, bool) {
	res, from_cache := RunAndParse(useCache, GetCacheKey("ProtocolsShort"), "protocols", parseProtocolsShort, nil)
	return res, from_cache
}

func Protocols(useCache bool) (Parsed, bool) {
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

		toCache(GetCacheKey("metaProtocol"), metaProtocol)
	}

	res, from_cache := RunAndParse(useCache, GetCacheKey("Protocols"), "protocols all", parseProtocols, createMetaCache)
	return res, from_cache
}

func ProtocolsBgp(useCache bool) (Parsed, bool) {
	protocols, from_cache := Protocols(useCache)
	if IsSpecial(protocols) {
		return protocols, from_cache
	}

	protocolsMeta, _ := fromCache(GetCacheKey("metaProtocol"))
	metaProtocol := protocolsMeta["protocols"].(Parsed)

	bgpProtocols := Parsed{}

	for key, protocol := range metaProtocol["bird_protocol"].(Parsed)["BGP"].(Parsed) {
		bgpProtocols[key] = *(protocol.(*Parsed))
	}

	return Parsed{"protocols": bgpProtocols,
		"ttl":       protocols["ttl"],
		"cached_at": protocols["cached_at"]}, from_cache
}

func Symbols(useCache bool) (Parsed, bool) {
	return RunAndParse(useCache, GetCacheKey("Symbols"), "symbols", parseSymbols, nil)
}

func RoutesPrefixed(useCache bool, prefix string) (Parsed, bool) {
	cmd := routeQueryForChannel(useCache, "route "+prefix+" all")
	return RunAndParse(useCache, GetCacheKey("RoutesPrefixed", prefix), cmd, parseRoutes, nil)
}

func RoutesProto(useCache bool, protocol string) (Parsed, bool) {
	cmd := routeQueryForChannel(useCache, "route all protocol "+protocol)
	return RunAndParse(useCache, GetCacheKey("RoutesProto", protocol), cmd, parseRoutes, nil)
}

func RoutesPeer(useCache bool, peer string) (Parsed, bool) {
	cmd := routeQueryForChannel(useCache, "route all where from="+peer)
	return RunAndParse(useCache, GetCacheKey("RoutesPeer", peer), cmd, parseRoutes, nil)
}

func RoutesTableAndPeer(useCache bool, table string, peer string) (Parsed, bool) {
	cmd := routeQueryForChannel(useCache, "route table "+table+" all where from="+peer)
	return RunAndParse(useCache, GetCacheKey("RoutesTableAndPeer", table, peer), cmd, parseRoutes, nil)
}

func RoutesProtoCount(useCache bool, protocol string) (Parsed, bool) {
	cmd := routeQueryForChannel(useCache, "route protocol "+protocol) + " count"
	return RunAndParse(useCache, GetCacheKey("RoutesProtoCount", protocol), cmd, parseRoutesCount, nil)
}

func RoutesProtoPrimaryCount(useCache bool, protocol string) (Parsed, bool) {
	cmd := routeQueryForChannel(useCache, "route primary protocol "+protocol) + " count"
	return RunAndParse(useCache, GetCacheKey("RoutesProtoPrimaryCount", protocol), cmd, parseRoutesCount, nil)
}

func PipeRoutesFilteredCount(useCache bool, pipe string, table string, neighborAddress string) (Parsed, bool) {
	cmd := "route table " + table + " noexport " + pipe + " where from=" + neighborAddress + " count"
	return RunAndParse(useCache, GetCacheKey("PipeRoutesFilteredCount", table, pipe, neighborAddress), cmd, parseRoutesCount, nil)
}

func PipeRoutesFiltered(useCache bool, pipe string, table string) (Parsed, bool) {
	cmd := routeQueryForChannel(useCache, "route table '"+table+"' noexport '"+pipe+"' all")
	return RunAndParse(useCache, GetCacheKey("PipeRoutesFiltered", table, pipe), cmd, parseRoutes, nil)
}

func RoutesFiltered(useCache bool, protocol string) (Parsed, bool) {
	cmd := routeQueryForChannel(useCache, "route all filtered protocol "+protocol)
	return RunAndParse(useCache, GetCacheKey("RoutesFiltered", protocol), cmd, parseRoutes, nil)
}

func RoutesExport(useCache bool, protocol string) (Parsed, bool) {
	cmd := routeQueryForChannel(useCache, "route all export "+protocol)
	return RunAndParse(useCache, GetCacheKey("RoutesExport", protocol), cmd, parseRoutes, nil)
}

func RoutesNoExport(useCache bool, protocol string) (Parsed, bool) {
	cmd := routeQueryForChannel(useCache, "route all noexport "+protocol)
	return RunAndParse(useCache, GetCacheKey("RoutesNoExport", protocol), cmd, parseRoutes, nil)
}

func RoutesExportCount(useCache bool, protocol string) (Parsed, bool) {
	cmd := routeQueryForChannel(useCache, "route export "+protocol) + " count"
	return RunAndParse(useCache, GetCacheKey("RoutesExportCount", protocol), cmd, parseRoutesCount, nil)
}

func RoutesTable(useCache bool, table string) (Parsed, bool) {
	return RunAndParse(useCache, GetCacheKey("RoutesTable", table), "route table "+table+" all", parseRoutes, nil)
}

func RoutesTableFiltered(useCache bool, table string) (Parsed, bool) {
	return RunAndParse(useCache, GetCacheKey("RoutesTableFiltered", table), "route table "+table+" filtered", parseRoutes, nil)
}

func RoutesTableCount(useCache bool, table string) (Parsed, bool) {
	return RunAndParse(useCache, GetCacheKey("RoutesTableCount", table), "route table "+table+" count", parseRoutesCount, nil)
}

func RoutesLookupTable(useCache bool, net string, table string) (Parsed, bool) {
	return RunAndParse(useCache, GetCacheKey("RoutesLookupTable", net, table), "route for "+net+" table "+table+" all", parseRoutes, nil)
}

func RoutesLookupProtocol(useCache bool, net string, protocol string) (Parsed, bool) {
	return RunAndParse(useCache, GetCacheKey("RoutesLookupProtocol", net, protocol), "route for "+net+" protocol "+protocol+" all", parseRoutes, nil)
}

func routeQueryForChannel(useCache bool, cmd string) string {
	status, _ := Status(useCache)
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
