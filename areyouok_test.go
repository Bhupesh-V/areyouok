package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/Bhupesh-V/areyouok/git"
	"github.com/Bhupesh-V/areyouok/links"
	"github.com/Bhupesh-V/areyouok/utils"
)

func Test_In(t *testing.T) {
	t.Run("in", func(t *testing.T) {
		ans := utils.In("test", []string{"test", "sample"})
		if ans != true {
			t.Errorf("In() want %s, want %s", "false", "true")
		}
	})
	t.Run("not in", func(t *testing.T) {
		ans := utils.In("nice", []string{"test", "sample"})
		if ans != false {
			t.Errorf("In() want %s, want %s", "true", "false")
		}
	})
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
		ans := links.GetURLRegex().MatchString(url)
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
		ans := links.GetURLRegex().MatchString(url)
		if ans != false {
			t.Errorf("RegEx %s want %s got %s", url, "false", "true")
		}
	}
}

// Compare 2 slices
func Equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func AreEqualJSON(s1, s2 string) (bool, error) {
	var o1 interface{}
	var o2 interface{}

	var err error
	err = json.Unmarshal([]byte(s1), &o1)
	if err != nil {
		return false, fmt.Errorf("Error mashalling string 1 :: %s", err.Error())
	}
	err = json.Unmarshal([]byte(s2), &o2)
	if err != nil {
		return false, fmt.Errorf("Error mashalling string 2 :: %s", err.Error())
	}

	return reflect.DeepEqual(o1, o2), nil
}

func Test_getGitDetails(t *testing.T) {
	userPath := "."

	correctRemote := "https://github.com/Bhupesh-V/areyouok"
	correctBranch := "master"

	ansBranch, _ := git.GetGitBranch(&userPath)
	ansRemote, _ := git.GetGitRemoteURL(&userPath)

	if correctBranch != ansBranch {
		t.Errorf("GitDetails want %s got %s", correctBranch, ansBranch)
	}
	if correctRemote != ansRemote {
		t.Errorf("GitDetails want %s got %s", correctRemote, ansRemote)
	}
}

func Test_getValidFiles(t *testing.T) {

	t.Run("With no ignore", func(t *testing.T) {
		valid := []string{
			".github/ISSUE_TEMPLATE/----bug-report.md",
			".github/ISSUE_TEMPLATE/---feature-request.md",
			".github/ISSUE_TEMPLATE/---question.md",
			".github/ISSUE_TEMPLATE/---say-thank-you.md",
			".github/PULL_REQUEST_TEMPLATE/pull_request_template.md",
			"CHANGELOG.md",
			"CODE_OF_CONDUCT.md",
			"CONTRIBUTING.md",
			"README.md",
		}
		ans := links.GetFiles("md", []string{""})
		if !Equal(valid, ans) {
			t.Errorf("GetValidFiles() want %s got %s", valid, ans)
		}
	})
	t.Run("With multiple ignore", func(t *testing.T) {
		valid := []string{
			"CODE_OF_CONDUCT.md",
			"CONTRIBUTING.md",
			"README.md",
		}
		ans := links.GetFiles("md", []string{".github", "CHANGELOG.md"})
		if !Equal(valid, ans) {
			t.Errorf("GetValidFiles() want %s got %s", valid, ans)
		}
	})
}

func Test_getLinks(t *testing.T) {
	links := links.GetLinks([]string{"CODE_OF_CONDUCT.md"})
	a1 := links.AllHyperlinks
	a2 := links.FileToListOfLinks

	t.Run("check file path based JSON", func(t *testing.T) {
		valid := []map[string]string{
			{
				"file": "CODE_OF_CONDUCT.md",
				"url":  "https://www.contributor-covenant.org/version/1/4/code-of-conduct.html",
			},
			{
				"file": "CODE_OF_CONDUCT.md",
				"url":  "https://www.contributor-covenant.org",
			},
			{
				"file": "CODE_OF_CONDUCT.md",
				"url":  "https://www.contributor-covenant.org/faq",
			},
		}
		j, _ := json.MarshalIndent(a1, "", "  ")
		v, _ := json.MarshalIndent(valid, "", "  ")
		result, err := AreEqualJSON(string(v), string(j))
		if !result || err != nil {
			t.Errorf("getLinks() want %s got %s", valid, a2)
		}
	})
	t.Run("check list of JSON", func(t *testing.T) {
		valid := map[string][]string{
			"CODE_OF_CONDUCT.md": {
				"https://www.contributor-covenant.org/version/1/4/code-of-conduct.html",
				"https://www.contributor-covenant.org",
				"https://www.contributor-covenant.org/faq",
			},
		}
		j2, _ := json.MarshalIndent(a2, "", "  ")
		v, _ := json.MarshalIndent(valid, "", "  ")
		result, err := AreEqualJSON(string(v), string(j2))
		if !result || err != nil {
			t.Errorf("getLinks() want %s got %s", valid, a2)
		}
	})
}
