package main

import (
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
    // "reflect"
	"log"
	"net/http"
	"net/url"
	"os"
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
	aroVersion string = "dev"
)

type HyperLink struct {
	Link               string
	File               string
	Response_time      string
	Http_response_code string
	Http_message       string
}

func checkLink(link string, wg *sync.WaitGroup, ch chan map[string]string) {
	defer wg.Done()
	goStart := time.Now()
	reqURL, _ := url.Parse(link)
	req := &http.Request{
		Method: "GET",
		URL:    reqURL,
		Header: map[string][]string{
			"User-Agent": {"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/77.0.3865.90 Safari/537.36"},
		},
	}
	resp, err := http.DefaultClient.Do(req)
	responseTime := fmt.Sprintf("%.2fs", time.Since(goStart).Seconds())
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

func In(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func GetFiles(userPath string, filetype string, ignore []string) []string {
	var validFiles []string

	err := filepath.Walk(userPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				if In(info.Name(), ignore) {
					return filepath.SkipDir
				}
			}
			if strings.HasSuffix(filepath.Base(path), filetype) && !In(info.Name(), ignore) {
				validFiles = append(validFiles, path)
			}
			return nil
		})
	if err != nil {
		log.Println(err)
	}
	return validFiles
}

func Indent(v interface{}) string {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Sprintf("%#v", v)
	}
	return string(b)
}

func GetLinks(files []string) []*HyperLink {
	hyperlinks := make(map[string][]string)
    var all_hyperlinks []*HyperLink
	var allLinks []string

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
    // fmt.Println(Indent(hyperlinks))
    fmt.Println(len(hyperlinks))
	for filepath, v := range hyperlinks {
		allLinks = append(allLinks, v...)
        for _, link := range v {
            // fmt.Println(r)
            all_hyperlinks = append(all_hyperlinks, &HyperLink{Link: link, File: filepath})
            // fmt.Println(all_hyperlinks[link].File)
            // fmt.Println(len(all_hyperlinks))
        }
	}
    // fmt.Println(reflect.ValueOf(all_hyperlinks).MapKeys())
    fmt.Println(len(all_hyperlinks))
    fmt.Println(all_hyperlinks[3].File)
    fmt.Println(all_hyperlinks[3].Link)
    // fmt.Println(Indent(allLinks))
	// yay! Jackpot!!
	totalFiles = len(hyperlinks)
	totalLinks = len(allLinks)
	fmt.Printf("%d links found across %d files\n\n", len(allLinks), len(hyperlinks))
	// fmt.Println(Indent(hyperlinks))

	return all_hyperlinks
}

func GenerateReport(data []*HyperLink, reportType string) {
	currentDir, _ := os.Getwd()
	//go:embed static/*
	var report_templates embed.FS

	if reportType == "html" || reportType == "github" {
		now := time.Now()
		t, err := template.ParseFS(report_templates, fmt.Sprintf("static/report_%s.html", reportType))
		if err != nil {
			fmt.Println(err)
		}
		f, _ := os.Create("report.html")
		templateData := struct {
			NotOkurls  []*HyperLink
			Date       string
			Time       string
			TotalLinks string
			TotalFiles string
			TotalTime  string
		}{
			NotOkurls:  data,
			Date:       now.Format("January 2, 2006"),
			Time:       now.Format(time.Kitchen),
			TotalLinks: strconv.Itoa(totalLinks),
			TotalFiles: strconv.Itoa(totalFiles),
			TotalTime:  totalTime,
		}
		t.Execute(f, templateData)
	} else if reportType == "json" {
		j, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			fmt.Println(err)
		}
		err = ioutil.WriteFile("report."+reportType, j, 0644)
	} else if reportType == "txt" {
		t := textTemplate.Must(textTemplate.New("t1").
			Parse(`{{.TotalLinks}} URLs were analyzed across {{.TotalFiles}} files in {{ println .TotalTime}}{{"\n"}}Following URLs were found not OK:{{"\n\n"}}{{range $_, $v := $.NotOkurls}}{{ if ne $v.message "OK" }}{{ println $v.url }}{{end}}{{end}}`))
		templateData := struct {
			NotOkurls  []*HyperLink
			TotalLinks string
			TotalFiles string
			TotalTime  string
		}{
			NotOkurls:  data,
			TotalLinks: strconv.Itoa(totalLinks),
			TotalFiles: strconv.Itoa(totalFiles),
			TotalTime:  totalTime,
		}
		f, _ := os.Create("report.txt")
		t.Execute(f, templateData)
	}
	fmt.Printf("\nReport Generated: %s.%s", filepath.Join(currentDir, "report"), reportType)
}

func Driver(data []*HyperLink) []*HyperLink {
	var wg sync.WaitGroup
	var analyzed_urls map[string]string
	start := time.Now()
	ch := make(chan map[string]string, len(data)) //unbuffered channel
	wg.Add(len(data))
	for _, d := range data {
		go checkLink(d.Link, &wg, ch)
	}
	go func() {
		wg.Wait()
		close(ch)
	}()
    count := 0
	for _, v := range data {
		analyzed_urls = <-ch
        // fmt.Println(analyzed_urls["code"])
        v.Http_response_code = analyzed_urls["code"]
        v.Http_message = analyzed_urls["message"]
        v.Response_time = analyzed_urls["response_time"]
        count += 1
        fmt.Printf("\rAnalyzing %d/%d URLs", count, totalLinks)
	}
	// fmt.Println(Indent(analyzed_urls))
	totalTime = fmt.Sprintf("%.2fs", time.Since(start).Seconds())
	fmt.Printf("\nTotal Time: %.2fs\n", time.Since(start).Seconds())
    return data
}

func main() {
	var (
		typeOfFile string
		ignoreDirs string
		userDir    string
		reportType string
		dirs       []string
	)
	flag.StringVar(&typeOfFile, "t", "md", "Specify type of files to scan")
	flag.StringVar(&ignoreDirs, "i", "", "Comma separated directory and/or file names to ignore")
	flag.StringVar(&reportType, "r", "html", "Generate report. Supported formats include json, html, txt & github")
	Version := flag.Bool("v", false, "Prints Current AreYouOk Version")

	flag.Usage = func() {
		fmt.Fprintf(os.Stdout, "AreYouOK URL Health Checker\n")
		fmt.Fprintf(os.Stdout, "Usage: areyouok [OPTIONS] <directory-path>\nFollowing options are available:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stdout, "\nExample: areyouok -i=.git,README.md -r=html Documents/some-dir/\n\n")
		fmt.Fprintf(os.Stdout, "Report Any Bugs to varshneybhupesh@gmail.com\n")
	}
	flag.Parse()
	if *Version {
		fmt.Println(aroVersion)
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
	if !In(reportType, []string{"github", "json", "txt", "html"}) {
		fmt.Printf("%s in not a supported report format\n", reportType)
		os.Exit(1)
	}
	var validFiles = GetFiles(userDir, typeOfFile, dirs)
    data := Driver(GetLinks(validFiles))
	GenerateReport(data, reportType)
}
