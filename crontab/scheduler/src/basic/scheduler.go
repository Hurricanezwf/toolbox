package basic

import (
	"fmt"
	"time"
)

// weekScheduler 以周为周期的定时调度
type weekScheduler struct {
	minute  uint64 // [0,59]
	hour    uint64 // [0,23]
	weekday uint64 // [1,7]
}

func newWeekScheduler(minute, hour, weekday uint64) *weekScheduler {
	if weekday == uint64(0) {
		panic("weekday can't be zero")
	}
	if hour == uint64(0) {
		panic("hour can't be zero")
	}
	if minute == uint64(0) {
		panic("minute can't be zero")
	}
	return &weekScheduler{
		minute:  minute,
		hour:    hour,
		weekday: weekday,
	}
}

func (s *weekScheduler) Prev(baseTime time.Time) time.Time {
	// << 如果当前星期对不上 >>
	if !isMatch(1<<uint64(baseTime.Weekday()), s.weekday) {
		// 直接找前一个星期的最晚的可执行时间
		t := findPrevWeekday(baseTime, s.weekday)
		return time.Date(t.Year(), t.Month(), t.Day(), findBiggestHour(s.hour), findBiggestMinute(s.minute), 0, 0, t.Location())
	}

	// << 如果当前小时匹配不上 >>
	if !isMatch(1<<uint64(baseTime.Hour()), s.hour) {
		// 找到了另一个可执行的小时
		if prevHour := findPrevHour(baseTime.Hour(), s.hour); prevHour >= 0 {
			return time.Date(baseTime.Year(), baseTime.Month(), baseTime.Day(), prevHour, findBiggestMinute(s.minute), 0, 0, baseTime.Location())
		}
		// 如果未找到另一个可执行的小时，那就找前一个可执行的weekday吧
		t := findPrevWeekday(baseTime, s.weekday)
		return time.Date(t.Year(), t.Month(), t.Day(), findBiggestHour(s.hour), findBiggestMinute(s.minute), 0, 0, t.Location())
	}

	// << 如果当前分钟匹配不上 >>
	if !isMatch(1<<uint64(baseTime.Minute()), s.minute) {
		// 找另一个可执行的分钟
		if prevMinute := findPrevMinute(baseTime.Minute(), s.minute); prevMinute >= 0 {
			return time.Date(baseTime.Year(), baseTime.Month(), baseTime.Day(), baseTime.Hour(), prevMinute, 0, 0, baseTime.Location())
		}
		// 找另一个可执行的小时的最晚的时间
		if prevHour := findPrevHour(baseTime.Hour(), s.hour); prevHour >= 0 {
			return time.Date(baseTime.Year(), baseTime.Month(), baseTime.Day(), prevHour, findBiggestMinute(s.minute), 0, 0, baseTime.Location())
		}
		// 找另一个可执行的weekday的最早时间
		t := findPrevWeekday(baseTime, s.weekday)
		return time.Date(t.Year(), t.Month(), t.Day(), findBiggestHour(s.hour), findBiggestMinute(s.minute), 0, 0, baseTime.Location())
	}
	// << 我靠，所有的都能匹配的上 >>
	return baseTime
}

func (s *weekScheduler) Next(baseTime time.Time) time.Time {
	// << 如果当前星期对不上 >>
	if !isMatch(1<<uint64(baseTime.Weekday()), s.weekday) {
		// 直接找下一个星期的最早的可执行时间
		t := findNextWeekday(baseTime, s.weekday)
		return time.Date(t.Year(), t.Month(), t.Day(), findSmallestHour(s.hour), findSmallestMinute(s.minute), 0, 0, t.Location())
	}

	// << 如果当前小时匹配不上 >>
	if !isMatch(1<<uint64(baseTime.Hour()), s.hour) {
		// 找到了另一个可执行的小时
		if nextHour := findNextHour(baseTime.Hour(), s.hour); nextHour >= 0 {
			return time.Date(baseTime.Year(), baseTime.Month(), baseTime.Day(), nextHour, findSmallestMinute(s.minute), 0, 0, baseTime.Location())
		}
		// 如果未找到另一个可执行的小时，那就找下一个可执行的weekday吧
		t := findNextWeekday(baseTime, s.weekday)
		return time.Date(t.Year(), t.Month(), t.Day(), findSmallestHour(s.hour), findSmallestMinute(s.minute), 0, 0, t.Location())
	}

	// << 如果当前分钟匹配不上 >>
	if !isMatch(1<<uint64(baseTime.Minute()), s.minute) {
		// 找另一个可执行的分钟
		if nextMinute := findNextMinute(baseTime.Minute(), s.minute); nextMinute >= 0 {
			return time.Date(baseTime.Year(), baseTime.Month(), baseTime.Day(), baseTime.Hour(), nextMinute, 0, 0, baseTime.Location())
		}
		// 找另一个可执行的小时的最早的时间
		if nextHour := findNextHour(baseTime.Hour(), s.hour); nextHour >= 0 {
			return time.Date(baseTime.Year(), baseTime.Month(), baseTime.Day(), nextHour, findSmallestMinute(s.minute), 0, 0, baseTime.Location())
		}
		// 找另一个可执行的weekday的最早时间
		t := findNextWeekday(baseTime, s.weekday)
		return time.Date(t.Year(), t.Month(), t.Day(), findSmallestHour(s.hour), findSmallestMinute(s.minute), 0, 0, t.Location())
	}
	// << 我靠，所有的都能匹配的上 >>
	return baseTime
}

