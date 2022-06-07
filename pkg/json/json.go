package json

import (
	"bytes"
	"sort"

	"github.com/json-iterator/go"
)

// 定义JSON操作
var (
	standard            = jsoniter.ConfigCompatibleWithStandardLibrary
	fast                = jsoniter.ConfigFastest
	unescape            = jsoniter.Config{EscapeHTML: false}.Froze()
	Get                 = standard.Get
	Valid               = standard.Valid
	NewDecoder          = standard.NewDecoder
	NewEncoder          = standard.NewEncoder
	Marshal             = standard.Marshal
	MarshalToString     = standard.MarshalToString
	MarshalIndent       = standard.MarshalIndent
	Unmarshal           = standard.Unmarshal
	UnmarshalFromString = standard.UnmarshalFromString
)

// MustTransform 转化数据
func MustTransform(data, out any) {
	bs := MustEncode(data)
	MustDecode(bs, out)
}

// MustEncode 编码
func MustEncode(data any) []byte {
	bs, err := Marshal(data)
	if err != nil {
		panic(err)
	}
	return bs
}

// MustEncodeToString 转为json字符串
func MustEncodeToString(data any) string {
	str, err := MarshalToString(data)
	if err != nil {
		panic(err)
	}
	return str
}

// MustEncodeFast 没有排序和UnescapeHTML
func MustEncodeFast(data any) string {
	str, err := fast.MarshalToString(data)
	if err != nil {
		panic(err)
	}
	return str
}

func MustDecodeFast(data string, out any) {
	err := fast.UnmarshalFromString(data, out)
	if err != nil {
		panic(err)
	}
}

func MustEncodeWithUnescapeHTML(data any) string {
	str, err := unescape.MarshalToString(data)
	if err != nil {
		panic(err)
	}
	return str
}

// func MustEncodeWithUnescapeHTML(data any) string {
// 	var buf bytes.Buffer
// 	encoder := NewEncoder(&buf)
// 	encoder.SetEscapeHTML(false)
// 	err := encoder.Encode(data)
// 	if err != nil {
// 		panic(err)
// 	}
// 	out := buf.String()
// 	if out[len(out)-1:] == "\n" {
// 		out = out[:len(out)-1]
// 	}
// 	return out
// }

// MustEncodePretty 编码
func MustEncodePretty(data any) []byte {
	bs, err := MarshalIndent(data, "", "  ")
	if err != nil {
		panic(err)
	}
	return bs
}

// MustEncodePrettyString 编码
func MustEncodePrettyString(data any) string {
	bs := MustEncodePretty(data)
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
	err := UnmarshalFromString(data, out)
	if err != nil {
		panic(err)
	}
}

// OrderMap  有序Map
type OrderMap struct {
	Map   map[string]any
	Order []string
}

func NewOrderMap(m map[string]any, order []string) *OrderMap {
	return &OrderMap{
		Map:   m,
		Order: order,
	}
}

func (om *OrderMap) Unmarshal(b []byte) {
	err := Unmarshal(b, &om.Map)
	if err != nil {
		panic(err)
	}

	index := make(map[string]int)
	for key := range om.Map {
		om.Order = append(om.Order, key)
		esc, _ := Marshal(key) // Escape the key
		index[key] = bytes.Index(b, esc)
	}

	sort.Slice(om.Order, func(i, j int) bool { return index[om.Order[i]] < index[om.Order[j]] })
}

func (om OrderMap) marshal() *bytes.Buffer {
	buf := new(bytes.Buffer)
	// var b []byte
	// buf := bytes.NewBuffer(b)
	buf.WriteRune('{')
	l := len(om.Order)
	for i, key := range om.Order {
		km, err := fast.Marshal(key)
		if err != nil {
			panic(err)
		}
		buf.Write(km)
		buf.WriteRune(':')
		if v, ok := om.Map[key].(*OrderMap); ok {
			buf.Write(v.Marshal())
		} else if v, ok := om.Map[key].(OrderMaps); ok {
			buf.Write(v.Marshal())
		} else {
			vm, err := fast.Marshal(om.Map[key])
			if err != nil {
				panic(err)
			}
			buf.Write(vm)
		}
		if i != l-1 {
			buf.WriteRune(',')
		}
	}
	buf.WriteRune('}')
	// fmt.Println(buf.String())
	return buf
}

func (om OrderMap) Marshal() []byte {
	return om.marshal().Bytes()
}

func (om OrderMap) MarshalToString() string {
	return om.marshal().String()
}

type OrderMaps []*OrderMap

func (list OrderMaps) marshal() *bytes.Buffer {
	buf := new(bytes.Buffer)
	buf.WriteRune('[')
	for _, om := range list {
		buf.Write(om.Marshal())
	}
	buf.WriteRune(']')
	return buf
}

func (list OrderMaps) Marshal() []byte {
	return list.marshal().Bytes()
}

func (list OrderMaps) MarshalToString() string {
	return list.marshal().String()
}
