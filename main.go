package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"unicode"
	"xarantolus/filtrite/util"
)

const (
	tmpDir  = "tmp"
	listDir = "lists"
	distDir = "dist"
	logDir  = "logs"

	// This is the "kMaxBodySize" from https://github.com/bromite/bromite/blob/master/build/patches/Bromite-AdBlockUpdaterService.patch
	bromiteMaxFilterSize = 1024 * 1024 * 10
)

// generateFilterList generates a filter list from the listTextFile
func generateFilterList(listTextFile string) (err error) {

	var listName = strings.Map(
		func(r rune) rune {
			if unicode.IsSpace(r) {
				return '_'
			}
			return r
		},
		strings.TrimSuffix(filepath.Base(listTextFile), ".txt"),
	)

	fmt.Printf("::group::List: %s\n", listName)
	defer fmt.Println("::endgroup::")

	var (
		outputFile = filepath.Join(distDir, listName+".dat")
		logFile    = filepath.Join(logDir, listName+".log")
	)

	// Load all URLs
	filterListURLs, err := util.ReadListFile(listTextFile)
	if err != nil {
		return fmt.Errorf("reading list file: %w", err)
	}

	// Create temporary directory and make sure we remove it afterwards
	err = os.MkdirAll(tmpDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("creating temp directory for filter lists: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	log.Printf("Downloading %d filter lists...\n", len(filterListURLs))

	// Actually download lists
	paths, err := util.DownloadURLs(filterListURLs, tmpDir)
	if err != nil {
		return fmt.Errorf("downloading filter lists: %w", err)
	}

	log.Printf("Got %d/%d\n", len(paths), len(filterListURLs))

	if len(paths) == 0 {
		return fmt.Errorf("all lists failed to download")
	}

	// make sure dir exists
	err = os.MkdirAll(filepath.Dir(outputFile), 0664)
	if err != nil {
		return fmt.Errorf("creating distribution directory: %w", err)
	}

	log.Println("Converting ruleset...")
	err = util.GenerateDistributableList(paths, outputFile, logFile)
	if err != nil {
		return fmt.Errorf("generating distributable list: %w", err)
	}

	// Check if output file is larger than 10mb
	fileInfo, err := os.Stat(outputFile)
	if err != nil {
		return fmt.Errorf("getting filter output file info: %w", err)
	}
	if fileInfo.Size() > bromiteMaxFilterSize {
		return fmt.Errorf("filter list is too large for Bromite (%d bytes > %d bytes)", fileInfo.Size(), bromiteMaxFilterSize)
	}

	err = util.AppendReleaseList(listTextFile, len(paths), len(filterListURLs))
	if err != nil {
		return fmt.Errorf("generating release list: %w", err)
	}

	return nil
}

func main() {
	err := filepath.Walk(listDir, func(path string, d os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Walk: Error: %s\n", err.Error())
			return nil
		}
		if d.IsDir() {
			return nil
		}

		if !strings.HasSuffix(path, ".txt") {
			log.Printf("File %q was not processed because it does not end with \".txt\"\n", path)
			return nil
		}

		err = generateFilterList(path)
		if err != nil {
			log.Printf("Error while generating filter for %q: %s\n", path, err.Error())
		}
		return err
	})
	if err != nil {
		panic("error while walking: " + err.Error())
	}
}
