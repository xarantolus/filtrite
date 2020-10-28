package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"xarantolus/filtrite/util"
)

const (
	tmpDir = "tmp"
)

func main() {
	flag.Parse()

	listTextFile, outputFile := flag.Arg(0), flag.Arg(1)
	if listTextFile == "" || outputFile == "" {
		log.Fatalln("need two args, first is URL input list and second one is output file")
	}

	log.Println("Reading lists...")

	// Load all URLs
	filterListURLs, err := util.ReadListFile(listTextFile)
	if err != nil {
		panic("reading list file: " + err.Error())
	}

	// Create temporary directory and make sure we remove it afterwards
	err = os.MkdirAll(tmpDir, 0644)
	if err != nil {
		panic("creating temp directory for filter lists:" + err.Error())
	}
	defer os.RemoveAll(tmpDir)

	log.Printf("Downloading %d filter lists...\n", len(filterListURLs))

	// Actually download lists
	paths, err := util.DownloadURLs(filterListURLs, tmpDir)
	if err != nil {
		panic("downloading filter lists: " + err.Error())
	}

	log.Printf("Got %d/%d\n", len(paths), len(filterListURLs))

	// make sure dir exists
	err = os.MkdirAll(filepath.Dir(outputFile), 0664)
	if err != nil {
		panic("creating distribution directory: " + err.Error())
	}

	log.Println("Converting ruleset...")
	err = util.GenerateDistributableList(paths, outputFile)
	if err != nil {
		panic("generating distributable list: " + err.Error())
	}
}
