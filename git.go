package main

import (
	"bytes"
	"os/exec"
	"strings"
)

func gitStatusForRepo(repoPath string) (res string, err error) {
	gitCmd := exec.Command(
		"git", "status",
		"--porcelain", // easy to parseable output
		"-b",          // show if branch is ahead of origin
		"--untracked-files=no", // don't show untracked-files
	)
	gitCmd.Dir = repoPath
	var data []byte
	data, err = gitCmd.Output()
	if err != nil {
		return
	}

	res = strings.TrimSpace(string(data))

	if res == "## master...origin/master" {
		res = ""
	}
	return
}

func gitDiff(repoPath string) (data []byte, err error) {
	gitCmd := exec.Command(
		"git", "diff",
		//"--porcelain",          // easy to parseable output
		// "--untracked-files=no", // don't show untracked-files
	)
	gitCmd.Dir = repoPath
	data, err = gitCmd.Output()
	data = bytes.TrimSpace(data)
	return
}

func gitPush(repoPath string) (data []byte, err error) {
	gitCmd := exec.Command(
		"git", "push",
		"origin", "master",
	)

	gitCmd.Dir = repoPath
	data, err = gitCmd.Output()
	return
}
