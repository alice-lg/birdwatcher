package bird

import (
  "regexp"
)

type Parsed map[string]interface{}

func parseStatus(input []byte) Parsed {
  res := Parsed{}
  line_sep := regexp.MustCompile(`((\r?\n)|(\r\n?))`)
  lines := line_sep.Split(string(input), -1)

  start_line_rx := regexp.MustCompile(`^BIRD\s([0-9\.]+)\s*$`)
  router_id_rx := regexp.MustCompile(`^Router\sID\sis\s([0-9\.]+)\s*$`)
  current_server_rx := regexp.MustCompile(`^Current\sserver\stime\sis\s([0-9\-]+)\s([0-9\:]+)\s*$`)
  last_reboot_rx := regexp.MustCompile(`^Last\sreboot\son\s([0-9\-]+)\s([0-9\:]+)\s*$`)
  last_reconfig_rx := regexp.MustCompile(`^Last\sreconfiguration\son\s([0-9\-]+)\s([0-9\:]+)\s*$`)
  for _, line := range lines {
    if (start_line_rx.MatchString(line)) {
      res["version"] = start_line_rx.FindStringSubmatch(line)[1]
    } else if (router_id_rx.MatchString(line)) {
      res["router_id"] = router_id_rx.FindStringSubmatch(line)[1]
    } else if (current_server_rx.MatchString(line)) {
      res["current_server"] = current_server_rx.FindStringSubmatch(line)[1]
    } else if (last_reboot_rx.MatchString(line)) {
      res["last_reboot"] = last_reboot_rx.FindStringSubmatch(line)[1]
    } else if (last_reconfig_rx.MatchString(line)) {
      res["last_reconfig"] = last_reconfig_rx.FindStringSubmatch(line)[1]
    } else {
      res["message"] = line
    }
  }
  return res
}

func parseProtocols(input []byte) Parsed {
  return Parsed{}
}

func parseSymbols(input []byte) Parsed {
  return Parsed{}
}

func parseRoutes(input []byte) Parsed {
  return Parsed{}
}

func parseRoutesCount(input []byte) Parsed {
  return Parsed{}
}
