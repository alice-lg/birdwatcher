package bird

import (
	"io/ioutil"
	"reflect"
	"testing"
)

func readSampleData(filename string) ([]byte, error) {
	sample := "../test/" + filename
	return ioutil.ReadFile(sample)
}

func TestParseBgpRoutes(t *testing.T) {

	inputs := []string{
		"1 imported, 0 filtered, 2 exported, 1 preferred",
		"0 imported, 2846 exported", // Bird 1.4.x
	}

	expected := []Parsed{
		Parsed{
			"imported":  int64(1),
			"filtered":  int64(0),
			"exported":  int64(2),
			"preferred": int64(1),
		},
		Parsed{
			"imported": int64(0),
			"exported": int64(2846),
		},
	}

	for i, in := range inputs {
		routes := parseBgpRoutes(in)
		if !reflect.DeepEqual(routes, expected[i]) {
			t.Error("Parse bgpRoutes:", routes, "expected:", expected[i])
		}
	}

}

func TestParseRoutesAll(t *testing.T) {
	sample, err := readSampleData("routes_bird1_ipv4.sample")
	if err != nil {
		t.Error(err)
	}

	result := parseRoutes(sample)
	routes, ok := result["routes"].([]Parsed)
	if !ok {
		t.Error("Error getting routes")
	}

	if len(routes) != 4 {
		t.Error("Expected number of routes to be 4")
	}

	assertRouteIsEqual(expectedRoute{
		network: "16.0.0.0/24",
		gateway: "1.2.3.16",
		asPath:  []string{"1340"},
		community: [][]int64{
			{65011, 3},
			{9033, 3251},
		},
		largeCommunities: [][]int64{
			{9033, 65666, 12},
			{9033, 65666, 9},
		},
		metric:   100,
		protocol: "ID8503_AS1340",
		primary:  true,
	}, routes[0], "Route 1", t)
	assertRouteIsEqual(expectedRoute{
		network: "200.0.0.0/24",
		gateway: "1.2.3.15",
		asPath:  []string{"1339"},
		community: [][]int64{
			{65011, 40},
			{9033, 3251},
		},
		largeCommunities: [][]int64{
			{9033, 65666, 12},
			{9033, 65666, 9},
		},
		metric:   100,
		protocol: "ID8503_AS1340",
		primary:  true,
	}, routes[1], "Route 2", t)
	assertRouteIsEqual(expectedRoute{
		network: "200.0.0.0/24",
		gateway: "1.2.3.16",
		asPath:  []string{"1340"},
		community: [][]int64{
			{65011, 3},
			{9033, 3251},
		},
		largeCommunities: [][]int64{
			{9033, 65666, 12},
			{9033, 65666, 9},
		},
		metric:   100,
		protocol: "ID8503_AS1340",
		primary:  false,
	}, routes[2], "Route 3", t)
	assertRouteIsEqual(expectedRoute{
		network: "16.0.0.0/24",
		gateway: "1.2.3.16",
		asPath:  []string{"1340"},
		community: [][]int64{
			{65011, 3},
			{9033, 3251},
		},
		largeCommunities: [][]int64{
			{9033, 65666, 12},
			{9033, 65666, 9},
		},
		metric:   100,
		protocol: "ID8503_AS1340",
		primary:  true,
	}, routes[3], "Route 4", t)
}

func TestParseRoutesAllBird1(t *testing.T) {
	sample, err := readSampleData("routes_bird1_ipv6.sample")
	if err != nil {
		t.Error(err)
	}

	result := parseRoutes(sample)
	routes, ok := result["routes"].([]Parsed)
	if !ok {
		t.Fatal("Error getting routes")
	}

	if len(routes) != 3 {
		t.Fatalf("Expected 3 routes but got %d", len(routes))
	}

	assertIpv6RoutesAsExpected(routes, t)
}

func TestParseRoutesAllBird2(t *testing.T) {
	sample, err := readSampleData("routes_bird2_ipv6.sample")
	if err != nil {
		t.Error(err)
	}

	result := parseRoutes(sample)
	routes, ok := result["routes"].([]Parsed)
	if !ok {
		t.Fatal("Error getting routes")
	}

	if len(routes) != 3 {
		t.Fatalf("Expected 3 routes but got %d", len(routes))
	}

	assertIpv6RoutesAsExpected(routes, t)
}

func assertIpv6RoutesAsExpected(routes []Parsed, t *testing.T) {
	assertRouteIsEqual(expectedRoute{
		network: "2001:4860::/32",
		gateway: "fe80:ffff:ffff::1",
		asPath:  []string{"15169"},
		community: [][]int64{
			{9033, 3001},
			{65000, 680},
		},
		largeCommunities: [][]int64{
			{48821, 0, 2000},
			{48821, 0, 2100},
		},
		metric:   500,
		primary:  true,
		protocol: "upstream1",
	}, routes[0], "Route 1", t)
	assertRouteIsEqual(expectedRoute{
		network: "2001:4860::/32",
		gateway: "fe80:ffff:ffff::2",
		asPath:  []string{"50629", "15169"},
		community: [][]int64{
			{50629, 200},
			{50629, 201},
		},
		largeCommunities: [][]int64{
			{48821, 0, 3000},
			{48821, 0, 3100},
		},
		metric:   100,
		primary:  false,
		protocol: "upstream2",
	}, routes[1], "Route 2", t)
	assertRouteIsEqual(expectedRoute{
		network: "2001:678:1e0::/48",
		gateway: "fe80:ffff:ffff::2",
		asPath:  []string{"202739"},
		community: [][]int64{
			{48821, 2000},
			{48821, 2100},
		},
		largeCommunities: [][]int64{
			{48821, 0, 2000},
			{48821, 0, 2100},
		},
		metric:   5000,
		primary:  true,
		protocol: "upstream2",
	}, routes[2], "Route 3", t)
}

func assertRouteIsEqual(expected expectedRoute, actual Parsed, name string, t *testing.T) {
	if prefix := actual["network"].(string); prefix != expected.network {
		t.Fatal(name, ": Expected network to be:", expected.network, "not", prefix)
	}

	if nextHop := actual["gateway"].(string); nextHop != expected.gateway {
		t.Fatal(name, ": Expected gateway to be:", expected.gateway, "not", nextHop)
	}

	if metric := actual["metric"].(int64); metric != expected.metric {
		t.Fatal(name, ": Expected metric to be:", expected.metric, "not", metric)
	}

	if protocol := actual["from_protocol"].(string); protocol != expected.protocol {
		t.Fatal(name, ": Expected protocol to be:", expected.protocol, "not", protocol)
	}

	bgp := actual["bgp"].(Parsed)
	if asPath := bgp["as_path"].([]string); !reflect.DeepEqual(asPath, expected.asPath) {
		t.Fatal(name, ": Expected as_path to be:", expected.asPath, "not", asPath)
	}

	if community := bgp["communities"].([][]int64); !reflect.DeepEqual(community, expected.community) {
		t.Fatal(name, ": Expected community to be:", expected.community, "not", community)
	}

	if largeCommunity := bgp["large_communities"].([][]int64); !reflect.DeepEqual(largeCommunity, expected.largeCommunities) {
		t.Fatal(name, ": Expected large_community to be:", expected.largeCommunities, "not", largeCommunity)
	}
}

type expectedRoute struct {
	network          string
	gateway          string
	asPath           []string
	community        [][]int64
	largeCommunities [][]int64
	metric           int64
	protocol         string
	primary          bool
}
