package links

import (
	"io/ioutil"
	"regexp"

	"github.com/Bhupesh-V/areyouok/utils"
)

type Links struct {
	// count of all valid links
	TotalLinks int
	// count of files with links
	TotalValidFiles int
	// Map of file to list of links
	FileToListOfLinks map[string][]string
	AllHyperlinks     []map[string]string
}

// Get all valid links
func GetLinks(files []string) *Links {
	var totalLinks int
	var linkRegex = regexp.MustCompile(`(http|https):\/\/([\w\-_]+(?:(?:\.[\w\-_]+)+))([\w\-\.,@?^=%&amp;:/~\+#]*[\w\-\@?^=%&amp;/~\+#])?`)
	// Map of file to list of links
	hyperlinks := make(map[string][]string)
	// List of JSON {file, url}
	var allHyperlinks []map[string]string

	for _, file := range files {
		data, err := ioutil.ReadFile(file)
		utils.CheckErr(err)
		fileContent := string(data)
		submatchall := linkRegex.FindAllString(fileContent, -1)
		if len(submatchall) > 0 {
			hyperlinks[file] = submatchall
		}
	}
	//TODO: figure this out
	for filepath, links := range hyperlinks {
		totalLinks = totalLinks + 1
		for _, link := range links {
			allHyperlinks = append(allHyperlinks, map[string]string{"file": filepath, "url": link})
		}
	}

	return &Links{
		TotalLinks:        totalLinks,
		TotalValidFiles:   len(hyperlinks),
		FileToListOfLinks: hyperlinks,
		AllHyperlinks:     allHyperlinks,
	}

	// return allHyperlinks, hyperlinks
}