// monthScheduler 以月为周期的定时调度
type monthScheduler struct {
	minute   uint64 // [0,59]
	hour     uint64 // [0,23]
	monthday uint64 // [1,31]
}

func newMonthScheduler(minute, hour, monthday uint64) *monthScheduler {
	if monthday == uint64(0) {
		panic("monthday can't be zero")
	}
	if hour == uint64(0) {
		panic("hour can't be zero")
	}
	if minute == uint64(0) {
		panic("minute can't be zero")
	}
	return &monthScheduler{
		minute:   minute,
		hour:     hour,
		monthday: monthday,
	}
}

func (s *monthScheduler) Prev(baseTime time.Time) time.Time {
	// << 如果当前日期对不上 >>
	if !isMatch(1<<uint64(baseTime.Day()), s.monthday) {
		// 直接找前一个最晚的可执行时间
		t := findPrevMonthday(baseTime, s.monthday)
		return time.Date(t.Year(), t.Month(), t.Day(), findBiggestHour(s.hour), findBiggestMinute(s.minute), 0, 0, t.Location())
	}

	// << 如果当前小时匹配不上 >>
	if !isMatch(1<<uint64(baseTime.Hour()), s.hour) {
		// 找到了另一个可执行的小时
		if prevHour := findPrevHour(baseTime.Hour(), s.hour); prevHour >= 0 {
			return time.Date(baseTime.Year(), baseTime.Month(), baseTime.Day(), prevHour, findBiggestMinute(s.minute), 0, 0, baseTime.Location())
		}
		// 如果未找到另一个可执行的小时，那就找前一个可执行的monthday吧
		t := findPrevMonthday(baseTime, s.monthday)
		return time.Date(t.Year(), t.Month(), t.Day(), findBiggestHour(s.hour), findBiggestMinute(s.minute), 0, 0, t.Location())
	}

	// << 如果当前分钟匹配不上 >>
	if !isMatch(1<<uint64(baseTime.Minute()), s.minute) {
		// 找另一个可执行的分钟
		if prevMinute := findPrevMinute(baseTime.Minute(), s.minute); prevMinute >= 0 {
			return time.Date(baseTime.Year(), baseTime.Month(), baseTime.Day(), baseTime.Hour(), prevMinute, 0, 0, baseTime.Location())
		}
		// 找另一个可执行的小时的最晚的时间
		if prevHour := findPrevHour(baseTime.Hour(), s.hour); prevHour >= 0 {
			return time.Date(baseTime.Year(), baseTime.Month(), baseTime.Day(), prevHour, findBiggestMinute(s.minute), 0, 0, baseTime.Location())
		}
		// 找另一个可执行的monthday的最早时间
		t := findPrevMonthday(baseTime, s.monthday)
		return time.Date(t.Year(), t.Month(), t.Day(), findBiggestHour(s.hour), findBiggestMinute(s.minute), 0, 0, baseTime.Location())
	}
	// << 我靠，所有的都能匹配的上 >>
	return baseTime
}

func (s *monthScheduler) Next(baseTime time.Time) time.Time {
	// << 如果当前日期对不上 >>
	if !isMatch(1<<uint64(baseTime.Day()), s.monthday) {
		// 直接找下一个最早的可执行时间
		t := findNextMonthday(baseTime, s.monthday)
		return time.Date(t.Year(), t.Month(), t.Day(), findSmallestHour(s.hour), findSmallestMinute(s.minute), 0, 0, t.Location())
	}

	// << 如果当前小时匹配不上 >>
	if !isMatch(1<<uint64(baseTime.Hour()), s.hour) {
		// 找到了另一个可执行的小时
		if nextHour := findNextHour(baseTime.Hour(), s.hour); nextHour >= 0 {
			return time.Date(baseTime.Year(), baseTime.Month(), baseTime.Day(), nextHour, findSmallestMinute(s.minute), 0, 0, baseTime.Location())
		}
		// 如果未找到另一个可执行的小时，那就找下一个可执行的monthday吧
		t := findNextMonthday(baseTime, s.monthday)
		return time.Date(t.Year(), t.Month(), t.Day(), findSmallestHour(s.hour), findSmallestMinute(s.minute), 0, 0, t.Location())
	}

	// << 如果当前分钟匹配不上 >>
	if !isMatch(1<<uint64(baseTime.Minute()), s.minute) {
		// 找另一个可执行的分钟
		if nextMinute := findNextMinute(baseTime.Minute(), s.minute); nextMinute >= 0 {
			return time.Date(baseTime.Year(), baseTime.Month(), baseTime.Day(), baseTime.Hour(), nextMinute, 0, 0, baseTime.Location())
		}
		// 找另一个可执行的小时的最早的时间
		if nextHour := findNextHour(baseTime.Hour(), s.hour); nextHour >= 0 {
			return time.Date(baseTime.Year(), baseTime.Month(), baseTime.Day(), nextHour, findSmallestMinute(s.minute), 0, 0, baseTime.Location())
		}
		// 找另一个可执行的monthday的最早时间
		t := findNextMonthday(baseTime, s.monthday)
		return time.Date(t.Year(), t.Month(), t.Day(), findSmallestHour(s.hour), findSmallestMinute(s.minute), 0, 0, t.Location())
	}
	// << 我靠，所有的都能匹配的上 >>
	return baseTime
}

