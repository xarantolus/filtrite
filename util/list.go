package util

import (
	"bufio"
	"net/url"
	"os"
	"sort"
	"strings"
)

// ReadListFile returns all URLs read from file `name` without duplicates, sorted
func ReadListFile(fn string) (entries []string, err error) {
	var entriesMap = map[string]bool{}

	f, err := os.Open(fn)
	if err != nil {
		return
	}
	defer f.Close()

	scan := bufio.NewScanner(f)

	for scan.Scan() {
		t := strings.TrimSpace(scan.Text())

		// comments
		if strings.HasPrefix(t, "#") || t == "" {
			continue
		}

		if _, uerr := url.ParseRequestURI(t); uerr == nil {
			if !entriesMap[t] {
				entriesMap[t] = true
			}
		}
	}

	for url := range entriesMap {
		entries = append(entries, url)
	}

	sort.Strings(entries)

	return
}
