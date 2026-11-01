package util

import (
	"bufio"
	"log"
	"os"
	"strings"
)

func CleanLists(paths []string) (err error) {
	var seen = map[string]struct{}{}

	for _, p := range paths {
		var dupes int
		dupes, err = cleanFile(p, seen)
		if err != nil {
			return
		}
		log.Printf("Cleaned %s: removed %d duplicate rules\n", p, dupes)
	}

	return nil
}

func cleanFile(path string, seen map[string]struct{}) (dupes int, err error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return
	}

	var (
		out     strings.Builder
		scanner = bufio.NewScanner(strings.NewReader(string(content)))
	)
	scanner.Buffer(make([]byte, 0, 64*1024), 16*1024*1024)

	for scanner.Scan() {
		line := strings.TrimRight(scanner.Text(), "\r")

		trimmed := strings.TrimSpace(line)
		isRule := trimmed != "" && !strings.HasPrefix(trimmed, "!") && !strings.HasPrefix(trimmed, "[")

		if !isRule {
			out.WriteString(line)
			out.WriteByte('\n')
			continue
		}

		if _, ok := seen[trimmed]; ok {
			dupes++
			continue
		}
		seen[trimmed] = struct{}{}

		out.WriteString(trimmed)
		out.WriteByte('\n')
	}
	if err = scanner.Err(); err != nil {
		return
	}

	err = os.WriteFile(path, []byte(out.String()), 0o644)
	return
}
