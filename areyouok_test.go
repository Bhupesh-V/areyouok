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
	goodUrls := []string{
		"https://good.com",
		"http://work.io",
		"http://www.website.gov.uk",
		"http://www.website.gov.uk/index.html",
		"http://website.in/843783787",
	}
	for _, url := range goodUrls {
		ans := re.MatchString(url)
		if ans != true {
			t.Errorf("RegEx %s want %s got %s", url, "true", "false")
		}
	}
}

func TestRegExBad(t *testing.T) {
	badUrls := []string{
		"emailto:test@gam.com",
		"ht#tp://www.website.gov.uk",
		"example.com/file[/].html",
	}
	for _, url := range badUrls {
		ans := re.MatchString(url)
		if ans != false {
			t.Errorf("RegEx %s want %s got %s", url, "false", "true")
		}
	}
}

func TestGit(t *testing.T) {
	getGitDetails(".")
	correctRemote := "https://github.com/Bhupesh-V/areyouok"
	correctBranch := "master"
	if branchName != correctBranch {
		t.Errorf("GitDetails want %s got %s", correctBranch, branchName)
	}
	if repoURL != correctRemote {
		t.Errorf("GitDetails want %s got %s", correctRemote, repoURL)
	}
}
