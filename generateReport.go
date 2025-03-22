package main

import (
	"encoding/json"
	"fmt"
	"os"
	"text/template"
)

type suite struct {
	Path        string `json:"SuitePath"`
	Description string `json:"SuiteDescription"`
	Succeeded   bool   `json:"SuiteSucceeded"`
	PreRunStats struct {
		TotalSpecs       int `json:"TotalSpecs"`
		SpecsThatWillRun int `json:"SpecsThatWillRun"`
	} `json:"PreRunStats"`
	SpecReports []spec
}

type spec struct {
	Name     string `json:"LeafNodeText"`
	Type     string `json:"LeafNodeType"`
	State    string `json:"State"`
	Attempts int    `json:"NumAttempts"`
	Failure  struct {
		Message  string `json:"Message"`
		Location struct {
			LineNumber int    `json:"LineNumber"`
			StackTrace string `json:"FullStackTrace"`
		} `json:"Location"`
	}
}

func generateReport(dirName string) error {
	entries, err := os.ReadDir(fmt.Sprintf("%s/%s", testsDir, dirName))
	if err != nil {
		return err
	}

	// TODO: Check for `reportsDir`/`dirName` first?
	err = os.MkdirAll(fmt.Sprintf("%s/%s", reportsDir, dirName), os.ModePerm)
	if err != nil {
		return err
	}
	f, err := os.Create(fmt.Sprintf("%s/%s/report.txt", reportsDir, dirName))
	defer f.Close()
	if err != nil {
		return err
	}

	tplName := "suite.tpl"
	tpl := template.Must(template.New(tplName).ParseFiles(tplName))

	// TODO: Check for `.json` files only.
	for _, entry := range entries {
		var s []suite
		b, err := os.ReadFile(fmt.Sprintf("%s/%s/%s", testsDir, dirName, entry.Name()))
		if err != nil {
			fmt.Println(err)
			continue
		}
		err = json.Unmarshal(b, &s)
		if err != nil {
			fmt.Println(err)
		}
		err = tpl.Execute(f, s)
		if err != nil {
			return err
		}
	}
	return nil
}
