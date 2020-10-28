package util

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	releaseFile = "release.md"
)

// AppendReleaseList appends
func AppendReleaseList(fn string, gotCount, fullCount int) (err error) {
	f, err := os.OpenFile(releaseFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()

	var (
		errorCount   = fullCount - gotCount
		friendlyName = strings.TrimSuffix(filepath.Base(fn), ".txt")
	)

	if errorCount == 1 {
		fmt.Fprintf(f, "* `%s`: updated %d/%d lists, one error\n", friendlyName, gotCount, fullCount)
	} else {
		fmt.Fprintf(f, "* `%s`: updated %d/%d lists, %d errors\n", friendlyName, gotCount, fullCount, errorCount)
	}

	return
}
