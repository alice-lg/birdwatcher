package main

import (
	"bufio"
	"fmt"
	"os"
)

// readLines reads a whole file into memory
// and returns a slice of its lines.
func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

// pattern looks for pattern matching regular expression re in the
// lines buffer. Returns the name of our next regular expression to
// look for.
func pattern(first string, lines []string) (maplist []map[string]string) {
	var re RE
	var fieldmap map[string]string

	fmt.Printf("In pattern\n")

	re, ok := conf.Res[first]
	if !ok {
		slog.Debug("Couldn't find first state.")
		return
	}

	// Store away all the fields in a map.
	fieldmap = make(map[string]string)

	for _, line := range lines {
		if debug > 1 {
			slog.Debug("Looking at: " + line)
		}

		fields := re.RE.FindStringSubmatch(line)

		if len(fields)-1 == len(re.Match.Fields) {
			if debug > 0 {
				line := fmt.Sprintf("Found match for a message of type '%v'", re.Match)
				slog.Debug(line)
			}
			for key, name := range re.Match.Fields {
				if debug > 0 {
					slog.Debug("Got " + re.Match.Fields[key])
				}
				fieldmap[name] = fields[key+1]
			}

			if re.Match.Action == "send" {
				// Finished for this item. Create a new field map and add this to the list.
				maplist = append(maplist, fieldmap)
				fieldmap = make(map[string]string)
			}

			// Go to the next state, if it exists. If it
			// doesn't we're finished here.
			re, ok = conf.Res[re.Match.Next]
			if !ok {
				break
			}
		}
	}

	return maplist
}
