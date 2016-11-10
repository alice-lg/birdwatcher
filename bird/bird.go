package bird

import (
	"os/exec"
	"strings"
)

func Run(args string) ([]byte, error) {
	args = "show " + args
	argsList := strings.Split(args, " ")
	return exec.Command("birdc", argsList...).Output()
}

func RunAndParse(cmd string, parser func([]byte) Parsed) Parsed {
	out, err := Run(cmd)

	if err != nil {
		// ignore errors for now
		return Parsed{}
	}

	return parser(out)
}

func Status() Parsed {
	return RunAndParse("status", parseStatus)
}

func Protocols() Parsed {
	return RunAndParse("protocols all", parseProtocols)
}

func Symbols() Parsed {
	return RunAndParse("symbols", parseSymbols)
}

func RoutesProto(protocol string) Parsed {
	return RunAndParse("route protocol "+protocol+" all",
		parseRoutes)
}

func RoutesProtoCount(protocol string) Parsed {
	return RunAndParse("route protocol "+protocol+" count",
		parseRoutesCount)
}

func RoutesExport(protocol string) Parsed {
	return RunAndParse("route export "+protocol+" all",
		parseRoutes)
}

func RoutesExportCount(protocol string) Parsed {
	return RunAndParse("route export "+protocol+" count",
		parseRoutesCount)
}

func RoutesTable(table string) Parsed {
	return RunAndParse("route table "+table+" all",
		parseRoutes)
}

func RoutesTableCount(table string) Parsed {
	return RunAndParse("route table "+table+" count",
		parseRoutesCount)
}

func RoutesLookupTable(net string, table string) Parsed {
	return RunAndParse("route for "+net+" table "+table+" all",
		parseRoutes)
}

func RoutesLookupProtocol(net string, protocol string) Parsed {
	return RunAndParse("route for "+net+" protocol "+protocol+" all",
		parseRoutes)
}
