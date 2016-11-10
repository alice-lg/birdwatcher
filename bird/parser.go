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

func parseRoutes(input []byte) Parsed {
	return Parsed{}
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
