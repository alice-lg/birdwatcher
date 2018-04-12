package main

import (
	"flag"
	"log"
	"net/http"
	"strings"

	"github.com/ecix/birdwatcher/bird"
	"github.com/ecix/birdwatcher/endpoints"

	"github.com/julienschmidt/httprouter"
)

//go:generate versionize
var VERSION = "1.11.0"

func isModuleEnabled(module string, modulesEnabled []string) bool {
	for _, enabled := range modulesEnabled {
		if enabled == module {
			return true
		}
	}

	return false
}

func makeRouter(config endpoints.ServerConfig) *httprouter.Router {
	whitelist := config.ModulesEnabled

	r := httprouter.New()
	if isModuleEnabled("status", whitelist) {
		r.GET("/version", endpoints.Version(VERSION))
		r.GET("/status", endpoints.Endpoint(endpoints.Status))
	}
	if isModuleEnabled("protocols", whitelist) {
		r.GET("/protocols", endpoints.Endpoint(endpoints.Protocols))
	}
	if isModuleEnabled("protocols_bgp", whitelist) {
		r.GET("/protocols/bgp", endpoints.Endpoint(endpoints.Bgp))
	}
	if isModuleEnabled("symbols", whitelist) {
		r.GET("/symbols", endpoints.Endpoint(endpoints.Symbols))
	}
	if isModuleEnabled("symbols_tables", whitelist) {
		r.GET("/symbols/tables", endpoints.Endpoint(endpoints.SymbolTables))
	}
	if isModuleEnabled("symbols_protocols", whitelist) {
		r.GET("/symbols/protocols", endpoints.Endpoint(endpoints.SymbolProtocols))
	}
	if isModuleEnabled("routes_protocol", whitelist) {
		r.GET("/routes/protocol/:protocol", endpoints.Endpoint(endpoints.ProtoRoutes))
	}
	if isModuleEnabled("routes_table", whitelist) {
		r.GET("/routes/table/:table", endpoints.Endpoint(endpoints.TableRoutes))
	}
	if isModuleEnabled("routes_count_protocol", whitelist) {
		r.GET("/routes/count/protocol/:protocol", endpoints.Endpoint(endpoints.ProtoCount))
	}
	if isModuleEnabled("routes_count_table", whitelist) {
		r.GET("/routes/count/table/:table", endpoints.Endpoint(endpoints.TableCount))
	}
	if isModuleEnabled("routes_filtered", whitelist) {
		r.GET("/routes/filtered/:protocol", endpoints.Endpoint(endpoints.RoutesFiltered))
	}
	if isModuleEnabled("routes_noexport", whitelist) {
		r.GET("/routes/noexport/:protocol", endpoints.Endpoint(endpoints.RoutesNoExport))
	}
	if isModuleEnabled("routes_prefixed", whitelist) {
		r.GET("/routes/prefix", endpoints.Endpoint(endpoints.RoutesPrefixed))
	}
	if isModuleEnabled("route_net", whitelist) {
		r.GET("/route/net/:net", endpoints.Endpoint(endpoints.RouteNet))
		r.GET("/route/net/:net/table/:table", endpoints.Endpoint(endpoints.RouteNetTable))
	}
	if isModuleEnabled("routes_peer", whitelist) {
		r.GET("/routes/peer", endpoints.Endpoint(endpoints.RoutesPeer))
	}
	if isModuleEnabled("routes_dump", whitelist) {
		r.GET("/routes/dump", endpoints.Endpoint(endpoints.RoutesDump))
	}
	return r
}

// Print service information like, listen address,
// access restrictions and configuration flags
func PrintServiceInfo(conf *Config, birdConf bird.BirdConfig) {
	// General Info
	log.Println("Starting Birdwatcher")
	log.Println("            Using:", birdConf.BirdCmd)
	log.Println("           Listen:", birdConf.Listen)
	log.Println("        Cache TTL:", birdConf.CacheTtl)

	// Endpoint Info
	if len(conf.Server.AllowFrom) == 0 {
		log.Println("        AllowFrom: ALL")
	} else {
		log.Println("        AllowFrom:", strings.Join(conf.Server.AllowFrom, ", "))
	}

	log.Println("   ModulesEnabled:")
	for _, m := range conf.Server.ModulesEnabled {
		log.Println("       -", m)
	}

	log.Println("   Per Peer Tables:", conf.Parser.PerPeerTables)
}

func main() {
	bird6 := flag.Bool("6", false, "Use bird6 instead of bird")
	workerPoolSize := flag.Int("worker-pool-size", 8, "Number of go routines used to parse routing tables concurrently")
	configfile := flag.String("config", "./etc/ecix/birdwatcher.conf", "Configuration file location")
	flag.Parse()

	bird.WorkerPoolSize = *workerPoolSize

	endpoints.VERSION = VERSION
	bird.InstallRateLimitReset()
	// Load configurations
	conf, err := LoadConfigs(ConfigOptions(*configfile))

	if err != nil {
		log.Fatal("Loading birdwatcher configuration failed:", err)
	}

	// Get config according to flags
	birdConf := conf.Bird
	if *bird6 {
		birdConf = conf.Bird6
		bird.IPVersion = "6"
	}

	PrintServiceInfo(conf, birdConf)

	// Configuration
	bird.ClientConf = birdConf
	bird.StatusConf = conf.Status
	bird.RateLimitConf.Conf = conf.Ratelimit
	bird.ParserConf = conf.Parser
	endpoints.Conf = conf.Server

	// Make server
	r := makeRouter(conf.Server)
	log.Fatal(http.ListenAndServe(birdConf.Listen, r))
}
