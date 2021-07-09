package util

import (
	"bufio"
	"log"
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

		// Remove comments and empty lines
		if strings.HasPrefix(t, "#") || t == "" {
			continue
		}

		// Warn about invalid URLs
		if _, uerr := url.ParseRequestURI(t); uerr != nil {
			log.Printf("Invalid URL %q\n", t)
			continue
		}

		if entriesMap[t] {
			// We have seen this URL before. Just warn about it, nothing else
			log.Printf("Duplicate URL %q will only be downloaded once\n", t)
		} else {
			// If we haven't seen this URL before, we add it to the list
			entriesMap[t] = true
		}
	}

	for url := range entriesMap {
		entries = append(entries, url)
	}

	sort.Strings(entries)

	return
}
