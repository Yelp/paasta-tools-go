package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sync"
)

// Store holds config data
type Store struct {
	// If you care about sanity, never write here, just read
	data  map[string]interface{}
	dir   string
	hints map[string]string
	sync.Mutex

	listFiles  func(string) ([]string, error)
	parseFile  func(string, interface{}) error
	fileExists func(string) (bool, error)
}

func listFiles(dirname string) ([]string, error) {
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		return nil, err
	}

	ret := make([]string, len(files))
	for idx, file := range files {
		ret[idx] = file.Name()
	}

	return ret, nil
}

func parseFile(path string, value interface{}) error {
	reader, err := os.Open(path)
	defer reader.Close()
	if err != nil {
		return err
	}

	return json.NewDecoder(reader).Decode(value)
}

func fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// NewStore creates a new config store
// `dir` is a directory from where data will be loaded
// `hints` is a dictionary to find which file to load for
// a given key
func NewStore(dir string, hints map[string]string) *Store {
	if hints == nil {
		hints = map[string]string{}
	}
	return &Store{
		dir:        dir,
		hints:      hints,
		listFiles:  listFiles,
		parseFile:  parseFile,
		fileExists: fileExists,
	}
}

// Decode `path` contents using `json`, copy `s.data` into
// a new map, merge loaded data into the copy and swap `s.data`
func (s *Store) loadPath(path string) error {
	value := map[string]interface{}{}
	err := s.parseFile(path, value)
	if err != nil {
		return err
	}

	s.Lock()
	defer s.Unlock()
	for key, val := range value {
		s.data[key] = val
	}

	return nil
}

// Walk `s.dir` and `s.loadPath` all the files
func (s *Store) loadAll() error {
	files, err := s.listFiles(s.dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		err := s.loadPath(path.Join(s.dir, file))
		if err != nil {
			return err
		}
	}
	return nil
}

// Look for `domain`.json file, if not found try loading all files and
// print a warning about hints
func (s *Store) load(domain string) error {
	path := path.Join(s.dir, fmt.Sprintf("%s.json", domain))
	exists, err := s.fileExists(path)
	if err != nil {
		return err
	}
	if !exists {
		log.Printf(
			"WARN: loading all configs, consider adding some hints" +
				fmt.Sprintf("for %s in %s", domain, s.dir),
		)
		return s.loadAll()
	}
	return s.loadPath(path)
}

// Get returns value for given `key`. If not found in `s.data`, call
// `s.load` function with `domain` from `s.hints` or `key` itself.
func (s *Store) Get(key string) (interface{}, bool) {
	if val, ok := s.data[key]; ok {
		return val, ok
	}

	var domain string
	if dom, ok := s.hints[key]; ok {
		domain = dom
	} else {
		domain = key
	}
	s.load(domain)

	val, ok := s.data[key]
	return val, ok
}
