package bird

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"

	pretty "github.com/tonnerre/golang-pretty"
)

func openFile(filename string) (*os.File, error) {
	sample := "../test/" + filename
	return os.Open(sample)
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
		routes := parseProtocolRoutes(in)
		if !reflect.DeepEqual(routes, expected[i]) {
			t.Error("Parse bgpRoutes:", routes, "expected:", expected[i])
		}
	}

}

func TestParseProtocolBgp(t *testing.T) {
	f, err := openFile("protocols_bgp_pipe.sample")
	if err != nil {
		t.Error(err)
	}
	defer f.Close()

	p := parseProtocols(f)
	log.Printf("%# v", pretty.Formatter(p))
	lines := p["protocols"].([]string)

	protocols := []Parsed{}

	for _, v := range lines {
		p2 := parseProtocol(v)
		protocols = append(protocols, p2)
	}

	if len(protocols) != 3 {
		//log.Printf("%# v", pretty.Formatter(protocols))
		t.Fatalf("Expected 3 protocols, found: %v", len(protocols))
	}

	fmt.Println(protocols)

}

func TestParseRoutesAllIpv4Bird1(t *testing.T) {
	runTestForIpv4WithFile("routes_bird1_ipv4.sample", t)
}

func TestParseRoutesAllIpv4Bird2(t *testing.T) {
	runTestForIpv4WithFile("routes_bird2_ipv4.sample", t)
}

func runTestForIpv4WithFile(file string, t *testing.T) {
	f, err := openFile(file)
	if err != nil {
		t.Error(err)
	}
	defer f.Close()

	result := parseRoutes(f)
	routes, ok := result["routes"].([]Parsed)
	if !ok {
		t.Fatal("Error getting routes")
	}

	if len(routes) != 4 {
		t.Fatal("Expected 4 routes but got ", len(routes))
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
		metric:    100,
		localPref: "100",
		protocol:  "ID8503_AS1340",
		primary:   true,
		iface:     "eno7",
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
		metric:    100,
		localPref: "100",
		protocol:  "ID8497_AS1339",
		primary:   true,
		iface:     "eno7",
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
		metric:    100,
		localPref: "100",
		protocol:  "ID8503_AS1340",
		primary:   false,
		iface:     "eno8",
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
		metric:    100,
		localPref: "100",
		protocol:  "ID8503_AS1340",
		primary:   true,
		iface:     "eno7",
	}, routes[3], "Route 4", t)
}

func TestParseRoutesAllIpv6Bird1(t *testing.T) {
	runTestForIpv6WithFile("routes_bird1_ipv6.sample", t)
}

func TestParseRoutesAllIpv6Bird2(t *testing.T) {
	runTestForIpv6WithFile("routes_bird2_ipv6.sample", t)
}

func runTestForIpv6WithFile(file string, t *testing.T) {
	f, err := openFile(file)
	if err != nil {
		t.Error(err)
	}
	defer f.Close()

	result := parseRoutes(f)
	routes, ok := result["routes"].([]Parsed)
	if !ok {
		t.Fatal("Error getting routes")
	}

	if len(routes) != 3 {
		t.Fatal("Expected 3 routes but got ", len(routes))
	}

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
		metric:    100,
		localPref: "500",
		primary:   true,
		protocol:  "upstream1",
		iface:     "eth2",
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
		localPref: "100",
		metric:    100,
		primary:   false,
		protocol:  "upstream2",
		iface:     "eth3",
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
		metric:    100,
		localPref: "5000",
		primary:   true,
		protocol:  "upstream2",
		iface:     "eth2",
	}, routes[2], "Route 3", t)
}

func assertRouteIsEqual(expected expectedRoute, actual Parsed, name string, t *testing.T) {
	if prefix := value(actual, "network", name, t).(string); prefix != expected.network {
		t.Fatal(name, ": Expected network to be:", expected.network, "not", prefix)
	}

	if nextHop := value(actual, "gateway", name, t).(string); nextHop != expected.gateway {
		t.Fatal(name, ": Expected gateway to be:", expected.gateway, "not", nextHop)
	}

	if metric := value(actual, "metric", name, t).(int64); metric != expected.metric {
		t.Fatal(name, ": Expected metric to be:", expected.metric, "not", metric)
	}

	if protocol := value(actual, "from_protocol", name, t).(string); protocol != expected.protocol {
		t.Fatal(name, ": Expected protocol to be:", expected.protocol, "not", protocol)
	}

	if iface := value(actual, "interface", name, t).(string); iface != expected.iface {
		t.Fatal(name, ": Expected interface to be:", expected.iface, "not", iface)
	}

	bgp := actual["bgp"].(Parsed)
	if localPref := value(bgp, "local_pref", name, t).(string); localPref != expected.localPref {
		t.Fatal(name, ": Expected local_pref to be:", expected.localPref, "not", localPref)
	}

	if asPath := value(bgp, "as_path", name, t).([]string); !reflect.DeepEqual(asPath, expected.asPath) {
		t.Fatal(name, ": Expected as_path to be:", expected.asPath, "not", asPath)
	}

	if community := value(bgp, "communities", name, t).([][]int64); !reflect.DeepEqual(community, expected.community) {
		t.Fatal(name, ": Expected community to be:", expected.community, "not", community)
	}

	if largeCommunity := value(bgp, "large_communities", name, t).([][]int64); !reflect.DeepEqual(largeCommunity, expected.largeCommunities) {
		t.Fatal(name, ": Expected large_community to be:", expected.largeCommunities, "not", largeCommunity)
	}
}

func value(parsed Parsed, key, name string, t *testing.T) interface{} {
	v, ok := parsed[key]
	if !ok {
		t.Fatal(name, ": Key not found", key)
	}

	return v
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
	localPref        string
	iface            string
}
