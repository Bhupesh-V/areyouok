package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

/*
TODO:
1. Concurrent reading of files?
2. Handle localhost URL
3. Handle ` char in url
4. Improve --help
*/

// its not perfect (look for edge cases)
// https://www.suon.co.uk/product/1/7/3/
var re = regexp.MustCompile(`https?:\/\/[^)\n,\s\"*]+`)

var (
	totalTime  string
	totalFiles int
	totalLinks int
)

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
		// fmt.Printf("Skipping %s due to Error: %s\n", link, err)
		ch <- map[string]string{"url": link, "message": err.Error()}
		return
	}
	ch <- map[string]string{
        "url": link, 
        "code": strconv.Itoa(resp.StatusCode), 
        "message": http.StatusText(resp.StatusCode),
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

func GetLinks(files []string) []string {
	hyperlinks := make(map[string][]string)
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
	for _, v := range hyperlinks {
		allLinks = append(allLinks, v...)
	}
	// yay! Jackpot!!
	totalFiles = len(hyperlinks)
	totalLinks = len(allLinks)
	fmt.Printf("%d links found across %d files\n\n", len(allLinks), len(hyperlinks))

	return allLinks
}

func GenerateReport(data []map[string]string, reportType string) {
    currentDir, err := os.Getwd()
    if err != nil {
        fmt.Println(err)
    }
	if reportType == "html" {
		now := time.Now()
		t, err := template.ParseFiles("static/report_template.html")
		if err != nil {
			fmt.Println(err)
		}
		f, err := os.Create("report.html")
		if err != nil {
			log.Println("File error: ", err)
			return
		}
		templateData := struct {
			NotOkurls  []map[string]string
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
        fmt.Printf("\nReport Generated at %s", filepath.Join(currentDir, "report.html"))
	} else if reportType == "json" {
		j, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			fmt.Println(err)
		}
		err = ioutil.WriteFile("report."+reportType, j, 0644)
        fmt.Printf("\nReport Generated at %s", filepath.Join(currentDir, "report.json"))
	}
}

func Driver(links []string) []map[string]string {
	var wg sync.WaitGroup
	var notoklinks []map[string]string
	start := time.Now()
	ch := make(chan map[string]string, len(links)) //unbuffered channel
	wg.Add(len(links))
	for _, url := range links {
		go checkLink(url, &wg, ch)
	}
	go func() {
		wg.Wait()
		close(ch)
	}()
	for i := range links {
		notoklinks = append(notoklinks, <-ch)
		fmt.Printf("\rAnalyzing %d/%d URLs", i+1, len(links))
	}
	// fmt.Println(Indent(notoklinks))
	totalTime = fmt.Sprintf("%.2fs", time.Since(start).Seconds())
	fmt.Printf("\nTotal Time: %.2fs\n", time.Since(start).Seconds())
	return notoklinks
}

func main() {
	var (
		typeOfFile string
		ignoreDirs string
		userDir   string
		reportType string
		dirs       []string
	)
	flag.StringVar(&typeOfFile, "t", "md", "Specify type of files to scan")
	flag.StringVar(&ignoreDirs, "i", "", "Comma separated directory and/or file names to ignore")
	flag.StringVar(&reportType, "r", "html", "Generate report. Supported formats include json/html")

	flag.Usage = func() {
		fmt.Fprintf(os.Stdout, "AreYouOK URL Health Checker\n")
		fmt.Fprintf(os.Stdout, "Usage: areyouok [OPTIONS] <directory-path>\nFollowing options are available:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stdout, "\nExample: areyouok -i=.git,README.md -r=html Documents/some-dir/\n\n")
		fmt.Fprintf(os.Stdout, "Report Any Bugs to varshneybhupesh@gmail.com\n")
	}
	flag.Parse()
	if ignoreDirs != "" {
		dirs = strings.Split(ignoreDirs, ",")
	}

	if len(flag.Args()) == 0 {
		userDir = "."
	} else {
		userDir = flag.Args()[0]
	}

	var validFiles = GetFiles(userDir, typeOfFile, dirs)
	data := Driver(GetLinks(validFiles))
	GenerateReport(data, reportType)
}
