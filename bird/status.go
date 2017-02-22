package bird

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

// Get last reconfig timestamp from file modification date
func lastReconfigTimestampFromFileStat(filename string) string {
	info, err := os.Stat(filename)
	if err != nil {
		return fmt.Sprintf("Could not fetch file modified timestamp: %s", err)
	}

	modTime := info.ModTime().UTC()
	buf, _ := modTime.MarshalJSON()

	return string(buf)
}

// Parse config file linewise, find matching line and extract date
func lastReconfigTimestampFromFileContent(filename string, regex string) string {
	rx := regexp.MustCompile(regex)

	// Read config file linewise
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Sprintf("Could not read: %s", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		matches := rx.FindStringSubmatch(scanner.Text())
		if len(matches) > 0 {
			return matches[1]
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Sprintf("Error reading config: %s", err)
	}

	return ""
}
