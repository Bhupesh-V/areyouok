package git

import (
	"log"
	"os/exec"
	"strings"
)

func GetGitBranch(userDir *string) (string, error) {
	var branchName string
	cmd, err := exec.Command("git", "-C", *userDir, "symbolic-ref", "--short", "HEAD").CombinedOutput()
	if err != nil {
		log.Fatalln(err)
		return "", err
	}
	branchName = strings.Trim(string(cmd[:]), "\r\n")
	return branchName, nil
}

func GetGitRemoteURL(userDir *string) (string, error) {
	var repoURL string
	config, err := exec.Command("git", "-C", *userDir, "config", "--get", "remote.origin.url").CombinedOutput()
	if err != nil {
		log.Fatalln(err)
		return "", err
	}
	repoURL = strings.Trim(string(config[:]), "\r\n")
	repoURL = strings.TrimSuffix(repoURL, ".git")
	return repoURL, nil
}
