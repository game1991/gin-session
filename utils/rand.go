package utils

import (
	"fmt"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandRangeNum 生成区间随机数
func RandRangeNum(min int, max int) int {
	if min > max || max == 0 || min == 0 {
		return max
	}
	return rand.Intn(max-min) + min
}

// IncrStr 返回前置补0的对应数字，length表示字符串总长度，例如"0001",此时length=4,i=1
func IncrStr(i int, length int) string {
	if length == 0 {
		return ""
	}

	/*
		// 方案一
		s := strconv.Itoa(i)
		len := len(s)
		for i := length; i > len; i-- {
			s = "0" + s
		}
		return s

	*/

	// 方案二

	return fmt.Sprintf("%0*d", length, i)

}
