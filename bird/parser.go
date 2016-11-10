package bird

import (
	"regexp"
	"strconv"
	"strings"
)

type Parsed map[string]interface{}

func emptyLine(line string) bool {
	return len(strings.TrimSpace(line)) == 0
}

func getLinesUnfiltered(input []byte) []string {
	line_sep := regexp.MustCompile(`((\r?\n)|(\r\n?))`)
	return line_sep.Split(string(input), -1)
}

func getLines(input []byte) []string {
	lines := getLinesUnfiltered(input)

	var filtered []string

	for _, line := range lines {
		if !emptyLine(line) {
			filtered = append(filtered, line)
		}
	}

	return filtered
}

func specialLine(line string) bool {
	return (strings.HasPrefix(line, "BIRD") ||
		strings.HasPrefix(line, "Access restricted"))
}

func parseStatus(input []byte) Parsed {
	res := Parsed{}
	lines := getLines(input)

	start_line_rx := regexp.MustCompile(`^BIRD\s([0-9\.]+)\s*$`)
	router_id_rx := regexp.MustCompile(`^Router\sID\sis\s([0-9\.]+)\s*$`)
	current_server_rx := regexp.MustCompile(`^Current\sserver\stime\sis\s([0-9\-]+)\s([0-9\:]+)\s*$`)
	last_reboot_rx := regexp.MustCompile(`^Last\sreboot\son\s([0-9\-]+)\s([0-9\:]+)\s*$`)
	last_reconfig_rx := regexp.MustCompile(`^Last\sreconfiguration\son\s([0-9\-]+)\s([0-9\:]+)\s*$`)

	for _, line := range lines {
		if start_line_rx.MatchString(line) {
			res["version"] = start_line_rx.FindStringSubmatch(line)[1]
		} else if router_id_rx.MatchString(line) {
			res["router_id"] = router_id_rx.FindStringSubmatch(line)[1]
		} else if current_server_rx.MatchString(line) {
			res["current_server"] = current_server_rx.FindStringSubmatch(line)[1]
		} else if last_reboot_rx.MatchString(line) {
			res["last_reboot"] = last_reboot_rx.FindStringSubmatch(line)[1]
		} else if last_reconfig_rx.MatchString(line) {
			res["last_reconfig"] = last_reconfig_rx.FindStringSubmatch(line)[1]
		} else {
			res["message"] = line
		}
	}
	return res
}

func parseProtocols(input []byte) Parsed {
	res := Parsed{}
	protocols := []string{}
	lines := getLinesUnfiltered(input)

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

	key_rx := regexp.MustCompile(`^([^\s]+)\s+(.+)\s*$`)
	for _, line := range lines {
		if specialLine(line) {
			continue
		}

		if key_rx.MatchString(line) {
			groups := key_rx.FindStringSubmatch(line)
			res[groups[2]] = groups[1]
		}
	}

	return res
}

func mainRouteDetail(groups []string, route Parsed) Parsed {
	route["network"] = groups[1]
	route["gateway"] = groups[2]
	route["interface"] = groups[3]
	route["from_protocol"] = groups[4]
	route["age"] = groups[5]
	route["learnt_from"] = groups[6]
	route["primary"] = groups[7] == "*"
	if val, err := strconv.ParseInt(groups[8], 10, 64); err != nil {
		route["metric"] = 0
	} else {
		route["metric"] = val
	}
	return route
}

func parseRoutes(input []byte) Parsed {
	res := Parsed{}
	lines := getLines(input)

	routes := []Parsed{}

	route := Parsed{}
	start_def_rx := regexp.MustCompile(`^([0-9a-f\.\:\/]+)\s+via\s+([0-9a-f\.\:]+)\s+on\s+(\w+)\s+\[(\w+)\s+([0-9\-\:]+)(?:\s+from\s+([0-9a-f\.\:\/]+)){0,1}\]\s+(?:(\*)\s+){0,1}\((\d+)(?:\/\d+){0,1}\).*$`)
	second_rx := regexp.MustCompile(`^\s+via\s+([0-9a-f\.\:]+)\s+on\s+(\w+)\s+\[(\w+)\s+([0-9\-\:]+)(?:\s+from\s+([0-9a-f\.\:\/]+)){0,1}\]\s+(?:(\*)\s+){0,1}\((\d+)(?:\/\d+){0,1}\).*$`)
	type_rx := regexp.MustCompile(`^\s+Type:\s+(.*)\s*$`)
	bgp_rx := regexp.MustCompile(`^\s+BGP.(\w+):\s+(\w+)\s*$`)
	community_rx := regexp.MustCompile(`^\((\d+),(\d+)\)`)
	for _, line := range lines {
		if specialLine(line) || (len(route) == 0 && emptyLine(line)) {
			continue
		}

		if start_def_rx.MatchString(line) {
			if len(route) > 0 {
				routes = append(routes, route)
				route = Parsed{}
			}
			route = mainRouteDetail(start_def_rx.FindStringSubmatch(line), route)
		} else if second_rx.MatchString(line) {
			routes = append(routes, route)
			var network string
			if tmp, ok := route["network"]; ok {
				if val, ok := tmp.(string); ok {
					network = val
				} else {
					continue
				}
			} else {
				continue
			}
			route = Parsed{}

			groups := second_rx.FindStringSubmatch(line)
			first, groups := groups[0], groups[1:]
			groups = append([]string{network}, groups...)
			groups = append([]string{first}, groups...)
			route = mainRouteDetail(groups, route)
		} else if type_rx.MatchString(line) {
			route["type"] = strings.Split(type_rx.FindStringSubmatch(line)[1],
				" ")
		} else if bgp_rx.MatchString(line) {
			groups := bgp_rx.FindStringSubmatch(line)
			bgp := Parsed{}

			if tmp, ok := route["bgp"]; ok {
				if val, ok := tmp.(Parsed); ok {
					bgp = val
				}
			}

			if groups[1] == "community" {
				communities := [][]int64{}
				for _, community := range strings.Split(groups[2], " ") {
					if community_rx.MatchString(community) {
						com_groups := community_rx.FindStringSubmatch(community)
						maj, err := strconv.ParseInt(com_groups[1], 10, 64)
						if err != nil {
							continue
						}
						min, err := strconv.ParseInt(com_groups[2], 10, 64)
						if err != nil {
							continue
						}
						communities = append(communities, []int64{maj, min})
					}
					bgp["communities"] = communities
				}
			} else {
				bgp[groups[1]] = groups[2]
			}

			route["bgp"] = bgp
		}
	}

	if len(route) > 0 {
		routes = append(routes, route)
	}

	res["routes"] = routes
	return res
}

func parseRoutesCount(input []byte) Parsed {
	res := Parsed{}
	lines := getLines(input)

	count_rx := regexp.MustCompile(`^(\d+)\s+of\s+(\d+)\s+routes.*$`)
	for _, line := range lines {
		if specialLine(line) {
			continue
		}

		if count_rx.MatchString(line) {
			i, err := strconv.ParseInt(count_rx.FindStringSubmatch(line)[1], 10, 64)
			if err != nil {
				// ignore for now
				continue
			}
			res["routes"] = i
		}
	}

	return res
}
