package bird

// Created: 2016-12-01 14:15:00

import (
	"testing"
)

func TestReconfigTimestampFromStat(t *testing.T) {

	// Just get the modification date of this file
	ts := lastReconfigTimestampFromFileStat("./status_test.go")
	t.Log(ts)

	ts = lastReconfigTimestampFromFileStat("./___i_do_not_exist___")
	t.Log(ts)
}

func TestReconfigTimestampFromContent(t *testing.T) {

	ts := lastReconfigTimestampFromFileContent("./status_test.go", "// Created: (.*)")
	t.Log(ts)
}
