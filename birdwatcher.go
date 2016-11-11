package main

import (
	"flag"
	"log"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/mchackorg/birdwatcher/endpoints"
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
	flag.Parse()

	r := makeRouter()

  realPort :=strings.Join([]string{":", *port}, "")
	log.Fatal(http.ListenAndServe(realPort, r))
}