func isMatch(v1, v2 uint64) bool {
	return v1&v2 != 0
}

func findSmallestHour(hourBits uint64) int {
	for i := 1; i <= 23; i++ {
		if isMatch(hourBits, 1<<uint64(i)) {
			return i
		}
	}
	panic(fmt.Sprintf("find smallest hour failed, not found, hour is %d", hourBits))
	return -1
}

func findBiggestHour(hourBits uint64) int {
	for i := 23; i >= 1; i-- {
		if isMatch(hourBits, 1<<uint64(i)) {
			return i
		}
	}
	panic(fmt.Sprintf("find latest hour failed, not found, hour is %d", hourBits))
	return -1
}

func findSmallestMinute(minuteBits uint64) int {
	for i := 1; i <= 59; i++ {
		if isMatch(minuteBits, 1<<uint64(i)) {
			return i
		}
	}
	panic(fmt.Sprintf("find smallest minute failed, not found, minute is %d", minuteBits))
	return -1
}

func findBiggestMinute(minuteBits uint64) int {
	for i := 59; i >= 0; i-- {
		if isMatch(minuteBits, 1<<uint64(i)) {
			return i
		}
	}
	panic(fmt.Sprintf("find biggest minute failed, not found, minute is %d", minuteBits))
	return -1
}

func findNextHour(baseHour int, hourBits uint64) int {
	for i := baseHour + 1; i <= 23; i++ {
		if isMatch(hourBits, 1<<uint64(i)) {
			return i
		}
	}
	return -1
}

func findPrevHour(baseHour int, hourBits uint64) int {
	for i := baseHour - 1; i >= 0; i-- {
		if isMatch(hourBits, 1<<uint64(i)) {
			return i
		}
	}
	return -1
}

func findNextMinute(baseMinute int, hourBits uint64) int {
	for i := baseMinute + 1; i <= 59; i++ {
		if isMatch(hourBits, 1<<uint64(i)) {
			return i
		}
	}
	return -1
}

func findPrevMinute(baseMinute int, hourBits uint64) int {
	for i := baseMinute - 1; i >= 0; i-- {
		if isMatch(hourBits, 1<<uint64(i)) {
			return i
		}
	}
	return -1
}

func findNextWeekday(baseTime time.Time, weekdayBits uint64) time.Time {
	t := baseTime
	for i := 1; i <= 7; i++ {
		t = t.Add(24 * time.Hour)
		if isMatch(1<<uint64(t.Weekday()), weekdayBits) {
			return t
		}
	}
	panic(fmt.Sprintf("find next weekday failed, not found, weekday=%d", weekdayBits))
	return time.Time{}
}

func findPrevWeekday(baseTime time.Time, weekdayBits uint64) time.Time {
	t := baseTime
	for i := 1; i <= 7; i++ {
		t = t.Add(-24 * time.Hour)
		if isMatch(1<<uint64(t.Weekday()), weekdayBits) {
			return t
		}
	}
	panic(fmt.Sprintf("find prev weekday failed, not found, weekday=%d", weekdayBits))
	return time.Time{}
}

func findNextMonthday(baseTime time.Time, monthdayBits uint64) time.Time {
	// 注意：不同月份对应的天数不一样
	t := baseTime
	for i := 0; i <= 365; i++ {
		t = t.Add(24 * time.Hour)
		if isMatch(1<<uint64(t.Day()), monthdayBits) {
			return t
		}
	}
	panic(fmt.Sprintf("find next monthday failed, not found, monthday=%d", monthdayBits))
	return time.Time{}
}

func findPrevMonthday(baseTime time.Time, monthdayBits uint64) time.Time {
	t := baseTime
	for i := 1; i <= 365; i++ {
		t = t.Add(-24 * time.Hour)
		if isMatch(1<<uint64(t.Day()), monthdayBits) {
			return t
		}
	}
	panic(fmt.Sprintf("find prev monthday failed, not found, monthday=%d", monthdayBits))
	return time.Time{}
}
