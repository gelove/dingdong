package yaml

import (
	yml "gopkg.in/yaml.v2"
)

var (
	NewDecoder = yml.NewDecoder
	NewEncoder = yml.NewEncoder
	Marshal    = yml.Marshal
	Unmarshal  = yml.Unmarshal
)

// MustTransform 转化数据
func MustTransform(data, out any) {
	bytes := MustEncode(data)
	MustDecode(bytes, out)
}

// MustEncode 编码
func MustEncode(data any) []byte {
	bytes, err := Marshal(data)
	if err != nil {
		panic(err)
	}
	return bytes
}

// MustEncodeToString 转为json字符串
func MustEncodeToString(data any) string {
	bs := MustEncode(data)
	return string(bs)
}

// MustDecode 解码
func MustDecode(data []byte, out any) {
	err := Unmarshal(data, out)
	if err != nil {
		panic(err)
	}
}

// MustDecodeFromString 从字符串转为json对象
func MustDecodeFromString(data string, out any) {
	MustDecode([]byte(data), out)
}
