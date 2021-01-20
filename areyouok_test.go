package main

import (
	"testing"
)

func TestInA(t *testing.T) {
	ans := in("test", []string{"test", "sample"})
	if ans != true {
		t.Errorf("In() want %s, want %s", "false", "true")
	}
}

func TestInB(t *testing.T) {
	ans := in("nice", []string{"test", "sample"})
	if ans != false {
		t.Errorf("In() want %s, want %s", "true", "false")
	}
}

func TestRegExGood(t *testing.T) {
	good_urls := []string{
		"https://good.com",
		"http://work.io",
		"http://www.website.gov.uk",
		"http://www.website.gov.uk/index.html",
		"http://website.in/843783787",
	}
	for _, url := range good_urls {
		ans := re.MatchString(url)
		if ans != true {
			t.Errorf("RegEx %s want %s got %s", url, "true", "false")
		}
	}
}

func TestRegExBad(t *testing.T) {
	good_urls := []string{
		"emailto:test@gam.com",
		"ht#tp://www.website.gov.uk",
		"http://news.sky.com/skynews/article/0,,30200-1303092,00.html",
		"example.com/file[/].html",
        "http://domain.com/$dfd",
	}
	for _, url := range good_urls {
		ans := re.MatchString(url)
		if ans != false {
			t.Errorf("RegEx %s want %s got %s", url, "false", "true")
		}
	}
}

