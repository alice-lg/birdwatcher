package bird

import (
	"bufio"
	"regexp"
	"strconv"
	"strings"
)

var (
	ParserConf ParserConfig
	regex      struct {
		lineSeperator *regexp.Regexp
		status        struct {
			startLine     *regexp.Regexp
			routerID      *regexp.Regexp
			currentServer *regexp.Regexp
			lastReboot    *regexp.Regexp
			lastReconfig  *regexp.Regexp
		}
		bgp struct {
			channel         *regexp.Regexp
			protocol        *regexp.Regexp
			numericValue    *regexp.Regexp
			routes          *regexp.Regexp
			stringValue     *regexp.Regexp
			importUpdates   *regexp.Regexp
			importWithdraws *regexp.Regexp
			exportUpdates   *regexp.Regexp
			exportWithdraws *regexp.Regexp
		}
		symbols struct {
			keyRx *regexp.Regexp
		}
		routeCount struct {
			countRx *regexp.Regexp
		}
		routes struct {
			startDefinition *regexp.Regexp
			second          *regexp.Regexp
			routeType       *regexp.Regexp
			bgp             *regexp.Regexp
			community       *regexp.Regexp
			largeCommunity  *regexp.Regexp
			origin          *regexp.Regexp
		}
	}
)

type Parsed map[string]interface{}

func init() {
	regex.status.startLine = regexp.MustCompile(`^BIRD\s([0-9\.]+)\s*$`)
	regex.status.routerID = regexp.MustCompile(`^Router\sID\sis\s([0-9\.]+)\s*$`)
	regex.status.currentServer = regexp.MustCompile(`^Current\sserver\stime\sis\s([0-9\-]+\s[0-9\:]+)\s*$`)
	regex.status.lastReboot = regexp.MustCompile(`^Last\sreboot\son\s([0-9\-]+\s[0-9\:]+)\s*$`)
	regex.status.lastReconfig = regexp.MustCompile(`^Last\sreconfiguration\son\s([0-9\-]+\s[0-9\:]+)\s*$`)

	regex.symbols.keyRx = regexp.MustCompile(`^([^\s]+)\s+(.+)\s*$`)

	regex.routeCount.countRx = regexp.MustCompile(`^(\d+)\s+of\s+(\d+)\s+routes.*$`)

	regex.bgp.channel = regexp.MustCompile("Channel ipv([46])")
	regex.bgp.protocol = regexp.MustCompile(`^([\w\.:]+)\s+BGP\s+(\w+)\s+(\w+)\s+([0-9]{4}-[0-9]{2}-[0-9]{2}\s+[0-9]{2}:[0-9]{2}:[0-9]{2})\s*(\w+)?.*$`)
	regex.bgp.numericValue = regexp.MustCompile(`^\s+([^:]+):\s+([\d]+)\s*$`)
	regex.bgp.routes = regexp.MustCompile(`^\s+Routes:\s+(.*)`)
	regex.bgp.stringValue = regexp.MustCompile(`^\s+([^:]+):\s+(.+)\s*$`)
	regex.bgp.importUpdates = regexp.MustCompile(`^\s+Import updates:\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s*$`)
	regex.bgp.importWithdraws = regexp.MustCompile(`^\s+Import withdraws:\s+(\d+)\s+(\d+)\s+\-\-\-\s+(\d+)\s+(\d+)\s*$`)
	regex.bgp.exportUpdates = regexp.MustCompile(`^\s+Export updates:\s+(\d+)\s+(\d+)\s+(\d+)\s+\-\-\-\s+(\d+)\s*$`)
	regex.bgp.exportWithdraws = regexp.MustCompile(`^\s+Export withdraws:\s+(\d+)(\s+\-\-\-){2}\s+(\d+)\s*$`)

	regex.routes.startDefinition = regexp.MustCompile(`^([0-9a-f\.\:\/]+)\s+via\s+([0-9a-f\.\:]+)\s+on\s+([\w\.]+)\s+\[([\w\.:]+)\s+([0-9\-\:\s]+)(?:\s+from\s+([0-9a-f\.\:\/]+)){0,1}\]\s+(?:(\*)\s+){0,1}\((\d+)(?:\/\d+){0,1}\).*`)
	regex.routes.second = regexp.MustCompile(`^\s+via\s+([0-9a-f\.\:]+)\s+on\s+([\w\.]+)\s+\[([\w\.:]+)\s+([0-9\-\:\s]+)(?:\s+from\s+([0-9a-f\.\:\/]+)){0,1}\]\s+(?:(\*)\s+){0,1}\((\d+)(?:\/\d+){0,1}\).*$`)
	regex.routes.routeType = regexp.MustCompile(`^\s+Type:\s+(.*)\s*$`)
	regex.routes.bgp = regexp.MustCompile(`^\s+BGP.(\w+):\s+(.+)\s*$`)
	regex.routes.community = regexp.MustCompile(`^\((\d+),\s*(\d+)\)`)
	regex.routes.largeCommunity = regexp.MustCompile(`^\((\d+),\s*(\d+),\s*(\d+)\)`)
	regex.routes.origin = regexp.MustCompile(`\([^\(]*\)\s*`)
}

