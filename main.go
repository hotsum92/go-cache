package main

import (
	"encoding/gob"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"reflect"
	"runtime"
)

func main() {
	v := Cache(fun)

	println(v)
}

func fun() string {
	return "vim-go"
}

func Cache[T any](fn func() T) T {
	usr, err := user.Current()

	if err != nil {
		panic(err)
	}

	cacheDir := filepath.Join(usr.HomeDir, ".cache")

	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		if err := os.Mkdir(cacheDir, 0755); err != nil {
			panic(err)
		}
	}

	funName, err := FuncName(fn)

	if err != nil {
		panic(err)
	}

	filePath := filepath.Join(cacheDir, funName)

	if _, err := os.Stat(filePath); os.IsExist(err) {
		file, err := os.Open(filePath)

		if err != nil {
			panic(err)
		}
		defer file.Close()

		decoder := gob.NewDecoder(file)
		var result T
		err = decoder.Decode(&result)

		if err != nil {
			panic(err)
		}

		return result
	}

	file, err := os.Create(filePath)

	if err != nil {
		panic(err)
	}
	defer file.Close()

	result := fn()

	encoder := gob.NewEncoder(file)
	err := encoder.Encode(result)

	if err != nil {
		panic(err)
	}

	return result
}

func VCache(_ func()) {
}

func FuncName(fn interface{}) (string, error) {
	if rf := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()); rf != nil {
		name := rf.Name()
		return name, nil
	}
	return "", fmt.Errorf("unknownFunc")
}

func Print[T any](v T) T {
	println(v)
	return v
}
