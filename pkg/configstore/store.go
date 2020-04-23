package configstore

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
	Data  *sync.Map
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
		return fmt.Errorf("Failed to open %s: %v", path, err)
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
		Data:       &sync.Map{},
		Dir:        dir,
		Hints:      hints,
		ListFiles:  listFiles,
		ParseFile:  parseFile,
		FileExists: fileExists,
	}
}

// Decode `path` contents using `json`, lock the store mutex, merge loaded data
// into s.Data, unlock the mutex
func (s *Store) loadPath(path string) error {
	value := map[string]interface{}{}
	err := s.ParseFile(path, value)
	if err != nil {
		return fmt.Errorf("Failed to parse %s: %v", path, err)
	}

	s.Lock()
	defer s.Unlock()
	for key, val := range value {
		s.Data.Store(key, val)
	}

	return nil
}

// Walk `s.dir` and `s.loadPath` all the files
func (s *Store) loadAll() error {
	files, err := s.ListFiles(s.Dir)
	if err != nil {
		return fmt.Errorf("Failed to list %s: %v", s.Dir, err)
	}

	for _, file := range files {
		filepath := path.Join(s.Dir, file)
		err := s.loadPath(filepath)
		if err != nil {
			return fmt.Errorf("Failed to load %s: %v", filepath, err)
		}
	}
	return nil
}

var extensions = []string{"json", "yaml"}

// Look for `file`.json or `file`.yaml, if not found try loading all files and
// print a warning about hints
func (s *Store) load(file string) error {
	for _, ext := range extensions {
		path := path.Join(s.Dir, fmt.Sprintf("%s.%s", file, ext))
		exists, err := s.FileExists(path)
		if err != nil {
			return fmt.Errorf("Failed to find %s: %v", path, err)
		}
		if exists {
			err = s.loadPath(path)
			if err != nil {
				return fmt.Errorf("Failed to load %s: %v", path, err)
			}
			return nil
		}
	}

	log.Printf(
		"WARN: loading all configs, consider adding some hints in %s",
		path.Join(s.Dir, file),
	)
	return s.loadAll()
}

// Get returns value for given `key`. If not found in `s.data`, call
// `s.load` function with `file` from `s.hints` or `key` itself.
func (s *Store) Get(key string) (interface{}, error) {
	if val, ok := s.Data.Load(key); ok {
		return val, nil
	}

	var file string
	if val, ok := s.Hints[key]; ok {
		file = val
	} else {
		file = key
	}
	s.load(file)

	val, ok := s.Data.Load(key)
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
		return fmt.Errorf("Failed to get %s: %v", key, err)
	}
	return mapstructure.Decode(val, dst)
}
