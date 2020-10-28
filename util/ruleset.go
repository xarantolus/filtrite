package util

import (
	"os"
	"os/exec"
	"strings"
)

const (
	rulesetExecutable = "deps/ruleset_converter"
)

func GenerateDistributableList(inputPaths []string, output string) (err error) {
	// ruleset_converter --input_format=filter-list \
	// --output_format=unindexed-ruleset \
	//         --input_files=easyprivacy.txt,easylist.txt \
	// --output_file=filters.dat
	// Example: https://www.bromite.org/custom-filters

	cmd := exec.Command(rulesetExecutable, "--input_format=filter-list", "--output_format=unindexed-ruleset", "--input_files="+strings.Join(inputPaths, ","), "--output_file="+output)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
