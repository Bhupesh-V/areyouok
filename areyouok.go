package main

import (
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	textTemplate "text/template"
	"time"
)

// its not perfect (look for edge cases)
var re = regexp.MustCompile(`(http|ftp|https):\/\/([\w\-_]+(?:(?:\.[\w\-_]+)+))([\w\-\.,@?^=%&amp;:/~\+#]*[\w\-\@?^=%&amp;/~\+#])?`)

var (
	totalTime  string
	totalFiles int
	totalLinks int
	//Do not modify, its done at compile time using ldflags
	aroVersion string = "dev" //aro Version
	aroDate    string = "dev" //aro Build Date
	branchName string
	repoURL    string
)

func checkLink(link string, wg *sync.WaitGroup, ch chan map[string]string) {
	defer wg.Done()
	goStart := time.Now()
	reqURL, _ := url.Parse(link)
	req, err := http.NewRequest("GET", reqURL.String(), nil)
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/77.0.3865.90 Safari/537.36")
	resp, err := http.DefaultClient.Do(req)
	responseTime := fmt.Sprintf("%.2fs", time.Now().Sub(goStart).Seconds())
	if err != nil {
		ch <- map[string]string{"url": link, "message": err.Error()}
		return
	}
	ch <- map[string]string{
		"url":           link,
		"code":          strconv.Itoa(resp.StatusCode),
		"message":       http.StatusText(resp.StatusCode),
		"response_time": responseTime,
	}
}

func in(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func getGitDetails(userDir string) {
	// get default branch name
	cmd, err := exec.Command("git", "-C", userDir, "symbolic-ref", "--short", "HEAD").CombinedOutput()
	if err == nil {
		branchName = strings.Trim(string(cmd[:]), "\r\n")
	}
	// get repo url
	config, err := exec.Command("git", "-C", userDir, "config", "--get", "remote.origin.url").CombinedOutput()
	if err == nil {
		repoURL = string(config[:])
		repoURL = repoURL[0 : len(repoURL)-5]
	}
}

func getFiles(userPath string, filetype string, ignore []string) []string {
	var validFiles []string

	err := filepath.Walk(userPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				if in(info.Name(), ignore) {
					return filepath.SkipDir
				}
			}
			if strings.HasSuffix(filepath.Base(path), filetype) && !in(info.Name(), ignore) {
				validFiles = append(validFiles, path)
			}
			return nil
		})
	if err != nil {
		log.Println(err)
	}
	return validFiles
}

func getLinks(files []string) ([]map[string]string, map[string][]string) {
	hyperlinks := make(map[string][]string)
	var allLinks []string
	var allHyperlinks []map[string]string

	for _, file := range files {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			fmt.Println("File reading error", err)
		}
		fileContent := string(data)
		submatchall := re.FindAllString(fileContent, -1)
		if len(submatchall) > 0 {
			hyperlinks[file] = submatchall
		}
	}
	for filepath, v := range hyperlinks {
		allLinks = append(allLinks, v...)
		for _, link := range v {
			allHyperlinks = append(allHyperlinks, map[string]string{"file": filepath, "url": link})
		}
	}
	// yay! Jackpot!!
	totalFiles = len(hyperlinks)
	totalLinks = len(allLinks)
	fmt.Printf("Found %d link(s) across %d file(s)\n\n", totalLinks, totalFiles)
	for file := range hyperlinks {
		fmt.Println(file)
	}
	fmt.Println()

	return allHyperlinks, hyperlinks
}

func generateReport(validfiles map[string][]string, linkfr map[string]map[string]string, reportType string) {
	currentDir, _ := os.Getwd()
	//go:embed static/*
	var reportTemplates embed.FS

	if reportType == "html" || reportType == "github" {
		now := time.Now()
		t, err := template.ParseFS(reportTemplates, fmt.Sprintf("static/report_%s.html", reportType))
		if err != nil {
			fmt.Println(err)
		}
		f, err := os.Create(fmt.Sprintf("report.%s", reportType))
		if err != nil {
			fmt.Errorf("open report.%s failed: %w", reportType, err)
		}
		templateData := struct {
			ValidFiles map[string][]string
			ReLinks    map[string]map[string]string
			Date       string
			Time       string
			TotalLinks string
			TotalFiles string
			TotalTime  string
			BranchName string
			RepoURL    string
		}{
			ValidFiles: validfiles,
			ReLinks:    linkfr,
			Date:       now.Format("January 2, 2006"),
			Time:       now.Format(time.Kitchen),
			TotalLinks: strconv.Itoa(totalLinks),
			TotalFiles: strconv.Itoa(totalFiles),
			TotalTime:  totalTime,
			BranchName: branchName,
			RepoURL:    repoURL,
		}
		t.Execute(f, templateData)
	} else if reportType == "json" {
		j, err := json.MarshalIndent(linkfr, "", "  ")
		if err != nil {
			fmt.Println(err)
		}
		err = ioutil.WriteFile("report."+reportType, j, 0644)
	} else if reportType == "txt" {
		t := textTemplate.Must(textTemplate.ParseFS(reportTemplates, fmt.Sprintf("static/report_%s.txt", reportType)))
		templateData := struct {
			ReLinks    map[string]map[string]string
			TotalLinks string
			TotalFiles string
			TotalTime  string
		}{
			ReLinks:    linkfr,
			TotalLinks: strconv.Itoa(totalLinks),
			TotalFiles: strconv.Itoa(totalFiles),
			TotalTime:  totalTime,
		}
		f, _ := os.Create("report.txt")
		t.Execute(f, templateData)
	}
	fmt.Printf("\nReport Generated: %s.%s", filepath.Join(currentDir, "report"), reportType)
}

func driver(links []map[string]string) []map[string]string {
	var wg sync.WaitGroup
	var notoklinks []map[string]string
	start := time.Now()
	ch := make(chan map[string]string, len(links)) //unbuffered channel
	wg.Add(len(links))
	for _, v := range links {
		go checkLink(v["url"], &wg, ch)
	}
	go func() {
		wg.Wait()
		close(ch)
	}()
	for i := range links {
		notoklinks = append(notoklinks, <-ch)
		fmt.Printf("\rAnalyzing %d/%d URLs", i+1, len(links))
	}
	// fmt.Println(indent(notoklinks))
	totalTime = fmt.Sprintf("%.2fs", time.Since(start).Seconds())
	fmt.Printf("\nTotal Time: %.2fs\n", time.Since(start).Seconds())
	return notoklinks
}

func main() {
	var (
		typeOfFile string
		ignoreDirs string
		userDir    string
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

	if len(flag.Args()) == 0 {
		userDir = "."
	} else {
		userDir = flag.Args()[0]
	}
	if !in(reportType, []string{"github", "json", "txt", "html", ""}) {
		fmt.Printf("%s in not a supported report format\n", reportType)
		os.Exit(1)
	}
	getGitDetails(userDir)
	validFiles := getFiles(userDir, typeOfFile, dirs)
	links, valid := getLinks(validFiles)
	data := driver(links)
	linkfr := make(map[string]map[string]string)
	for _, v := range data {
		urlMap := map[string]string{
			"code":          v["code"],
			"message":       v["message"],
			"response_time": v["response_time"],
		}
		linkfr[v["url"]] = urlMap
	}
	if reportType != "" {
		generateReport(valid, linkfr, reportType)
	}
}
