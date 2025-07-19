package basic

import (
	"time"
)

var (
	DayDuration = time.Hour * 24
)

// CheckTime 检查时间是否在边界之内
func CheckTime(startBoundary, endBoundary, timeStamp time.Time) bool {
	return !(timeStamp.Before(startBoundary) || timeStamp.After(endBoundary))
}

// CheckTimeRange 检查start、end所表示的时间范围是否在指定的边界之内
func CheckTimeRange(startBoundary, endBoundary, start, end time.Time) bool {
	return !(start.Before(startBoundary) || end.After(endBoundary) || start.After(end))
}

// CheckDurationLimit ...
func CheckDurationLimit(startBoundary, endBoundary time.Time, maxDuration time.Duration) bool {
	return startBoundary.Equal(endBoundary) ||
		(startBoundary.Before(endBoundary) && startBoundary.Add(maxDuration).After(endBoundary))
}

// GetTimestampByTimeStr 根据当地时区时间字符串获取当前时间戳
// 2022-04-25 00:00:00  ecmbasic.TIMEFORMAT
func GetTimestampByTimeStr(timeStr string, format string) int64 {
	stamp, _ := time.ParseInLocation(format, timeStr, time.Local)
	return stamp.Unix()
}

// GetTimeStrBeforeNMins 获取当前时间前N分钟的时间字符串
func GetTimeStrBeforeNMins(timeStr string, format string, before int) string {
	stamp, _ := time.ParseInLocation(format, timeStr, time.Local)
	return stamp.Add(-time.Duration(before) * time.Minute).Format(TIMEFORMAT)
}

// GetCurrentTimeStr 获取当前时间的时间字符串
func GetCurrentTimeStr(format string) string {
	return time.Now().Format(format)
}

// Today 获取今天零点的时间
func Today() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
}
