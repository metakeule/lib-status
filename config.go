package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/go-on/queue/q"
)

var (
	home       = os.Getenv("HOME")
	gopath     = os.Getenv("GOPATH")
	configPath = filepath.Join(home, ".lib-status.json")
	config     = &Config{Mutex: &sync.Mutex{}}
)

type Config struct {
	LibPaths    []string
	*sync.Mutex `json:"-"`
}

// findLibraryPath returns the position of the given path in the LibPaths
// slice of the config. If the path is not found, it returns -1.
func (c *Config) findLibraryPath(path string) (pos int) {
	c.Lock()
	defer c.Unlock()
	for i, lp := range c.LibPaths {
		if lp == path {
			return i
		}
	}
	return -1
}

func (c *Config) addLibraryPath(path string) {
	// if the path is already in the config.LibPaths, don't do anything
	if c.findLibraryPath(path) >= 0 {
		return
	}
	c.Lock()
	defer c.Unlock()
	c.LibPaths = append(c.LibPaths, path)
	c.write()
}

func (c *Config) deleteLibraryPath(path string) {
	pos := c.findLibraryPath(path)

	// not found, do nothing
	if pos == -1 {
		return
	}

	c.Lock()
	defer c.Unlock()

	// SNIPPET delete, stolen from http://code.google.com/p/go-wiki/wiki/SliceTricks
	// deletes the element at position pos from slice c.LibPaths
	copy(c.LibPaths[pos:], c.LibPaths[pos+1:])
	c.LibPaths[len(c.LibPaths)-1] = ""
	c.LibPaths = c.LibPaths[:len(c.LibPaths)-1]
	// /SNIPPET delete

	c.write()
}

func (c *Config) write() {
	q.Err(q.PANIC)(
		json.Marshal, c,
	)(
		ioutil.WriteFile, configPath, q.V, os.FileMode(0644),
	).Run()
}

func (c *Config) read() {
	// fmt.Println(configPath)
	_, err := os.Stat(configPath)

	if err == nil {
		// fmt.Println("config found")
		q.Err(q.PANIC)(
			ioutil.ReadFile, configPath,
		)(
			json.Unmarshal, q.V, c,
		).Run()
	} else {
		// fmt.Println("create config")
		c.write()
		// writeConfig()
	}

	// fmt.Printf("%#v\n", config)
}

func init() {
	config.read()
}
