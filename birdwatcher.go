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

func makeRouter() *httprouter.Router {
	r := httprouter.New()
	r.GET("/status", endpoints.Endpoint(endpoints.Status))
	r.GET("/protocols/bgp", endpoints.Endpoint(endpoints.Bgp))
	r.GET("/symbols", endpoints.Endpoint(endpoints.Symbols))
	r.GET("/symbols/tables", endpoints.Endpoint(endpoints.SymbolTables))
	r.GET("/symbols/protocols", endpoints.Endpoint(endpoints.SymbolProtocols))
	r.GET("/routes/protocol/:protocol", endpoints.Endpoint(endpoints.ProtoRoutes))
	r.GET("/routes/table/:table", endpoints.Endpoint(endpoints.TableRoutes))
	r.GET("/routes/count/protocol/:protocol", endpoints.Endpoint(endpoints.ProtoCount))
	r.GET("/routes/count/table/:table", endpoints.Endpoint(endpoints.TableCount))
	r.GET("/route/net/:net", endpoints.Endpoint(endpoints.RouteNet))
	r.GET("/route/net/:net/table/:table", endpoints.Endpoint(endpoints.RouteNetTable))
	r.GET("/protocols", endpoints.Endpoint(endpoints.Protocols))
	return r
}

// Print service information like, listen address,
// access restrictions and configuration flags
func PrintServiceInfo(conf *Config, birdConf BirdConfig) {
	// General Info
	log.Println("Starting Birdwatcher")
	log.Println("     Using:", birdConf.BirdCmd)
	log.Println("    Listen:", birdConf.Listen)
}

func main() {
	bird6 := flag.Bool("6", false, "Use bird6 instead of bird")
	flag.Parse()

	// Load configurations
	conf, err := LoadConfigs([]string{
		"./etc/ecix/birdwatcher.conf",
		"/etc/ecix/birdwatcher.conf",
		"./etc/ecix/birdwatcher.local.conf",
	})

	if err != nil {
		log.Fatal("Loading birdwatcher configuration failed:", err)
	}

	// Get config according to flags
	birdConf := conf.Bird
	if *bird6 {
		birdConf = conf.Bird6
	}

	PrintServiceInfo(conf, birdConf)

	// Configure client
	bird.BirdCmd = birdConf.BirdCmd

	r := makeRouter()

	realPort := strings.Join([]string{":", "23022"}, "")
	log.Fatal(http.ListenAndServe(realPort, r))
}
