package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func runTests() {
	//	entries, err := os.ReadDir("repos")
	b, err := os.ReadFile("operators.list")
	ifErr(err)
	for _, repo := range strings.Split(string(b), "\n") {
		fmt.Println("repo", repo)
		cmd := exec.Command(
			"ginkgo",
			"--tags=e2e,osde2e",
			"--flake-attempts=3",
			"--procs=4",
			"-vv",
			"--trace",
			"--output-dir=/home/btoll/projects/redhat/test-runner/tests",
			fmt.Sprintf("--json-report=%s.json", repo),
			".",
		)
		fmt.Println(cmd.Path, cmd.Args)
		cmd.Dir = fmt.Sprintf("repos/%s/test/e2e", repo)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		_ = cmd.Run()
		//		outputBytes, err := cmd.CombinedOutput()
		//		ifErr(err)
		//		fmt.Println("outputBytes", outputBytes)
		//		ifErr(cmd.Process.Signal(syscall.SIGTERM))
		//		return string(outputBytes), err
	}
}
