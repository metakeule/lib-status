package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

func depTrackForRepo(repoPath string) (data []byte, err error) {
	depCmd := exec.Command(
		"dep", "track",
	)
	depCmd.Env = os.Environ()
	depCmd.Dir = repoPath
	data, err = depCmd.Output()
	return
}

func depGDFForRepo(repoPath string) (data []byte, err error) {
	depCmd := exec.Command(
		"dep", "gdf",
	)

	// fmt.Println("os.Environ", os.Environ())
	depCmd.Dir = repoPath
	depCmd.Env = os.Environ()
	data, err = depCmd.Output()
	data = bytes.TrimSpace(data)

	if err != nil {
		fmt.Println(repoPath)
		fmt.Println(err.Error())
		fmt.Println(string(data))
	}

	if string(data) != "" {
		// fmt.Println("empty", repoPath)
	}

	if err != nil && string(data) != "" {
		// fmt.Printf("len gdf: %d\n", len(data))
		gdfFile := filepath.Join(repoPath, "gdf.json")
		fmt.Println("writing", gdfFile)
		err = ioutil.WriteFile(gdfFile, data, os.FileMode(0644))

		if err != nil {
			fmt.Println(err.Error())
		}
	}
	return
}
