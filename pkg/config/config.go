package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sync"

	"github.com/mitchellh/mapstructure"
)

// Store holds config data
type Store struct {
	// If you care about sanity, never write here, just read
	Data  map[string]interface{}
	Dir   string
	Hints map[string]string
	sync.Mutex

	ListFiles  func(string) ([]string, error)
	ParseFile  func(string, interface{}) error
	FileExists func(string) (bool, error)
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
		Dir:        dir,
		Hints:      hints,
		ListFiles:  listFiles,
		ParseFile:  parseFile,
		FileExists: fileExists,
	}
}

// Decode `path` contents using `json`, copy `s.data` into
// a new map, merge loaded data into the copy and swap `s.data`
func (s *Store) loadPath(path string) error {
	value := map[string]interface{}{}
	err := s.ParseFile(path, value)
	if err != nil {
		return err
	}

	s.Lock()
	defer s.Unlock()
	for key, val := range value {
		s.Data[key] = val
	}

	return nil
}

// Walk `s.dir` and `s.loadPath` all the files
func (s *Store) loadAll() error {
	files, err := s.ListFiles(s.Dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		err := s.loadPath(path.Join(s.Dir, file))
		if err != nil {
			return err
		}
	}
	return nil
}

// Look for `domain`.json file, if not found try loading all files and
// print a warning about hints
func (s *Store) load(domain string) error {
	path := path.Join(s.Dir, fmt.Sprintf("%s.json", domain))
	exists, err := s.FileExists(path)
	if err != nil {
		return err
	}
	if !exists {
		log.Printf(
			"WARN: loading all configs, consider adding some hints" +
				fmt.Sprintf("for %s in %s", domain, s.Dir),
		)
		return s.loadAll()
	}
	return s.loadPath(path)
}

// Get returns value for given `key`. If not found in `s.data`, call
// `s.load` function with `domain` from `s.hints` or `key` itself.
func (s *Store) Get(key string) (interface{}, error) {
	if val, ok := s.Data[key]; ok {
		return val, nil
	}

	var domain string
	if dom, ok := s.Hints[key]; ok {
		domain = dom
	} else {
		domain = key
	}
	s.load(domain)

	val, ok := s.Data[key]
	if !ok {
		return nil, fmt.Errorf("key not found: %s", key)
	}
	return val, nil
}

// Load uses mapstructure.Decode to parse result of a Get into provided
// destination value
func (s *Store) Load(key string, dst interface{}) error {
	val, err := s.Get(key)
	if err != nil {
		return err
	}
	mapstructure.Decode(val, dst)
	return nil
}
