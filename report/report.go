package report

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	textTemplate "text/template"
	"time"

	"github.com/Bhupesh-V/areyouok/git"
	"github.com/Bhupesh-V/areyouok/utils"
)

//go:embed formats/*
var reportTemplates embed.FS

type ReportMetaData struct {
	Date string
	Time string
}

type ReportData struct {
	// count of links who were checked
	TotalLinks int
	// count of files where all valid links were found
	TotalFiles int
	// total time it took to query all the links
	TotalTime *string
	// Current Git branch name
	BranchName *string
	// URL for remote code host github/gitlab etc
	RepoURL    *string
	ValidFiles map[string][]string
	// Links which have some problem, which didn't return status 200
	NotOkLinks         interface{}
	CompleteHealthData interface{}
	ReportMetaData     ReportMetaData
}

type Report struct {
	ReportType *string
	ReportData *ReportData
}

func (r *Report) GenerateReport() {
	now := time.Now()
	currentDir, err := os.Getwd()
	userDir := utils.GetUserDirectory()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if *r.ReportType == "html" || *r.ReportType == "github" {
		parsedTemplate, err := template.ParseFS(reportTemplates, fmt.Sprintf("formats/report_%s.html", *r.ReportType))
		utils.CheckErr(err)
		reportFile, err := os.Create(fmt.Sprintf("report.%s", *r.ReportType))
		utils.CheckErr(err)

		r.ReportData.ReportMetaData.Date = now.Format("January 2, 2006")
		r.ReportData.ReportMetaData.Time = now.Format(time.Kitchen)

		branch, err := git.GetGitBranch(&userDir)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		r.ReportData.BranchName = &branch
		remoteURL, err := git.GetGitRemoteURL(&userDir)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		r.ReportData.RepoURL = &remoteURL

		parsedTemplate.Execute(reportFile, r.ReportData)
	} else if *r.ReportType == "json" {
		jsonReport, err := json.MarshalIndent(r.ReportData.CompleteHealthData, "", "  ")
		utils.CheckErr(err)
		if err = ioutil.WriteFile("report."+*r.ReportType, jsonReport, 0644); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else if *r.ReportType == "txt" {
		textReport := textTemplate.Must(textTemplate.ParseFS(reportTemplates, fmt.Sprintf("formats/report_%s.txt", *r.ReportType)))
		f, _ := os.Create("report.txt")
		textReport.Execute(f, r.ReportData)
	}
	fmt.Printf("\nReport Generated: %s.%s", filepath.Join(currentDir, "report"), *r.ReportType)
}
