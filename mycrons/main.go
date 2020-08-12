package mycrons

import (
	"fmt"
	"time"
)

type Scheduler interface {
	// Prev 找寻基准时间之前理论上的执行时间
	Prev(baseTime time.Time) time.Time
	// Next 找寻基准时间之后理论上的执行时间
	Next(baseTime time.Time) time.Time
}

type WeekOfDayScheduler struct {
	// minute 可取值 0, 15, 30, 45 中的任意一个或多个
	minute uint64
	// hour 可取值 0-23 中的任意一个或多个
	hour uint64
	// weekday 可取值 0-6 中的任意一个或多个
	weekday uint64
}

func NewWeekOfDayScheduler(minute, hour, weekday uint64) *WeekOfDayScheduler {
	return &WeekOfDayScheduler{
		minute:  minute,
		hour:    hour,
		weekday: weekday,
	}
}

func (s *WeekOfDayScheduler) Prev(baseTime time.Time) time.Time {
	// 前置检测
	if s.weekday == uint64(0) {
		panic("weekday can't be zero")
	}
	if s.hour == uint64(0) {
		panic("hour can't be zero")
	}
	if s.minute == uint64(0) {
		panic("minute can't be zero")
	}

	// << 如果当前星期对不上 >>
	if !isMatch(1<<uint64(baseTime.Weekday()), s.weekday) {
		// 直接找前一个星期的最晚的可执行时间
		t := s.findPrevWeekday(baseTime)
		return time.Date(t.Year(), t.Month(), t.Day(), s.findBiggestHour(), s.findBiggestMinute(), 0, 0, t.Location())
	}

	// << 如果当前小时匹配不上 >>
	if !isMatch(1<<uint64(baseTime.Hour()), s.hour) {
		// 找到了另一个可执行的小时
		if prevHour := s.findPrevHour(baseTime.Hour()); prevHour >= 0 {
			return time.Date(baseTime.Year(), baseTime.Month(), baseTime.Day(), prevHour, s.findBiggestMinute(), 0, 0, baseTime.Location())
		}
		// 如果未找到另一个可执行的小时，那就找前一个可执行的weekday吧
		t := s.findPrevWeekday(baseTime)
		return time.Date(t.Year(), t.Month(), t.Day(), s.findBiggestHour(), s.findBiggestMinute(), 0, 0, t.Location())
	}

	// << 如果当前分钟匹配不上 >>
	if !isMatch(1<<uint64(baseTime.Minute()), s.minute) {
		// 找另一个可执行的分钟
		if prevMinute := s.findPrevMinute(baseTime.Minute()); prevMinute >= 0 {
			return time.Date(baseTime.Year(), baseTime.Month(), baseTime.Day(), baseTime.Hour(), prevMinute, 0, 0, baseTime.Location())
		}
		// 找另一个可执行的小时的最晚的时间
		if prevHour := s.findPrevHour(baseTime.Hour()); prevHour >= 0 {
			return time.Date(baseTime.Year(), baseTime.Month(), baseTime.Day(), prevHour, s.findBiggestMinute(), 0, 0, baseTime.Location())
		}
		// 找另一个可执行的weekday的最早时间
		t := s.findPrevWeekday(baseTime)
		return time.Date(t.Year(), t.Month(), t.Day(), s.findBiggestHour(), s.findBiggestMinute(), 0, 0, baseTime.Location())
	}
	// << 我靠，所有的都能匹配的上 >>
	return baseTime
}

