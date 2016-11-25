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

func main() {
	port := flag.String("port",
		"29184",
		"The port the birdwatcher should run on")
	birdc := flag.String("birdc",
		"birdc",
		"The birdc command to use (for IPv6, use birdc6)")
	flag.Parse()

	bird.BirdCmd = *birdc
	bird.InstallRateLimitReset()

	r := makeRouter()

	realPort := strings.Join([]string{":", *port}, "")
	log.Fatal(http.ListenAndServe(realPort, r))
}
