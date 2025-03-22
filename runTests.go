package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func runTests(dirName string, repoNames []byte) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	outputDir := fmt.Sprintf("--output-dir=%s/%s/%s", cwd, testsDir, dirName)
	for _, repo := range strings.Split(string(repoNames), "\n") {
		replacedDirName := strings.ReplaceAll(repo, "/", "_")
		cmd := exec.Command(
			"ginkgo",
			"--tags=e2e,osde2e",
			"--flake-attempts=3",
			"--procs=4",
			"-vv",
			"--trace",
			outputDir,
			fmt.Sprintf("--json-report=%s.json", replacedDirName),
			".",
		)
		cmd.Dir = fmt.Sprintf("%s/%s/test/e2e", reposDir, replacedDirName)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		_ = cmd.Run()
		//		outputBytes, err := cmd.CombinedOutput()
		//		ifErr(err)
		//		fmt.Println("outputBytes", outputBytes)
		//		ifErr(cmd.Process.Signal(syscall.SIGTERM))
		//		return string(outputBytes), err
	}
	return nil
}
