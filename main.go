package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/Bhupesh-V/areyouok/links"
	"github.com/Bhupesh-V/areyouok/report"
	"github.com/Bhupesh-V/areyouok/utils"
)

var (
	//Do not modify, its done at compile time using ldflags
	aroVersion string = "dev" //aro Version
	aroDate    string = "dev" //aro Build Date
)

func driver(validLinks []map[string]string) ([]map[string]string, string) {
	var wg sync.WaitGroup
	var notoklinks []map[string]string
	start := time.Now()
	ch := make(chan map[string]string, len(validLinks)) //unbuffered channel
	wg.Add(len(validLinks))
	for _, v := range validLinks {
		go links.CheckHealth(v["url"], &wg, ch)
	}
	go func() {
		wg.Wait()
		close(ch)
	}()
	for i := range validLinks {
		notoklinks = append(notoklinks, <-ch)
		fmt.Printf("\rAnalyzing %d/%d URLs", i+1, len(validLinks))
	}
	totalTime := fmt.Sprintf("%.2fs", time.Since(start).Seconds())
	fmt.Printf("\nTotal Time: %.2fs\n", time.Since(start).Seconds())

	return notoklinks, totalTime
}

func displayLinks(valid *map[string][]string) {
	userDir := utils.GetUserDirectory()
	// fmt.Printf("Found %s URL(s) across %s file(s)\n\n", strconv.Itoa(totalLinks), strconv.Itoa(totalFiles))
	for file, urls := range *valid {
		rel, _ := filepath.Rel(userDir, file)
		fmt.Printf("%d 🔗️ %s\r\n", len(urls), rel)
	}
	fmt.Println()
}

func main() {
	var (
		typeOfFile string
		ignoreDirs string
		reportType string
		dirs       []string
	)
	flag.StringVar(&typeOfFile, "t", "md", "Specify `type` of files to scan")
	flag.StringVar(&ignoreDirs, "i", "", "Comma separated directory and/or file names to `ignore`")
	flag.StringVar(&reportType, "r", "", "Generate `report`. Supported formats include json, html, txt & github")
	Version := flag.Bool("v", false, "Prints current AreYouOk `version`")

	flag.Usage = func() {
		fmt.Fprintf(os.Stdout, "AreYouOK URL Health Checker\n")
		fmt.Fprintf(os.Stdout, "Usage: areyouok [OPTIONS] <directory-path>\nFollowing options are available:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stdout, "\nExample: areyouok -t=html -i=.git,README.md -r=json Documents/some-dir/\n\n")
		fmt.Fprintf(os.Stdout, "Report Any Bugs via \nEmail  : varshneybhupesh@gmail.com\nGitHub : https://github.com/Bhupesh-V/areyouok/issues/new/choose\n")
	}
	flag.Parse()
	if *Version {
		fmt.Printf("AreYouOk %s built on %s", aroVersion, aroDate)
		os.Exit(0)
	}
	if ignoreDirs != "" {
		dirs = strings.Split(ignoreDirs, ",")
	}

	if !utils.In(reportType, []string{"github", "json", "txt", "html", ""}) {
		fmt.Printf("%s in not a supported report format\n", reportType)
		os.Exit(1)
	}
	validFiles := links.GetFiles(typeOfFile, dirs)
	links := links.GetLinks(validFiles)

	displayLinks(&links.FileToListOfLinks)

	data, totalTime := driver(links.AllHyperlinks)
	healthData := make(map[string]interface{})

	for _, v := range data {
		urlMap := map[string]string{
			"code":          v["code"],
			"message":       v["message"],
			"response_time": v["response_time"],
		}
		healthData[v["url"]] = urlMap
	}
	if reportType != "" {
		rData := report.Report{
			ReportType: &reportType,
			ReportData: &report.ReportData{
				TotalTime:          &totalTime,
				TotalLinks:         &links.TotalLinks,
				TotalFiles:         &links.TotalValidFiles,
				ValidFiles:         links.FileToListOfLinks,
				NotOkLinks:         data,
				CompleteHealthData: healthData,
			},
		}
		rData.GenerateReport()
		// generateReport(links, healthData, reportType)
	}
}
