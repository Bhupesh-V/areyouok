package links

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
)

func CheckHealth(link string, wg *sync.WaitGroup, ch chan map[string]string) {
	defer wg.Done()

	goStart := time.Now()
	reqURL, err := url.Parse(link)
	if err != nil {
		fmt.Println(err)
	}
	req, err := http.NewRequest("GET", reqURL.String(), nil)
	if err != nil {
		fmt.Println(err)
	}

	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/77.0.3865.90 Safari/537.36")

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
