package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"

	"github.com/google/go-github/v68/github"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	client   *github.Client
	ctx      context.Context = context.Background()
	generate bool
	org      = "openshift"
	test     bool
)

func clone(r *github.Repository, cloneTo string) error {
	fmt.Printf("Cloning repository: %s\n", *r.FullName)
	cmd := exec.Command("git", "clone", "--depth=1", *r.CloneURL, fmt.Sprintf("repos/%s", cloneTo))
	return cmd.Run()
}

func ifErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	flag.BoolVar(&generate, "generateReport", false, "Generate a report from the results of running Ginkgo.")
	flag.BoolVar(&test, "runTests", false, "Run the Ginkgo tests.")
	flag.Parse()

	if generate {
		generateReport()
	} else if test {
		runTests()
	} else {
		apiToken, isSet := os.LookupEnv("GITHUB_TOKEN")
		if apiToken == "" || !isSet {
			panic("[ERROR] Must set $GITHUB_TOKEN!")
		}
		client = github.NewClient(nil).WithAuthToken(apiToken)

		ifErr(os.MkdirAll("repos", os.ModePerm))

		if !terminal.IsTerminal(syscall.Stdin) {
			b, err := io.ReadAll(os.Stdin)
			ifErr(err)
			// TODO: Add concurrency.
			for _, repo := range strings.Split(string(b), "\n") {
				if repo != "" {
					s := strings.Split(repo, "/")
					r, _, err := client.Repositories.Get(ctx, s[0], s[1])
					ifErr(err)
					ifErr(clone(r, strings.ReplaceAll(repo, "/", "_")))
				}
			}
		} else {
			options := &github.RepositoryListByOrgOptions{
				ListOptions: github.ListOptions{PerPage: 100},
				Sort:        "full_name",
			}
			var wg sync.WaitGroup
			for {
				repos, resp, err := client.Repositories.ListByOrg(ctx, org, options)
				ifErr(err)
				for _, r := range repos {
					if strings.Contains(strings.ToLower(*r.Name), "operator") {
						wg.Add(1)
						go func(r *github.Repository) {
							ifErr(clone(r, *r.Name))
							wg.Done()
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
	}
}
