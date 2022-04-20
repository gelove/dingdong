package js

import (
	"io/ioutil"
	"sync"

	"github.com/robertkrimen/otto"
)

var cache sync.Map

func Call(filePath string, functionName string, args ...interface{}) (result otto.Value, err error) {
	if v, ok := cache.Load(filePath); ok {
		return call(v.([]byte), functionName, args...)
	}

	bs, err := ioutil.ReadFile(filePath)
	if err != nil {
		return
	}
	cache.Store(filePath, bs)
	return call(bs, functionName, args...)
}

func call(src []byte, functionName string, args ...interface{}) (result otto.Value, err error) {
	vm := otto.New()
	_, err = vm.Run(src)
	if err != nil {
		return
	}
	result, err = vm.Call(functionName, nil, args...)
	return
}
