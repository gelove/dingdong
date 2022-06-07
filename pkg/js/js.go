package js

import (
	"sync"

	"github.com/robertkrimen/otto"

	"dingdong/assets"
)

var cache sync.Map

func Call(filePath string, functionName string, args ...any) (result otto.Value, err error) {
	if v, ok := cache.Load(filePath); ok {
		return call(v.([]byte), functionName, args...)
	}

	bs, err := assets.ReadFile(filePath)
	if err != nil {
		return
	}
	cache.Store(filePath, bs)
	return call(bs, functionName, args...)
}

func call(src []byte, functionName string, args ...any) (result otto.Value, err error) {
	vm := otto.New()
	_, err = vm.Run(src)
	if err != nil {
		return
	}
	result, err = vm.Call(functionName, nil, args...)
	return
}
