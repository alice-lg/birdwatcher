package bird

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"

	pretty "github.com/kr/pretty"
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
	protocols := p["protocols"].(Parsed)

	if len(protocols) != 3 {
		//log.Printf("%# v", pretty.Formatter(protocols))
		t.Fatalf("Expected 3 protocols, found: %v", len(protocols))
	}

	fmt.Println(protocols)
}

func TestParseProtocolShort(t *testing.T) {
	f, err := openFile("protocols_short.sample")
	if err != nil {
		t.Error(err)
	}
	defer f.Close()

	p := parseProtocolsShort(f)
	log.Printf("%# v", pretty.Formatter(p))
	protocols := p["protocols"].(Parsed)

	if len(protocols) != 27 {
		//log.Printf("%# v", pretty.Formatter(protocols))
		t.Fatalf("Expected 27 protocols, found: %v", len(protocols))
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
			{0, 5464},
			{0, 8339},
			{0, 8741},
			{0, 8823},
			{0, 12387},
			{0, 13101},
			{0, 16097},
			{0, 16316},
			{0, 20546},
			{0, 20686},
			{0, 20723},
			{0, 21083},
			{0, 21385},
			{0, 24940},
			{0, 25504},
			{0, 28876},
			{0, 29545},
			{0, 30058},
			{0, 31103},
			{0, 31400},
			{0, 39090},
			{0, 39392},
			{0, 39912},
			{0, 42473},
			{0, 43957},
			{0, 44453},
			{0, 47297},
			{0, 47692},
			{0, 48200},
			{0, 50629},
			{0, 51191},
			{0, 51839},
			{0, 51852},
			{0, 54113},
			{0, 56719},
			{0, 57957},
			{0, 60517},
			{0, 60574},
			{0, 61303},
			{0, 62297},
			{0, 62336},
			{0, 62359},
			{33891, 33892},
			{33891, 50673},
			{48793, 48793},
			{50673, 500},
			{65101, 11077},
			{65102, 11000},
			{65103, 724},
			{65104, 150},
		},
		largeCommunities: [][]int64{
			{9033, 65666, 12},
			{9033, 65666, 9},
		},
		extendedCommunities: []interface{}{
			[]interface{}{"rt", "42", "1234"},
			[]interface{}{"generic", "0x43000000", "0x1"},
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
		extendedCommunities: []interface{}{
			[]interface{}{"ro", "21414", "52001"},
			[]interface{}{"ro", "21414", "52004"},
			[]interface{}{"ro", "21414", "64515"},
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
		extendedCommunities: []interface{}{
			[]interface{}{"ro", "21414", "52001"},
			[]interface{}{"ro", "21414", "52004"},
			[]interface{}{"ro", "21414", "64515"},
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
		extendedCommunities: []interface{}{
			[]interface{}{"rt", "42", "1234"},
			[]interface{}{"generic", "0x43000000", "0x1"},
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
			{0, 5464},
			{0, 8339},
			{0, 8741},
			{0, 8823},
			{0, 12387},
			{0, 13101},
			{0, 16097},
			{0, 16316},
			{0, 20546},
			{0, 20686},
			{0, 20723},
			{0, 21083},
			{0, 21385},
			{0, 24940},
			{0, 25504},
			{0, 28876},
			{0, 29545},
			{0, 30058},
			{0, 31103},
			{0, 31400},
			{0, 39090},
			{0, 39392},
			{0, 39912},
			{0, 42473},
			{0, 43957},
			{0, 44453},
			{0, 47297},
			{0, 47692},
			{0, 48200},
			{0, 50629},
			{0, 51191},
			{0, 51839},
			{0, 51852},
			{0, 54113},
			{0, 56719},
			{0, 57957},
			{0, 60517},
			{0, 60574},
			{0, 61303},
			{0, 62297},
			{0, 62336},
			{0, 62359},
			{33891, 33892},
			{33891, 50673},
			{48793, 48793},
			{50673, 500},
			{65101, 11077},
			{65102, 11000},
			{65103, 724},
			{65104, 150},
		},
		largeCommunities: [][]int64{
			{48821, 0, 2000},
			{48821, 0, 2100},
		},
		extendedCommunities: []interface{}{
			[]interface{}{"ro", "21414", "52001"},
			[]interface{}{"ro", "21414", "52004"},
			[]interface{}{"ro", "21414", "64515"},
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
		extendedCommunities: []interface{}{
			[]interface{}{"ro", "21414", "52001"},
			[]interface{}{"ro", "21414", "52004"},
			[]interface{}{"ro", "21414", "64515"},
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
		extendedCommunities: []interface{}{
			[]interface{}{"unknown 0x4300", "0", "1"},
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

	if extendedCommunity, ok := bgp["ext_communities"]; ok {
		if !reflect.DeepEqual(extendedCommunity.([]interface{}), expected.extendedCommunities) {
			t.Fatal(name, ": Expected ext_community to be:", expected.extendedCommunities, "not", extendedCommunity)
		}
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
	network             string
	gateway             string
	asPath              []string
	community           [][]int64
	largeCommunities    [][]int64
	extendedCommunities []interface{}
	metric              int64
	protocol            string
	primary             bool
	localPref           string
	iface               string
}
