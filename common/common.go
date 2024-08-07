package common

import "time"

// GetCurrentTime 获取当前时间
func GetCurrentTime() string {
	// 获取当前时间
	now := time.Now()

	// 定义所需的日期时间格式
	dateTimeFormat := "2006-01-02 15:04:05"

	// 格式化时间
	formattedTime := now.Format(dateTimeFormat)

	// 输出结果
	return formattedTime
}

// GetCurrentTimestamp 获取当前时间戳
func GetCurrentTimestamp() int64 {
	// Get the current time
	currentTime := time.Now()

	// Get the Unix timestamp (in seconds)
	unixTimestamp := currentTime.Unix()

	// Print the result
	return unixTimestamp
}
