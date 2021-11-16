package types

import (
	"goblog/pkg/logger"
	"strconv"
)

// Uint64ToString 将 uint64 转换为 string
func Uint64ToString(num uint64) string {
	return strconv.FormatUint(num, 10)
}

// StringToUint 将 string 转成 uint64
func StringToUint(str string) (i uint64) {
	i, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		logger.Error(err)
	}
	return
}
