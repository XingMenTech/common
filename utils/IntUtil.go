// Int类型工具
// Neo

package utils

import (
	"fmt"
	"strconv"
	"strings"
)

type Int64Array []int64

// /转换为字符串数组
func (object Int64Array) ToStringArray() (arr []string) {
	if nil != object && 0 != len(object) {
		arr = make([]string, len(object))
		for i := 0; i < len(object); i++ {
			arr[i] = fmt.Sprint(object[i])
		}
	}
	return
}

// 字符串转换 int
func Int(ceil string) int {
	i, err := strconv.Atoi(ceil)
	if err != nil {
		return 0
	}
	return i
}

// 字符串转int64
func Int64(ceil string) int64 {
	ib, err := strconv.ParseInt(strings.TrimSpace(ceil), 10, 64)

	//ib, err := strconv.ParseUint(ceil, 10, 32)
	if err == nil {
		return ib
	}
	return int64(0)
}

// 字符串转uint32
func Uint32(ceil string) uint32 {
	ib, err := strconv.ParseUint(ceil, 10, 32)
	if err == nil {
		return uint32(ib)
	}
	return uint32(0)
}
