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
func pattern(re RE, lines []string) (next string) {
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
		}
	}

	return re.Match.Next
}
