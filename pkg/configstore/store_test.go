package configstore

import (
	"fmt"
	"sync"
	"testing"
)

func errorIf(test *testing.T, pred bool, format string, args ...interface{}) {
	if pred {
		test.Fatalf(format, args...)
	}
}

func errorUnexpected(test *testing.T, expected, actual interface{}) {
	errorIf(test, expected != actual, "expected %+v, actual %+v", expected, actual)
}

func unexpectedParseFile(test *testing.T) func(string, interface{}) error {
	return func(file string, val interface{}) error {
		test.Fatalf("unexpected call to parseFile(%s, _)", file)
		return nil
	}
}

func unexpectedListFiles(test *testing.T) func(string) ([]string, error) {
	return func(dirname string) ([]string, error) {
		test.Fatalf("unexpected call to listFiles(%s)", dirname)
		return []string{}, nil
	}
}

func TestStore_loadPath(test *testing.T) {
	key := "test key"
	s := &Store{
		Data: &sync.Map{},
		ParseFile: func(file string, val interface{}) error {
			v, ok := val.(*map[string]interface{})
			if !ok {
				panic("assert failed")
			}
			(*v)[key] = file
			return nil
		},
	}
	expected := "old value"
	s.Data.Store(key, expected)
	s.loadPath("new value")

	if actual, ok := s.Data.Load(key); ok {
		if actual != "new value" {
			test.Fatalf("%+v was expected, got %+v", expected, actual)
		}
	} else {
		test.Fatalf("key %+v not found", key)
	}
}

func TestStore_loadAll(test *testing.T) {
	s := &Store{
		Dir:  "zero",
		Data: &sync.Map{},
		ParseFile: func(file string, val interface{}) error {
			fmt.Printf("parse file called: %s\n", file)
			v, ok := val.(*map[string]interface{})
			if !ok {
				panic("assert failed")
			}
			(*v)[file] = "loaded"
			return nil
		},
		ListFiles: func(dirname string) ([]string, error) {
			return []string{"one", "two"}, nil
		},
	}
	s.loadAll()

	for _, key := range []string{"zero/one", "zero/two"} {
		if actual, ok := s.Data.Load(key); ok {
			if actual != "loaded" {
				test.Fatalf("%s wasn't loaded correctly", key)
			}
		} else {
			test.Fatalf("key %+v wasn't loaded", key)
		}
	}
}

func TestStore_load(test *testing.T) {
	s := &Store{
		Dir:       "zero",
		Data:      &sync.Map{},
		ParseFile: func(file string, val interface{}) error { return nil },
		ListFiles: unexpectedListFiles(test),
		FileExists: func(path string) (bool, error) {
			expected := "zero/one.json"
			if path != expected {
				test.Fatalf("expected path=%s, got path=%s", expected, path)
			}
			return true, nil
		},
	}
	s.load("one")

	s.FileExists = func(path string) (bool, error) { return false, nil }
	s.ParseFile = unexpectedParseFile(test)
	s.ListFiles = func(dirname string) ([]string, error) { return []string{}, nil }
	s.load("one")
}

func TestStore_Get(test *testing.T) {
	s := &Store{
		Dir:       "zero",
		Data:      &sync.Map{},
		ParseFile: unexpectedParseFile(test),
		ListFiles: unexpectedListFiles(test),
		FileExists: func(path string) (bool, error) {
			test.Fatalf("unexpected call to fileExists(%s)", path)
			return true, nil
		},
	}
	s.Data.Store("one", "two")
	val, ok, err := s.Get("one")
	errorIf(test, !ok, "key not found")
	errorIf(test, err != nil, fmt.Sprintf("err when loading: %v", err))
	errorUnexpected(test, "two", val)

	// key is missing, file with same name exists
	s.FileExists = func(path string) (bool, error) { return true, nil }
	s.ParseFile = func(file string, val interface{}) error {
		v := val.(*map[string]interface{})
		(*v)["two"] = "three"
		return nil
	}
	val, ok, err = s.Get("two")
	errorIf(test, !ok, "key not found")
	errorIf(test, err != nil, fmt.Sprintf("err when loading: %v", err))
	errorUnexpected(test, "three", val)

	// key is missing, file corresponding to a hint exists
	s.Hints = map[string]string{"three": "four"}
	s.FileExists = func(path string) (bool, error) {
		errorUnexpected(test, "zero/four.json", path)
		return true, nil
	}
	s.ParseFile = func(file string, val interface{}) error {
		errorUnexpected(test, "zero/four.json", file)
		v := val.(*map[string]interface{})
		(*v)["three"] = "four"
		return nil
	}
	val, ok, err = s.Get("three")
	errorIf(test, !ok, "key not found")
	errorIf(test, err != nil, fmt.Sprintf("err when loading: %v", err))
	errorUnexpected(test, "four", val)
}

func TestStore_GetLoadAllWhenNoKey(test *testing.T) {
	s := &Store{
		Dir:       "zero",
		Data:      &sync.Map{},
		ParseFile: unexpectedParseFile(test),
		ListFiles: unexpectedListFiles(test),
		FileExists: func(path string) (bool, error) {
			test.Fatalf("unexpected call to fileExists(%s)", path)
			return true, nil
		},
	}

	// key is missing, file is missing, hint is missing, loaded from all
	listFilesCalled := false
	s.ListFiles = func(string) ([]string, error) {
		listFilesCalled = true
		return []string{}, nil
	}
	s.FileExists = func(string) (bool, error) {
		return false, nil
	}
	s.Get("four")
	errorIf(test, !listFilesCalled, "listFiles wasn't called")
}

func TestStore_GetCallsListAllWhenMatchingFileHasNoKey(test *testing.T) {
	s := &Store{
		Dir:       "zero",
		Data:      &sync.Map{},
		ParseFile: unexpectedParseFile(test),
		ListFiles: unexpectedListFiles(test),
		FileExists: func(path string) (bool, error) {
			test.Fatalf("unexpected call to fileExists(%s)", path)
			return true, nil
		},
	}

	// file exists, but key in different file
	listFilesCalled := false
	s.ListFiles = func(string) ([]string, error) {
		listFilesCalled = true
		return []string{}, nil
	}
	s.FileExists = func(string) (bool, error) {
		return true, nil
	}
	s.ParseFile = func(file string, val interface{}) error {
		v := val.(*map[string]interface{})
		(*v)["five"] = "five"
		return nil
	}
	s.Get("five")
	errorIf(test, listFilesCalled, "listFiles wasn't called")
}
