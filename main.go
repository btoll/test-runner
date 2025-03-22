package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/google/go-github/v68/github"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	client     *github.Client
	ctx        context.Context = context.Background()
	attempts   int
	generate   bool
	org        string = "openshift"
	reposDir   string = "repos"
	reportsDir string = "reports"
	testsDir   string = "tests"
)

func clone(r *github.Repository, cloneTo string) error {
	fmt.Printf("Cloning repository: %s\n", *r.FullName)
	cmd := exec.Command("git", "clone", "--depth=1", *r.CloneURL, fmt.Sprintf("%s/%s", reposDir, cloneTo))
	//	cmd.Stdout = os.Stdout
	//	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func cloneRunGenerate() error {
	err := os.MkdirAll(reposDir, os.ModePerm)
	if err != nil {
		return err
	}
	if !terminal.IsTerminal(syscall.Stdin) {
		b, err := io.ReadAll(os.Stdin)
		if err != nil {
			return err
		}
		for _, repo := range strings.Split(string(b), "\n") {
			if repo != "" {
				s := strings.Split(repo, "/")
				b, _ := exists(fmt.Sprintf("repos/%s_%s", s[0], s[1]))
				if b {
					continue
				}
				r, _, err := client.Repositories.Get(ctx, s[0], s[1])
				if err != nil {
					return err
				}
				err = clone(r, strings.ReplaceAll(repo, "/", "_"))
				if err != nil {
					// if exitError, ok := err.(*exec.ExitError); ok {
					// A 128 in this context means that the repo has already been cloned and the directory
					// exists.  For example:
					// fatal: destination path 'repos/openshift_must-gather-operator' already exists and is not an empty directory.
					//     if exitError.ExitCode() == 128 {
					//         continue
					//     }
					// }
					return err
				}
			}
		}
		err = runTestsAndGenerateReport(attempts, b)
		if err != nil {
			return err
		}
	} else {
		options := &github.RepositoryListByOrgOptions{
			ListOptions: github.ListOptions{PerPage: 100},
			Sort:        "full_name",
		}
		var wg sync.WaitGroup
		for {
			repos, resp, err := client.Repositories.ListByOrg(ctx, org, options)
			if err != nil {
				return err
			}
			for _, r := range repos {
				if strings.Contains(strings.ToLower(*r.Name), "operator") {
					wg.Add(1)
					go func(r *github.Repository) error {
						err := clone(r, *r.Name)
						if err != nil {
							return err
						}
						wg.Done()
						return nil
					}(r)
				}
			}
			if resp.NextPage == 0 {
				break
			}
			options.Page = resp.NextPage
		}
		wg.Wait()
	}
	return nil
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}
	return false, err
}

func getDirName() string {
	return time.Now().Format(time.RFC3339)
}

func ifErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func runTestsAndGenerateReport(attempts int, b []byte) error {
	for _ = range attempts {
		dirName := getDirName()
		err := runTests(dirName, b)
		if err != nil {
			return err
		}
		err = generateReport(dirName)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	flag.BoolVar(&generate, "generateReport", false, "Generate a report from the results of running Ginkgo.")
	flag.IntVar(&attempts, "attempts", 1, "Number of times to run the Ginkgo tests.")
	flag.Parse()

	if generate {
		ifErr(generateReport("2025-04-07T23:18:51-04:00"))
	} else {
		apiToken, isSet := os.LookupEnv("GITHUB_TOKEN")
		if apiToken == "" || !isSet {
			panic("[ERROR] Must set $GITHUB_TOKEN!")
		}
		client = github.NewClient(nil).WithAuthToken(apiToken)
		ifErr(cloneRunGenerate())
	}
}
