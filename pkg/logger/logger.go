package logger

import "fmt"

// Error 记录错误日志
func Error(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
