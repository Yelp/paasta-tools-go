package config

import (
	"os"
)

func ReadSystemPaaSTAConfig(f string) *os.File {
	r, err := os.Open(f)
	if err != nil {
		panic(err)
	}
	return r
}
