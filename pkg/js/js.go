package js

import (
	"io/ioutil"

	"github.com/robertkrimen/otto"
)

var cache = make(map[string][]byte)

func Call(filePath string, functionName string, args ...interface{}) (result otto.Value, err error) {
	if _, ok := cache[filePath]; !ok {
		cache[filePath], err = ioutil.ReadFile(filePath)
		if err != nil {
			return
		}
	}

	vm := otto.New()
	_, err = vm.Run(cache[filePath])
	if err != nil {
		return
	}
	result, err = vm.Call(functionName, nil, args...)
	return
}
