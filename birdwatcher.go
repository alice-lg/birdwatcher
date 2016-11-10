package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"log/syslog"
	"net/http"
	"os"
	"regexp"

	"github.com/julienschmidt/httprouter"
	yaml "gopkg.in/yaml.v2"
)

var debug int = 0
var slog *syslog.Writer // Our syslog connection
var conf *Config

type Match struct {
	Expr   string   // The regular expression as a string.
	Fields []string // The named fields for grouped expressions.
	Next   string   // The next regular expression in the flow.
	Action string   // What to do with the stored fields: "store" or "send".
}

// Compiled regular expression and it's corresponding match data.
type RE struct {
	RE    *regexp.Regexp
	Match Match
}

// The configuration found in the configuration file.
type FileConfig struct {
	Matches  map[string]Match // All our regular expressions and related data.
	Listen   string           // Listen to this address:port for HTTP.
	FileName string           // File to look for patterns

}

type Config struct {
	Conf FileConfig
	Res  map[string]RE
}

// Parse the configuration file. Returns the configuration.
func parseconfig(filename string) (conf *Config, err error) {
	conf = new(Config)

	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}

	if err = yaml.Unmarshal(contents, &conf.Conf); err != nil {
		return
	}

	conf.Res = make(map[string]RE)

	// Build the regexps from the configuration.
	for key, match := range conf.Conf.Matches {
		var err error
		var re RE

		re.Match = match
		re.RE, err = regexp.Compile(match.Expr)
		if err != nil {
			slog.Err("Couldn't compile re: " + match.Expr)
			os.Exit(-1)
		}

		// Check that the number of capturing groups matches the number of expected fields.
		lengroups := len(re.RE.SubexpNames()) - 1
		lenfields := len(re.Match.Fields)

		if lengroups != lenfields {
			line := fmt.Sprintf("Number of capturing groups (%v) not equal to number of fields (%v): %s", lengroups, lenfields, re.Match.Expr)
			slog.Err(line)
			os.Exit(-1)
		}

		conf.Res[key] = re
	}

	return
}

func main() {
	var configfile = flag.String("config", "birdwatcher.yaml", "Path to configuration file")
	var flagdebug = flag.Int("debug", 0, "Be more verbose")
	flag.Parse()

	debug = *flagdebug

	slog, err := syslog.New(syslog.LOG_ERR, "birdwatcher")
	if err != nil {
		fmt.Printf("Couldn't open syslog")
		os.Exit(-1)
	}

	slog.Debug("birdwatcher starting")

	config, err := parseconfig(*configfile)
	if err != nil {
		slog.Err("Couldn't parse configuration file: " + err.Error())
		os.Exit(-1)
	}
	conf = config

	fmt.Printf("%v\n", conf)

	r := httprouter.New()
	r.GET("/status", Status)
	r.GET("/routes", Routes)
	r.GET("/protocols", Protocols)

	log.Fatal(http.ListenAndServe(":29184", r))
}
