package util

import (
	"os"
	"os/exec"
	"strings"
)

const (
	rulesetExecutable = "deps/ruleset_converter"
)

func GenerateDistributableList(inputPaths []string, output string, logPath string) (err error) {
	// ruleset_converter --input_format=filter-list \
	// --output_format=unindexed-ruleset \
	//         --input_files=easyprivacy.txt,easylist.txt \
	// --output_file=filters.dat
	// Example: https://www.bromite.org/custom-filters

	cmd := exec.Command(rulesetExecutable, "--input_format=filter-list", "--output_format=unindexed-ruleset", "--input_files="+strings.Join(inputPaths, ","), "--output_file="+output)

	if logPath != "" {
		f, err := os.Create(logPath)
		if err != nil {
			return err
		}
		defer f.Close()

		cmd.Stdout = f
		cmd.Stderr = f
	}

	return cmd.Run()
}
