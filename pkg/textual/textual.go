package textual

import (
	cRand "crypto/rand"
	"encoding/hex"
	"io"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

const (
	letterBytes   = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

// RandomString 获取指定长度随机字符串
func RandomString(length int) string {
	b := make([]byte, length)
	src := rand.NewSource(time.Now().UnixNano())
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := length-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}

// RandomKey 获取指定长度随机字符串
func RandomKey(length int) string {
	k := make([]byte, length)
	if _, err := io.ReadFull(cRand.Reader, k); err != nil {
		panic(err)
	}
	s := hex.EncodeToString(k)
	return s[:length]
}

// Split 分割字符串
func Split(s, sep string) []string {
	list := strings.Split(s, sep)
	res := make([]string, 0, len(list))
	for _, v := range list {
		res = append(res, strings.TrimSpace(v))
	}
	return res
}

// ArrayShift shift an element off the beginning of slice
func ArrayShift(list *[]string) string {
	if len(*list) == 0 {
		return ""
	}
	elem := (*list)[0]
	*list = (*list)[1:]
	return elem
}

// ArrayPop pop an element off the last of slice
func ArrayPop(list *[]string) string {
	l := len(*list)
	if l == 0 {
		return ""
	}
	elem := (*list)[l-1]
	*list = (*list)[:l-1]
	return elem
}

// InArray 判断是否在切片中
func InArray(needle string, list []string) bool {
	needle = strings.TrimSpace(needle)
	for _, v := range list {
		if strings.TrimSpace(v) == needle {
			return true
		}
	}
	return false
}

// IndexOf 获取在切片中的索引 不存在则为-1
func IndexOf(needle string, list []string) int {
	needle = strings.TrimSpace(needle)
	for i, v := range list {
		if strings.TrimSpace(v) == needle {
			return i
		}
	}
	return -1
}

// PrefixInArray 判断是否在切片中, 或前缀在切片中
func PrefixInArray(needle string, list []string) bool {
	needle = strings.TrimSpace(needle)
	for _, v := range list {
		v = strings.TrimSpace(v)
		if v == needle || strings.HasPrefix(needle, v) {
			return true
		}
	}
	return false
}

func PrefixIndexOf(needle string, list []string) int {
	needle = strings.TrimSpace(needle)
	for i, v := range list {
		v = strings.TrimSpace(v)
		if v == needle || strings.HasPrefix(needle, v) {
			return i
		}
	}
	return -1
}

// TrimSpace 去除元素前后的空格
func TrimSpace(list []string) {
	for index := range list {
		list[index] = strings.TrimSpace(list[index])
	}
	return
}

// FilterByWhiteList 必须在白名单内否则过滤
// list []string
// whiteList []string 为空时不用过滤
func FilterByWhiteList(list []string, whiteList []string) []string {
	if len(whiteList) == 0 {
		return list
	}
	return FilterBySlice(list, whiteList, true)
}

// FilterByBlackList 必须在白名单内否则过滤
// list []string
// blackList []string
func FilterByBlackList(list []string, blackList []string) []string {
	return FilterBySlice(list, blackList, false)
}

// FilterBySlice 过滤指定的字符串并过滤空字符串
// list []string
// specified []string
// isWhiteList 白名单还是黑名单
func FilterBySlice(list []string, specified []string, isWhiteList bool) []string {
	ret := make([]string, 0, len(list))
	for _, val := range list {
		val = strings.TrimSpace(val)
		if val == "" {
			continue
		}
		if InArray(val, specified) == isWhiteList {
			ret = append(ret, val)
		}
	}
	return ret
}

// FilterSpace 过滤空格元素
func FilterSpace(list []string) []string {
	ret := make([]string, 0, len(list))
	for _, val := range list {
		val = strings.TrimSpace(val)
		if val != "" {
			ret = append(ret, val)
		}
	}
	return ret
}

// Intersect 判断是否有交集
func Intersect(a []string, b []string) (rs bool) {
	for _, v := range a {
		if InArray(v, b) {
			rs = true
			return
		}
	}
	return
}

// PrefixIntersect 判断前缀是否有交集
func PrefixIntersect(a []string, b []string) (rs bool) {
	for _, v := range a {
		if PrefixInArray(v, b) {
			rs = true
			return
		}
	}
	return
}

// Unique 去重
func Unique(list []string) []string {
	result := make([]string, 0, len(list))
	temp := map[string]struct{}{}
	for _, v := range list {
		v = strings.TrimSpace(v)
		if v == "" {
			continue
		}
		if _, ok := temp[v]; ok {
			continue
		}
		temp[v] = struct{}{}
		result = append(result, v)
	}
	return result
}

// ToIntList 字符串切片转为整形切片
func ToIntList(list []string) []int {
	result := make([]int, 0, len(list))
	for _, v := range list {
		v = strings.TrimSpace(v)
		item, err := strconv.Atoi(v)
		if err != nil {
			panic(err)
		}
		result = append(result, item)
	}
	return result
}

// SortStingNumber 对数字字符串数组排序并转为float64
// isAscend 是否升序
func SortStingNumber(list []string, isAscend bool) []float64 {
	res := make([]float64, 0, len(list))
	for _, v := range list {
		v = strings.TrimSpace(v)
		val, err := strconv.ParseFloat(v, 64)
		if err != nil {
			panic(err)
		}
		res = append(res, val)
	}
	if isAscend {
		sort.Sort(sort.Float64Slice(res))
	} else {
		sort.Sort(sort.Reverse(sort.Float64Slice(res)))
	}
	return res
}
