// Package configstore provides an object to fetch config data from disk
// according to PaaSTA conventions. Loaded values are cached to avoid repeated
// disk access. For more details, see docs for `configstore.Store`.
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
	"gopkg.in/yaml.v2"
)

// Store is a container that fetches and caches config values from disk and
// destructures them into PaaSTA value structures.
//
// `Store` object will mimic config loading from original paasta-tools, which
// works as following for `store.Get("foo")` call:
//
//  1. look on disk for file `foo.json` or `foo.yaml`
//  2. if file exists, load it and fetch `foo` key from top-level dictionary
//  3. if file is missing, load all `.json` or `.yaml` files and merge them
//     into single dictionary and look for `foo` key in there
//
// To avoid eagerly loading all existing configuration files, `Store` object
// accepts optional `hints` dictionary, with mapping from keys to file paths,
// where to look for those keys. If requested key is missing from hints, it will
// trigger eager loading as per default functionality.
//
// There are two ways to get config data, via `Get` or `Load` methods. `Get`
// will parse the config value and return it as `interface{}` type, while `Load`
// will accept destination pointer and use `mapstructure.Decode` to destructure
// the value.
//
// TODO: since `Store` is meant to be a long-lived object we need to keep track
// of updated configs and a possibility to manually reset the cache.
type Store struct {
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
		return nil, fmt.Errorf("Failed to list directory %v: %v", dirname, err)
	}

	ret := make([]string, len(files))
	for idx, file := range files {
		ret[idx] = file.Name()
	}

	return ret, nil
}

func parseFile(filepath string, value interface{}) error {
	reader, err := os.Open(filepath)
	defer reader.Close()
	if err != nil {
		return fmt.Errorf("Failed to open %s: %v", filepath, err)
	}

	ext := path.Ext(filepath)
	switch ext {
	case ".json":
		return json.NewDecoder(reader).Decode(value)
	case ".yaml":
		return yaml.NewDecoder(reader).Decode(value)
	default:
		return fmt.Errorf("unknown extension: %v", ext)
	}
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
	err := s.ParseFile(path, &value)
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

	return nil
}

// Get returns value for given `key`. If not found in `s.data`, call
// `s.load` function with `file` from `s.hints` or `key` itself.
func (s *Store) Get(key string) (interface{}, bool, error) {
	if val, ok := s.Data.Load(key); ok {
		return val, ok, nil
	}

	var file string
	var fromHint bool
	if val, ok := s.Hints[key]; ok {
		file = val
		fromHint = true
	} else {
		file = key
		fromHint = false
	}
	err := s.load(file)
	if err != nil {
		return nil, false, fmt.Errorf("Failed to load %v: %v", file, err)
	}

	val, ok := s.Data.Load(key)
	if !ok {
		if !fromHint {
			log.Printf(
				"WARN: loading all configs, consider adding some hints in %s",
				path.Join(s.Dir, file),
			)
			err := s.loadAll()
			if err != nil {
				return nil, false, fmt.Errorf("failed to load all configs: %v", err)
			}
			val, ok = s.Data.Load(key)
		}
	}

	return val, ok, nil
}

// Load uses mapstructure.Decode to parse result of a Get into provided
// destination value
func (s *Store) Load(key string, dst interface{}) (bool, error) {
	val, ok, err := s.Get(key)
	if err != nil {
		return false, fmt.Errorf("Failed to get %s: %v", key, err)
	}
	if !ok {
		return false, nil
	}
	return true, mapstructure.Decode(val, dst)
}
