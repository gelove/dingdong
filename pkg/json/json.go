package json

import (
	"github.com/json-iterator/go"
)

// 定义JSON操作
var (
	target              = jsoniter.ConfigCompatibleWithStandardLibrary
	Get                 = target.Get
	Valid               = target.Valid
	NewDecoder          = target.NewDecoder
	NewEncoder          = target.NewEncoder
	Marshal             = target.Marshal
	MarshalToString     = target.MarshalToString
	MarshalIndent       = target.MarshalIndent
	Unmarshal           = target.Unmarshal
	UnmarshalFromString = target.UnmarshalFromString
)

// MustTransform 转化数据
func MustTransform(data, out interface{}) {
	bytes := MustEncode(data)
	MustDecode(bytes, out)
}

// MustEncode 编码
func MustEncode(data interface{}) []byte {
	bytes, err := Marshal(data)
	if err != nil {
		panic(err)
	}
	return bytes
}

// MustEncodePretty 编码
func MustEncodePretty(data interface{}) []byte {
	bytes, err := MarshalIndent(data, "", "  ")
	if err != nil {
		panic(err)
	}
	return bytes
}

// MustEncodeToString 转为json字符串
func MustEncodeToString(data interface{}) string {
	str, err := MarshalToString(data)
	if err != nil {
		panic(err)
	}
	return str
}

// MustDecode 解码
func MustDecode(data []byte, out interface{}) {
	err := Unmarshal(data, out)
	if err != nil {
		panic(err)
	}
}

// MustDecodeFromString 从字符串转为json对象
func MustDecodeFromString(data string, out interface{}) {
	err := UnmarshalFromString(data, out)
	if err != nil {
		panic(err)
	}
}