func dirtyContains(l []string, e string) bool {
	for _, c := range l {
		if c == e {
			return true
		}
	}

	return false
}

func emptyLine(line string) bool {
	return len(strings.TrimSpace(line)) == 0
}

func getLinesUnfiltered(input string) []string {
	lines := make([]string, 0)

	reader := strings.NewReader(input)
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines
}

func getLinesFromString(input string) []string {
	lines := getLinesUnfiltered(input)

	var filtered []string

	for _, line := range lines {
		if !emptyLine(line) {
			filtered = append(filtered, line)
		}
	}

	return filtered
}

func getLines(input []byte) []string {
	return getLinesFromString(string(input))
}

func specialLine(line string) bool {
	return (strings.HasPrefix(line, "BIRD") || strings.HasPrefix(line, "Access restricted"))
}

func parseStatus(input []byte) Parsed {
	res := Parsed{}
	lines := getLines(input)

	for _, line := range lines {
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

func parseProtocols(input []byte) Parsed {
	res := Parsed{}
	protocols := []string{}
	lines := getLinesUnfiltered(string(input))

	proto := ""
	for _, line := range lines {
		if emptyLine(line) {
			if !emptyLine(proto) {
				protocols = append(protocols, proto)
			}
			proto = ""
		} else {
			proto += (line + "\n")
		}
	}

	res["protocols"] = protocols
	return res
}

func parseSymbols(input []byte) Parsed {
	res := Parsed{}
	lines := getLines(input)

	for _, line := range lines {
		if specialLine(line) {
			continue
		}

		if regex.symbols.keyRx.MatchString(line) {
			groups := regex.symbols.keyRx.FindStringSubmatch(line)
			res[groups[2]] = groups[1]
		}
	}

	return Parsed{"symbols": res}
}

func parseRoutes(input []byte) Parsed {
	res := Parsed{}
	lines := getLines(input)

	routes := []Parsed{}
	route := Parsed{}

	for _, line := range lines {
		if specialLine(line) || (len(route) == 0 && emptyLine(line)) {
			continue
		}

		if regex.routes.startDefinition.MatchString(line) {
			if len(route) > 0 {
				routes = append(routes, route)
				route = Parsed{}
			}

			route = parseMainRouteDetail(regex.routes.startDefinition.FindStringSubmatch(line), route)
		} else if regex.routes.second.MatchString(line) {
			routes = append(routes, route)

			route = parseRoutesSecond(line, route)
		} else if regex.routes.routeType.MatchString(line) {
			submatch := regex.routes.routeType.FindStringSubmatch(line)[1]
			route["type"] = strings.Split(submatch, " ")
		} else if regex.routes.bgp.MatchString(line) {
			bgp := Parsed{}
			if tmp, ok := route["bgp"]; ok {
				if val, ok := tmp.(Parsed); ok {
					bgp = val
				}
			}

			parseRoutesBgp(line, bgp)
			route["bgp"] = bgp
		}
	}

	if len(route) > 0 {
		routes = append(routes, route)
	}

	res["routes"] = routes
	return res
}

func parseMainRouteDetail(groups []string, route Parsed) Parsed {
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

	return route
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

	return parseMainRouteDetail(groups, route)
}

func parseRoutesBgp(line string, bgp Parsed) {
	groups := regex.routes.bgp.FindStringSubmatch(line)

	if groups[1] == "community" {
		parseRoutesCommunities(groups, bgp)
	} else if groups[1] == "large_community" {
		parseRoutesLargeCommunities(groups, bgp)
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

func parseRoutesCount(input []byte) Parsed {
	res := Parsed{}
	lines := getLines(input)

	for _, line := range lines {
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

	return currentIPVersion == IPVersion
}

func parseBgp(input string) Parsed {
	res := Parsed{}
	routeChanges := Parsed{}

	handlers := []func(string) bool{
		func(l string) bool { return parseBgpProtocol(l, res) },
		func(l string) bool { return parseBgpRouteLine(l, res) },
		func(l string) bool { return parseBgpImportUpdates(l, routeChanges) },
		func(l string) bool { return parseBgpImportWithdraws(l, routeChanges) },
		func(l string) bool { return parseBgpExportUpdates(l, routeChanges) },
		func(l string) bool { return parseBgpExportWithdraws(l, routeChanges) },
		func(l string) bool { return parseBgpNumberValuesRx(l, res) },
		func(l string) bool { return parseBgpStringValuesRx(l, res) },
	}

	lines := getLinesFromString(input)
	ipVersion := ""
	for _, line := range lines {
		if m := regex.bgp.channel.FindStringSubmatch(line); len(m) > 0 {
			ipVersion = m[1]
		}

		if isCorrectChannel(ipVersion) {
			parseLine(line, handlers)
		}
	}

	res["route_changes"] = routeChanges

	if _, ok := res["routes"]; !ok {
		routes := Parsed{}
		routes["accepted"] = 0
		routes["filtered"] = 0
		routes["exported"] = 0
		routes["preferred"] = 0

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

func parseBgpProtocol(line string, res Parsed) bool {
	groups := regex.bgp.protocol.FindStringSubmatch(line)
	if groups == nil {
		return false
	}

	res["protocol"] = groups[1]
	res["bird_protocol"] = "BGP"
	res["table"] = groups[2]
	res["state"] = groups[3]
	res["state_changed"] = groups[4]
	res["connection"] = groups[5]
	return true
}

func parseBgpRouteLine(line string, res Parsed) bool {
	groups := regex.bgp.routes.FindStringSubmatch(line)
	if groups == nil {
		return false
	}

	routes := parseBgpRoutes(groups[1])
	res["routes"] = routes
	return true
}

func parseBgpImportUpdates(line string, res Parsed) bool {
	groups := regex.bgp.importUpdates.FindStringSubmatch(line)
	if groups == nil {
		return false
	}

	updates := Parsed{}
	updates["received"] = parseInt(groups[1])
	updates["rejected"] = parseInt(groups[2])
	updates["filtered"] = parseInt(groups[3])
	updates["ignored"] = parseInt(groups[4])
	updates["accepted"] = parseInt(groups[5])

	res["import_updates"] = updates
	return true
}

func parseBgpImportWithdraws(line string, res Parsed) bool {
	groups := regex.bgp.importWithdraws.FindStringSubmatch(line)
	if groups == nil {
		return false
	}

	updates := Parsed{}
	updates["received"] = parseInt(groups[1])
	updates["rejected"] = parseInt(groups[2])
	updates["filtered"] = parseInt(groups[3])
	updates["accepted"] = parseInt(groups[4])

	res["import_withdraws"] = updates
	return true
}

func parseBgpExportUpdates(line string, res Parsed) bool {
	groups := regex.bgp.exportUpdates.FindStringSubmatch(line)
	if groups == nil {
		return false
	}

	updates := Parsed{}
	updates["received"] = parseInt(groups[1])
	updates["rejected"] = parseInt(groups[2])
	updates["ignored"] = parseInt(groups[3])
	updates["accepted"] = parseInt(groups[4])

	res["export_updates"] = updates
	return true
}

func parseBgpExportWithdraws(line string, res Parsed) bool {
	groups := regex.bgp.exportWithdraws.FindStringSubmatch(line)
	if groups == nil {
		return false
	}

	updates := Parsed{}
	updates["received"] = parseInt(groups[1])
	updates["accepted"] = parseInt(groups[3])

	res["export_withdraws"] = updates
	return true
}

func parseBgpNumberValuesRx(line string, res Parsed) bool {
	groups := regex.bgp.numericValue.FindStringSubmatch(line)
	if groups == nil {
		return false
	}

	key := treatKey(groups[1])
	res[key] = parseInt(groups[2])
	return true
}

func parseBgpStringValuesRx(line string, res Parsed) bool {
	groups := regex.bgp.stringValue.FindStringSubmatch(line)
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
		return 0
	}

	return val
}

func parseBgpRoutes(input string) Parsed {
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
