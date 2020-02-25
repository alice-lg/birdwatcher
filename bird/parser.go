package bird

import (
	"bufio"
	"io"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

// WorkerPoolSize is the number of go routines used to parse routing tables concurrently
var WorkerPoolSize = 8

var (
	ParserConf ParserConfig
	regex      struct {
		status struct {
			startLine     *regexp.Regexp
			routerID      *regexp.Regexp
			currentServer *regexp.Regexp
			lastReboot    *regexp.Regexp
			lastReconfig  *regexp.Regexp
		}
		protocol struct {
			channel      *regexp.Regexp
			protocol     *regexp.Regexp
			numericValue *regexp.Regexp
			routes       *regexp.Regexp
			stringValue  *regexp.Regexp
			routeChanges *regexp.Regexp
			short        *regexp.Regexp
		}
		symbols struct {
			keyRx *regexp.Regexp
		}
		routeCount struct {
			countRx *regexp.Regexp
		}
		routes struct {
			startDefinition   *regexp.Regexp
			second            *regexp.Regexp
			routeType         *regexp.Regexp
			bgp               *regexp.Regexp
			community         *regexp.Regexp
			largeCommunity    *regexp.Regexp
			extendedCommunity *regexp.Regexp
			origin            *regexp.Regexp
			prefixBird2       *regexp.Regexp
			gatewayBird2      *regexp.Regexp
		}
	}
)

type Parsed map[string]interface{}

func init() {
	regex.status.startLine = regexp.MustCompile(`^BIRD\sv?([0-9\.]+)\s*$`)
	regex.status.routerID = regexp.MustCompile(`^Router\sID\sis\s([0-9\.]+)\s*$`)
	regex.status.currentServer = regexp.MustCompile(`^Current\sserver\stime\sis\s([0-9\-]+\s[0-9\:]+)\s*$`)
	regex.status.lastReboot = regexp.MustCompile(`^Last\sreboot\son\s([0-9\-]+\s[0-9\:]+)\s*$`)
	regex.status.lastReconfig = regexp.MustCompile(`^Last\sreconfiguration\son\s([0-9\-]+\s[0-9\:]+)\s*$`)

	regex.symbols.keyRx = regexp.MustCompile(`^([^\s]+)\s+(.+)\s*$`)

	regex.routeCount.countRx = regexp.MustCompile(`^(\d+)\s+of\s+(\d+)\s+routes.*$`)

	regex.protocol.channel = regexp.MustCompile("Channel ipv([46])")
	// regex.protocol.protocol = regexp.MustCompile(`^(?:1002\-)?([^\s]+)\s+(BGP|RPKI|Pipe|BFD|Direct|Device|Kernel)\s+([^\s]+)\s+([^\s]+)\s+(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}|[^\s]+)(?:\s+(.*?)\s*)?$`)
	regex.protocol.protocol = regexp.MustCompile(`^(?:1002\-)?([^\s]+)\s+(\w+)\s+([^\s]+)\s+([^\s]+)\s+(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}|[^\s]+)(?:\s+(.*?)\s*)?$`)
	regex.protocol.numericValue = regexp.MustCompile(`^\s+([^:]+):\s+([\d]+)\s*$`)
	regex.protocol.routes = regexp.MustCompile(`^\s+Routes:\s+(.*)`)
	regex.protocol.stringValue = regexp.MustCompile(`^\s+([^:]+):\s+(.+)\s*$`)
	regex.protocol.routeChanges = regexp.MustCompile(`(Import|Export) (updates|withdraws):\s+(\d+|---)\s+(\d+|---)\s+(\d+|---)\s+(\d+|---)\s+(\d+|---)\s*$`)

	regex.routes.startDefinition = regexp.MustCompile(`^([0-9a-f\.\:\/]+)\s+via\s+([0-9a-f\.\:]+)\s+on\s+([\w\.]+)\s+\[([\w\.:]+)\s+([0-9\-\:\s]+)(?:\s+from\s+([0-9a-f\.\:\/]+)){0,1}\]\s+(?:(\*)\s+){0,1}\((\d+)(?:\/\d+){0,1}\).*`)
	regex.protocol.short = regexp.MustCompile(`^(?:1002\-)?([^\s]+)\s+([^\s]+)\s+([^\s]+)\s+([^\s]+)\s+([0-9\-]+\s+[0-9\:]+?|[0-9\-]+)\s+(.*?)\s*?$`)
	regex.routes.second = regexp.MustCompile(`^\s+via\s+([0-9a-f\.\:]+)\s+on\s+([\w\.]+)\s+\[([\w\.:]+)\s+([0-9\-\:\s]+)(?:\s+from\s+([0-9a-f\.\:\/]+)){0,1}\]\s+(?:(\*)\s+){0,1}\((\d+)(?:\/\d+){0,1}\).*$`)
	regex.routes.routeType = regexp.MustCompile(`^\s+Type:\s+(.*)\s*$`)
	regex.routes.bgp = regexp.MustCompile(`^\s+BGP.(\w+):\s+(.+)\s*$`)
	regex.routes.community = regexp.MustCompile(`^\((\d+),\s*(\d+)\)`)
	regex.routes.largeCommunity = regexp.MustCompile(`^\((\d+),\s*(\d+),\s*(\d+)\)`)
	regex.routes.extendedCommunity = regexp.MustCompile(`^\(([^,]+),\s*([^,]+),\s*([^,]+)\)`)
	regex.routes.origin = regexp.MustCompile(`\([^\(]*\)\s*`)
	regex.routes.prefixBird2 = regexp.MustCompile(`^([0-9a-f\.\:\/]+)?\s+unicast\s+\[([\w\.:]+)\s+([0-9\-\:\s]+)(?:\s+from\s+([0-9a-f\.\:\/]+))?\]\s+(?:(\*)\s+)?\((\d+)(?:\/\d+)?(?:\/[^\)]*)?\).*$`)
	regex.routes.gatewayBird2 = regexp.MustCompile(`^\s+via\s+([0-9a-f\.\:]+)\s+on\s+([\w\.]+)\s*$`)
}

func dirtyContains(l []string, e string) bool {
	for _, c := range l {
		if c == e {
			return true
		}
	}

	return false
}

func emptyString(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

func specialLine(line string) bool {
	return (strings.HasPrefix(line, "BIRD") || strings.HasPrefix(line, "Access restricted"))
}

func parseStatus(reader io.Reader) Parsed {
	res := Parsed{}

	lines := newLineIterator(reader, true)
	for lines.next() {
		line := lines.string()

		if regex.status.startLine.MatchString(line) {
			res["version"] = regex.status.startLine.FindStringSubmatch(line)[1]
		} else if regex.status.routerID.MatchString(line) {
			res["router_id"] = regex.status.routerID.FindStringSubmatch(line)[1]
		} else if regex.status.currentServer.MatchString(line) {
			res["current_server"] = regex.status.currentServer.FindStringSubmatch(line)[1]
		} else if regex.status.lastReboot.MatchString(line) {
			res["last_reboot"] = regex.status.lastReboot.FindStringSubmatch(line)[1]
		} else if regex.status.lastReconfig.MatchString(line) {
			res["last_reconfig"] = regex.status.lastReconfig.FindStringSubmatch(line)[1]
		} else {
			res["message"] = line
		}
	}

	for k := range res {
		if dirtyContains(ParserConf.FilterFields, k) {
			res[k] = nil
		}
	}

	return Parsed{"status": res}
}

func parseProtocolsShort(reader io.Reader) Parsed {
	res := Parsed{}

	lines := newLineIterator(reader, false)
	for lines.next() {
		line := lines.string()

		if specialLine(line) {
			continue
		}

		if regex.protocol.short.MatchString(line) {
			// The header is skipped, because the regular expression does not
			// match if the "since" field does not contain digits
			matches := regex.protocol.short.FindStringSubmatch(line)

			res[matches[1]] = Parsed{
				"proto": matches[2],
				"table": matches[3],
				"state": matches[4],
				"since": matches[5],
				"info":  matches[6],
			}
		}
	}

	return Parsed{"protocols": res}
}

func parseProtocols(reader io.Reader) Parsed {
	res := Parsed{}

	proto := ""

	lines := newLineIterator(reader, false)
	for lines.next() {
		line := lines.string()

		if emptyString(line) {
			if !emptyString(proto) {
				parsed := parseProtocol(proto)

				res[parsed["protocol"].(string)] = parsed
			}
			proto = ""
		} else {
			proto += (line + "\n")
		}
	}

	return Parsed{"protocols": res}
}

func parseSymbols(reader io.Reader) Parsed {
	res := Parsed{}

	lines := newLineIterator(reader, true)
	for lines.next() {
		line := lines.string()

		if specialLine(line) {
			continue
		}

		if regex.symbols.keyRx.MatchString(line) {
			groups := regex.symbols.keyRx.FindStringSubmatch(line)

			if _, ok := res[groups[2]]; !ok {
				res[groups[2]] = []string{}
			}

			res[groups[2]] = append(res[groups[2]].([]string), groups[1])
		}
	}

	return Parsed{"symbols": res}
}

type blockJob struct {
	lines    []string
	position int
}

type blockParsed struct {
	items    []Parsed
	position int
}

func parseRoutes(reader io.Reader) Parsed {
	jobs := make(chan blockJob)
	out := startRouteWorkers(jobs)

	res := startRouteConsumer(out)
	defer close(res)

	pos := 0
	block := []string{}
	lines := newLineIterator(reader, true)

	for lines.next() {
		line := lines.string()

		if line[0] != 32 && line[0] != 9 && len(block) > 0 {
			jobs <- blockJob{block, pos}
			pos++
			block = []string{}
		}

		block = append(block, line)
	}

	if len(block) > 0 {
		jobs <- blockJob{block, pos}
	}

	close(jobs)

	return <-res
}

func startRouteWorkers(jobs chan blockJob) chan blockParsed {
	out := make(chan blockParsed)

	wg := &sync.WaitGroup{}
	wg.Add(WorkerPoolSize)
	go func() {
		for i := 0; i < WorkerPoolSize; i++ {
			go workerForRouteBlockParsing(jobs, out, wg)
		}
		wg.Wait()
		close(out)
	}()

	return out
}

func startRouteConsumer(out <-chan blockParsed) chan Parsed {
	res := make(chan Parsed)

	go func() {
		byBlock := map[int][]Parsed{}
		count := 0
		for r := range out {
			count++
			byBlock[r.position] = r.items
		}
		res <- Parsed{"routes": sortedSliceForRouteBlocks(byBlock, count)}
	}()

	return res
}

func sortedSliceForRouteBlocks(byBlock map[int][]Parsed, numBlocks int) []Parsed {
	res := []Parsed{}

	for i := 0; i < numBlocks; i++ {
		routes, ok := byBlock[i]
		if !ok {
			continue
		}

		res = append(res, routes...)
	}

	return res
}

func workerForRouteBlockParsing(jobs <-chan blockJob, out chan<- blockParsed, wg *sync.WaitGroup) {
	for j := range jobs {
		parseRouteLines(j.lines, j.position, out)
	}
	wg.Done()
}

func parseRouteLines(lines []string, position int, ch chan<- blockParsed) {
	route := Parsed{}
	routes := []Parsed{}

	for i := 0; i < len(lines); {
		line := lines[i]

		if specialLine(line) {
			i++
			continue
		}

		if regex.routes.prefixBird2.MatchString(line) {
			formerPrefix := ""
			if len(route) > 0 {
				routes = append(routes, route)

				formerPrefix = route["network"].(string)
				route = Parsed{}
			}

			parseMainRouteDetailBird2(regex.routes.prefixBird2.FindStringSubmatch(line), route, formerPrefix)
		} else if regex.routes.startDefinition.MatchString(line) {
			if len(route) > 0 {
				routes = append(routes, route)
				route = Parsed{}
			}

			parseMainRouteDetail(regex.routes.startDefinition.FindStringSubmatch(line), route)
		} else if regex.routes.gatewayBird2.MatchString(line) {
			parseRoutesGatewayBird2(regex.routes.gatewayBird2.FindStringSubmatch(line), route)
		} else if regex.routes.second.MatchString(line) {
			routes = append(routes, route)

			route = parseRoutesSecond(line, route)
		} else if regex.routes.routeType.MatchString(line) {
			submatch := regex.routes.routeType.FindStringSubmatch(line)[1]
			route["type"] = strings.Split(submatch, " ")
		} else if regex.routes.bgp.MatchString(line) {
			// BIRD has a static buffer to hold information which is sent to the client (birdc)
			// If there is more information to be sent to the client than the buffer can hold,
			// the output is split into multiple lines and the continuation of the previous
			// line is indicated by 2 tab characters at the beginning of the next line
			joinLines := func() {
				for c := i + 1; c < len(lines); c++ {
					if strings.HasPrefix(lines[c], "\x09\x09") {
						line += lines[c][2:]
						i++
					} else {
						break
					}
				}
			}

			// The aforementioned behaviour was only observed for the *community fields
			if strings.HasPrefix(line, "\x09BGP.community") ||
				strings.HasPrefix(line, "\x09BGP.large_community") ||
				strings.HasPrefix(line, "\x09BGP.ext_community") {
				joinLines()
			}

			bgp := Parsed{}
			if tmp, ok := route["bgp"]; ok {
				if val, ok := tmp.(Parsed); ok {
					bgp = val
				}
			}

			parseRoutesBgp(line, bgp)
			route["bgp"] = bgp
		}

		i++
	}

	if len(route) > 0 {
		routes = append(routes, route)
	}

	ch <- blockParsed{routes, position}
}

func parseMainRouteDetail(groups []string, route Parsed) {
	route["network"] = groups[1]
	route["gateway"] = groups[2]
	route["interface"] = groups[3]
	route["from_protocol"] = groups[4]
	route["age"] = groups[5]
	route["learnt_from"] = groups[6]
	route["primary"] = groups[7] == "*"
	route["metric"] = parseInt(groups[8])

	for k := range route {
		if dirtyContains(ParserConf.FilterFields, k) {
			route[k] = nil
		}
	}
}

func parseMainRouteDetailBird2(groups []string, route Parsed, formerPrefix string) {
	if len(groups[1]) > 0 {
		route["network"] = groups[1]
	} else {
		route["network"] = formerPrefix
	}

	route["from_protocol"] = groups[2]
	route["age"] = groups[3]
	route["learnt_from"] = groups[4]
	route["primary"] = groups[5] == "*"
	route["metric"] = parseInt(groups[6])

	for k := range route {
		if dirtyContains(ParserConf.FilterFields, k) {
			route[k] = nil
		}
	}
}

func parseRoutesGatewayBird2(groups []string, route Parsed) {
	route["gateway"] = groups[1]
	route["interface"] = groups[2]
}

func parseRoutesSecond(line string, route Parsed) Parsed {
	tmp, ok := route["network"]
	if !ok {
		return route
	}

	var network string
	if network, ok = tmp.(string); !ok {
		return route
	}

	route = Parsed{}
	groups := regex.routes.second.FindStringSubmatch(line)
	first, groups := groups[0], groups[1:]
	groups = append([]string{network}, groups...)
	groups = append([]string{first}, groups...)

	parseMainRouteDetail(groups, route)
	return route
}

func parseRoutesBgp(line string, bgp Parsed) {
	groups := regex.routes.bgp.FindStringSubmatch(line)

	if groups[1] == "community" {
		parseRoutesCommunities(groups, bgp)
	} else if groups[1] == "large_community" {
		parseRoutesLargeCommunities(groups, bgp)
	} else if groups[1] == "ext_community" {
		parseRoutesExtendedCommunities(groups, bgp)
	} else if groups[1] == "as_path" {
		bgp["as_path"] = strings.Split(groups[2], " ")
	} else {
		bgp[groups[1]] = groups[2]
	}
}

func parseRoutesCommunities(groups []string, res Parsed) {
	communities := [][]int64{}
	for _, community := range regex.routes.origin.FindAllString(groups[2], -1) {
		if regex.routes.community.MatchString(community) {
			communityGroups := regex.routes.community.FindStringSubmatch(community)
			maj := parseInt(communityGroups[1])
			min := parseInt(communityGroups[2])
			communities = append(communities, []int64{maj, min})
		}
	}

	res["communities"] = communities
}

func parseRoutesLargeCommunities(groups []string, res Parsed) {
	communities := [][]int64{}
	for _, community := range regex.routes.origin.FindAllString(groups[2], -1) {
		if regex.routes.largeCommunity.MatchString(community) {
			communityGroups := regex.routes.largeCommunity.FindStringSubmatch(community)
			maj := parseInt(communityGroups[1])
			min := parseInt(communityGroups[2])
			pat := parseInt(communityGroups[3])
			communities = append(communities, []int64{maj, min, pat})
		}
	}

	res["large_communities"] = communities
}

func parseRoutesExtendedCommunities(groups []string, res Parsed) {
	communities := []interface{}{}
	for _, community := range regex.routes.origin.FindAllString(groups[2], -1) {
		if regex.routes.extendedCommunity.MatchString(community) {
			communityGroups := regex.routes.extendedCommunity.FindStringSubmatch(community)
			communities = append(communities, []interface{}{communityGroups[1], communityGroups[2], communityGroups[3]})
		}
	}

	res["ext_communities"] = communities
}

func parseRoutesCount(reader io.Reader) Parsed {
	res := Parsed{}

	lines := newLineIterator(reader, true)
	for lines.next() {
		line := lines.string()

		if specialLine(line) {
			continue
		}

		if regex.routeCount.countRx.MatchString(line) {
			count := regex.routeCount.countRx.FindStringSubmatch(line)[1]
			res["routes"] = parseInt(count)
		}
	}

	return res
}

func isCorrectChannel(currentIPVersion string) bool {
	if len(currentIPVersion) == 0 {
		return true
	}

	if getBirdVersion() == 2 {
		return currentIPVersion == "4" || currentIPVersion == "6"
	}

	return currentIPVersion == IPVersion
}

func parseProtocol(lines string) Parsed {
	res := Parsed{}
	routeChanges := Parsed{}

	handlers := []func(string) bool{
		func(l string) bool { return parseProtocolHeader(l, res) },
		func(l string) bool { return parseProtocolRouteLine(l, res) },
		func(l string) bool { return parseProtocolRouteChanges(l, routeChanges) },
		func(l string) bool { return parseProtocolNumberValuesRx(l, res) },
		func(l string) bool { return parseProtocolStringValuesRx(l, res) },
	}

	ipVersion := ""

	reader := strings.NewReader(lines)
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()

		if m := regex.protocol.channel.FindStringSubmatch(line); len(m) > 0 {
			ipVersion = m[1]
		}

		if isCorrectChannel(ipVersion) {
			parseLine(line, handlers)
		}
	}

	res["route_changes"] = routeChanges

	if _, ok := res["routes"]; !ok {
		routes := Parsed{}
		routes["accepted"] = int64(0)
		routes["filtered"] = int64(0)
		routes["imported"] = int64(0)
		routes["exported"] = int64(0)
		routes["preferred"] = int64(0)

		res["routes"] = routes
	}

	return res
}

func parseLine(line string, handlers []func(string) bool) {
	for _, h := range handlers {
		if h(line) {
			return
		}
	}
}

func parseProtocolHeader(line string, res Parsed) bool {
	groups := regex.protocol.protocol.FindStringSubmatch(line)
	if groups == nil {
		return false
	}

	res["protocol"] = groups[1]
	res["bird_protocol"] = groups[2]
	res["table"] = groups[3]
	res["state"] = groups[4]
	res["state_changed"] = groups[5]
	res["connection"] = groups[6] // TODO eliminate
	if groups[2] == "Pipe" {
		res["peer_table"] = groups[6][3:]
	}
	return true
}

func parseProtocolRouteLine(line string, res Parsed) bool {
	groups := regex.protocol.routes.FindStringSubmatch(line)
	if groups == nil {
		return false
	}

	routes := parseProtocolRoutes(groups[1])
	res["routes"] = routes
	return true
}

func setChangeCount(name string, value string, res Parsed) {
	if value == "---" { // field not available for protocol
		return
	}
	res[name] = parseInt(value)
}

func parseProtocolRouteChanges(line string, res Parsed) bool {
	groups := regex.protocol.routeChanges.FindStringSubmatch(line)
	if groups == nil {
		return false
	}

	updates := Parsed{}
	setChangeCount("received", groups[3], updates)
	setChangeCount("rejected", groups[4], updates)
	setChangeCount("filtered", groups[5], updates)
	setChangeCount("ignored", groups[6], updates)
	setChangeCount("accepted", groups[7], updates)

	key := strings.ToLower(groups[1]) + "_" + groups[2]
	res[key] = updates
	return true
}

func parseProtocolNumberValuesRx(line string, res Parsed) bool {
	groups := regex.protocol.numericValue.FindStringSubmatch(line)
	if groups == nil {
		return false
	}

	key := treatKey(groups[1])
	res[key] = parseInt(groups[2])
	return true
}

func parseProtocolStringValuesRx(line string, res Parsed) bool {
	groups := regex.protocol.stringValue.FindStringSubmatch(line)
	if groups == nil {
		return false
	}

	key := treatKey(groups[1])
	res[key] = groups[2]
	return true
}

// Will snake_case a value like that:
// I am a Weird stRiNg -> i_am_a_weird_string
func treatKey(key string) string {
	spaces := regexp.MustCompile(`\s+`)
	key = spaces.ReplaceAllString(key, "_")
	return strings.ToLower(key)
}

func parseInt(from string) int64 {
	val, err := strconv.ParseInt(from, 10, 64)
	if err != nil {
		return int64(0)
	}

	return val
}

func parseProtocolRoutes(input string) Parsed {
	routes := Parsed{}

	// Input: 1 imported, 0 filtered, 2 exported, 1 preferred
	tokens := strings.Split(input, ",")
	for _, token := range tokens {
		token = strings.TrimSpace(token)
		content := strings.Split(token, " ")
		routes[content[1]] = parseInt(content[0])
	}

	return routes
}
