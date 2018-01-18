package bird

import (
	"bufio"
	"io"
)

type lineIterator struct {
	scanner        *bufio.Scanner
	skipEmptyLines bool
}

func newLineIterator(reader io.Reader, skipEmptyLines bool) *lineIterator {
	scanner := bufio.NewScanner(reader)
	return &lineIterator{scanner: scanner, skipEmptyLines: skipEmptyLines}
}

func (l *lineIterator) next() bool {
	res := l.scanner.Scan()
	if !res || !l.skipEmptyLines {
		return res
	}

	if emptyString(l.scanner.Text()) {
		return l.next()
	}

	return res
}

func (l *lineIterator) string() string {
	return l.scanner.Text()
}