func (s *WeekOfDayScheduler) Next(baseTime time.Time) time.Time {
	// 前置检测
	if s.weekday == uint64(0) {
		panic("weekday can't be zero")
	}
	if s.hour == uint64(0) {
		panic("hour can't be zero")
	}
	if s.minute == uint64(0) {
		panic("minute can't be zero")
	}

	// << 如果当前星期对不上 >>
	if !isMatch(1<<uint64(baseTime.Weekday()), s.weekday) {
		// 直接找下一个星期的最早的可执行时间
		t := s.findNextWeekday(baseTime)
		return time.Date(t.Year(), t.Month(), t.Day(), s.findSmallestHour(), s.findSmallestMinute(), 0, 0, t.Location())
	}

	// << 如果当前小时匹配不上 >>
	if !isMatch(1<<uint64(baseTime.Hour()), s.hour) {
		// 找到了另一个可执行的小时
		if nextHour := s.findNextHour(baseTime.Hour()); nextHour >= 0 {
			return time.Date(baseTime.Year(), baseTime.Month(), baseTime.Day(), nextHour, s.findSmallestMinute(), 0, 0, baseTime.Location())
		}
		// 如果未找到另一个可执行的小时，那就找下一个可执行的weekday吧
		t := s.findNextWeekday(baseTime)
		return time.Date(t.Year(), t.Month(), t.Day(), s.findSmallestHour(), s.findSmallestMinute(), 0, 0, t.Location())
	}

	// << 如果当前分钟匹配不上 >>
	if !isMatch(1<<uint64(baseTime.Minute()), s.minute) {
		// 找另一个可执行的分钟
		if nextMinute := s.findNextMinute(baseTime.Minute()); nextMinute >= 0 {
			return time.Date(baseTime.Year(), baseTime.Month(), baseTime.Day(), baseTime.Hour(), nextMinute, 0, 0, baseTime.Location())
		}
		// 找另一个可执行的小时的最早的时间
		if nextHour := s.findNextHour(baseTime.Hour()); nextHour >= 0 {
			return time.Date(baseTime.Year(), baseTime.Month(), baseTime.Day(), nextHour, s.findSmallestMinute(), 0, 0, baseTime.Location())
		}
		// 找另一个可执行的weekday的最早时间
		t := s.findNextWeekday(baseTime)
		return time.Date(t.Year(), t.Month(), t.Day(), s.findSmallestHour(), s.findSmallestMinute(), 0, 0, t.Location())
	}
	// << 我靠，所有的都能匹配的上 >>
	return baseTime
}

func (s *WeekOfDayScheduler) findNextWeekday(baseTime time.Time) time.Time {
	t := baseTime
	for i := 1; i <= 7; i++ {
		t = t.Add(24 * time.Hour)
		if isMatch(1<<uint64(t.Weekday()), s.weekday) {
			return t
		}
	}
	panic(fmt.Sprintf("find next weekday failed, not found, weekday=%d", s.weekday))
	return time.Time{}
}

func (s *WeekOfDayScheduler) findPrevWeekday(baseTime time.Time) time.Time {
	t := baseTime
	for i := 1; i <= 7; i++ {
		t = t.Add(-24 * time.Hour)
		if isMatch(1<<uint64(t.Weekday()), s.weekday) {
			return t
		}
	}
	panic(fmt.Sprintf("find prev weekday failed, not found, weekday=%d", s.weekday))
	return time.Time{}
}

func (s *WeekOfDayScheduler) findSmallestHour() int {
	for i := 1; i <= 23; i++ {
		if (s.hour & (1 << uint64(i))) > 0 {
			return i
		}
	}
	panic(fmt.Sprintf("find smallest hour failed, not found, hour is %d", s.hour))
	return -1
}

func (s *WeekOfDayScheduler) findBiggestHour() int {
	for i := 23; i >= 1; i-- {
		if (s.hour & (1 << uint64(i))) > 0 {
			return i
		}
	}
	panic(fmt.Sprintf("find latest hour failed, not found, hour is %d", s.hour))
	return -1
}

func (s *WeekOfDayScheduler) findNextHour(baseHour int) int {
	for i := baseHour + 1; i <= 23; i++ {
		if (s.hour & (1 << uint64(i))) > 0 {
			return i
		}
	}
	return -1
}

func (s *WeekOfDayScheduler) findPrevHour(baseHour int) int {
	for i := baseHour - 1; i >= 0; i-- {
		if (s.hour & (1 << uint64(i))) > 0 {
			return i
		}
	}
	return -1
}

func (s *WeekOfDayScheduler) findSmallestMinute() int {
	for i := 1; i <= 59; i++ {
		if (s.minute & (1 << uint64(i))) > 0 {
			return i
		}
	}
	panic(fmt.Sprintf("find smallest minute failed, not found, minute is %d", s.minute))
	return -1
}

func (s *WeekOfDayScheduler) findBiggestMinute() int {
	for i := 59; i >= 0; i-- {
		if (s.minute & (1 << uint64(i))) > 0 {
			return i
		}
	}
	panic(fmt.Sprintf("find biggest minute failed, not found, minute is %d", s.minute))
	return -1
}

func (s *WeekOfDayScheduler) findNextMinute(baseMinute int) int {
	for i := baseMinute + 1; i <= 59; i++ {
		if (s.minute & (1 << uint64(i))) > 0 {
			return i
		}
	}
	return -1
}

func (s *WeekOfDayScheduler) findPrevMinute(baseMinute int) int {
	for i := baseMinute - 1; i >= 0; i-- {
		if (s.minute & (1 << uint64(i))) > 0 {
			return i
		}
	}
	return -1
}

func isMatch(v1, v2 uint64) bool {
	return v1&v2 != 0
}

//type DayOfMonthScheduler struct {
//	hour     uint64
//	monthday uint64
//}
//
