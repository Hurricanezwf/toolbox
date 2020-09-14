// +build unittest

package basic

import (
	"testing"
	"time"
)

func TestWeekScheduler_Prev(t *testing.T) {
	var (
		baseTime time.Time
		result   time.Time
		err      error
		loc      *time.Location
	)

	if loc, err = time.LoadLocation("Asia/Shanghai"); err != nil {
		t.Fatalf("load location failed, %v", err)
	}

	// week 匹配不上
	baseTime = time.Date(2020, 8, 12, 16, 2, 0, 0, loc)
	result = time.Date(2020, 8, 11, 22, 45, 0, 0, loc)
	if !newWeekScheduler(1<<30|1<<45, 1<<10|1<<22, 1<<2).Prev(baseTime).Equal(result) {
		t.Fatalf("failed")
	}

	// hour 匹配不上
	baseTime = time.Date(2020, 8, 11, 16, 2, 0, 0, loc)
	result = time.Date(2020, 8, 11, 10, 45, 0, 0, loc)
	if !newWeekScheduler(1<<30|1<<45, 1<<10|1<<22, 1<<2).Prev(baseTime).Equal(result) {
		t.Fatalf("failed")
	}
	baseTime = time.Date(2020, 8, 11, 23, 2, 0, 0, loc)
	result = time.Date(2020, 8, 11, 22, 45, 0, 0, loc)
	if !newWeekScheduler(1<<30|1<<45, 1<<10|1<<22, 1<<2).Prev(baseTime).Equal(result) {
		t.Fatalf("failed")
	}

	// minute 匹配不上
	baseTime = time.Date(2020, 8, 11, 10, 2, 0, 0, loc)
	result = time.Date(2020, 8, 4, 22, 45, 0, 0, loc)
	if !newWeekScheduler(1<<30|1<<45, 1<<10|1<<22, 1<<2).Prev(baseTime).Equal(result) {
		t.Fatalf("failed")
	}
	baseTime = time.Date(2020, 8, 11, 10, 32, 0, 0, loc)
	result = time.Date(2020, 8, 11, 10, 30, 0, 0, loc)
	if !newWeekScheduler(1<<30|1<<45, 1<<10|1<<22, 1<<2).Prev(baseTime).Equal(result) {
		t.Fatalf("failed")
	}
	baseTime = time.Date(2020, 8, 11, 10, 50, 0, 0, loc)
	result = time.Date(2020, 8, 11, 10, 45, 0, 0, loc)
	if !newWeekScheduler(1<<30|1<<45, 1<<10|1<<22, 1<<2).Prev(baseTime).Equal(result) {
		t.Fatalf("failed")
	}
	baseTime = time.Date(2020, 8, 11, 10, 20, 0, 0, loc)
	result = time.Date(2020, 8, 11, 10, 00, 0, 0, loc)
	if !newWeekScheduler(1<<0|1<<30, 1<<10|1<<22, 1<<2).Prev(baseTime).Equal(result) {
		t.Fatalf("failed")
	}

	// 都能匹配的上的
	baseTime = time.Date(2020, 8, 11, 22, 45, 0, 0, loc)
	result = time.Date(2020, 8, 11, 22, 30, 0, 0, loc)
	if !newWeekScheduler(1<<30|1<<45, 1<<10|1<<22, 1<<2).Prev(baseTime).Equal(result) {
		t.Fatalf("failed")
	}

	// 来个跨年的
	baseTime = time.Date(2020, 1, 1, 9, 2, 0, 0, loc)
	result = time.Date(2019, 12, 31, 22, 45, 0, 0, loc)
	if !newWeekScheduler(1<<30|1<<45, 1<<10|1<<22, 1<<2).Prev(baseTime).Equal(result) {
		t.Fatalf("failed")
	}
}

func TestWeekScheduler_Next(t *testing.T) {
	var (
		baseTime time.Time
		result   time.Time
		err      error
		loc      *time.Location
	)

	if loc, err = time.LoadLocation("Asia/Shanghai"); err != nil {
		t.Fatalf("load location failed, %v", err)
	}

	// week 匹配不上
	baseTime = time.Date(2020, 8, 12, 16, 2, 0, 0, loc)
	result = time.Date(2020, 8, 18, 10, 30, 0, 0, loc)
	if !newWeekScheduler(1<<30, 1<<10, 1<<2).Next(baseTime).Equal(result) {
		t.Fatalf("failed")
	}

	// hour匹配不上
	baseTime = time.Date(2020, 8, 11, 16, 2, 0, 0, loc)
	result = time.Date(2020, 8, 12, 10, 30, 0, 0, loc)
	if !newWeekScheduler(1<<30, 1<<10, 1<<2|1<<3).Next(baseTime).Equal(result) {
		t.Fatalf("failed")
	}
	baseTime = time.Date(2020, 8, 11, 7, 2, 0, 0, loc)
	result = time.Date(2020, 8, 11, 10, 30, 0, 0, loc)
	if !newWeekScheduler(1<<30, 1<<10, 1<<2).Next(baseTime).Equal(result) {
		t.Fatalf("failed")
	}
	baseTime = time.Date(2020, 8, 11, 8, 39, 0, 0, loc)
	result = time.Date(2020, 8, 11, 10, 30, 0, 0, loc)
	if !newWeekScheduler(1<<30, 1<<10, 1<<2).Next(baseTime).Equal(result) {
		t.Fatalf("failed")
	}
	baseTime = time.Date(2020, 8, 11, 18, 39, 0, 0, loc)
	result = time.Date(2020, 8, 18, 10, 30, 0, 0, loc)
	if !newWeekScheduler(1<<30, 1<<10, 1<<2).Next(baseTime).Equal(result) {
		t.Fatalf("failed")
	}

	// minute 匹配不上
	baseTime = time.Date(2020, 8, 11, 10, 2, 0, 0, loc)
	result = time.Date(2020, 8, 11, 10, 30, 0, 0, loc)
	if !newWeekScheduler(1<<30, 1<<10, 1<<2).Next(baseTime).Equal(result) {
		t.Fatalf("failed")
	}
	baseTime = time.Date(2020, 8, 11, 10, 39, 0, 0, loc)
	result = time.Date(2020, 8, 11, 18, 30, 0, 0, loc)
	if !newWeekScheduler(1<<30, 1<<10|1<<18, 1<<2).Next(baseTime).Equal(result) {
		t.Fatalf("failed")
	}
	baseTime = time.Date(2020, 8, 11, 10, 39, 0, 0, loc)
	result = time.Date(2020, 8, 11, 10, 45, 0, 0, loc)
	if !newWeekScheduler(1<<30|1<<45, 1<<10|1<<18, 1<<2).Next(baseTime).Equal(result) {
		t.Fatalf("failed")
	}
	baseTime = time.Date(2020, 8, 11, 10, 50, 0, 0, loc)
	result = time.Date(2020, 8, 11, 18, 30, 0, 0, loc)
	if !newWeekScheduler(1<<30|1<<45, 1<<10|1<<18, 1<<2).Next(baseTime).Equal(result) {
		t.Fatalf("failed")
	}
	baseTime = time.Date(2020, 8, 11, 10, 50, 0, 0, loc)
	result = time.Date(2020, 8, 11, 18, 00, 0, 0, loc)
	if !newWeekScheduler(1<<0|1<<30, 1<<10|1<<18, 1<<2).Next(baseTime).Equal(result) {
		t.Fatalf("failed")
	}

	// 都能匹配的上
	baseTime = time.Date(2020, 8, 11, 10, 45, 0, 0, loc)
	result = time.Date(2020, 8, 11, 18, 30, 0, 0, loc)
	if !newWeekScheduler(1<<30|1<<45, 1<<10|1<<18, 1<<2).Next(baseTime).Equal(result) {
		t.Fatalf("failed")
	}

	// 来个跨年的
	baseTime = time.Date(2019, 12, 31, 10, 48, 0, 0, loc)
	result = time.Date(2020, 1, 1, 10, 30, 0, 0, loc)
	if !newWeekScheduler(1<<30, 1<<10, 1<<3).Next(baseTime).Equal(result) {
		t.Fatalf("failed")
	}

	// extra
	baseTime = time.Date(2020, 9, 14, 11, 24, 23, 0, loc)
	result = time.Date(2020, 9, 16, 10, 0, 0, 0, loc)
	if !newWeekScheduler(1<<0, 1<<10, 1<<1|1<<3|1<<5).Next(baseTime).Equal(result) {
		t.Fatalf("failed")
	}
}

func TestMonthScheduler_Prev(t *testing.T) {
	var (
		baseTime time.Time
		result   time.Time
		err      error
		loc      *time.Location
	)

	if loc, err = time.LoadLocation("Asia/Shanghai"); err != nil {
		t.Fatalf("load location failed, %v", err)
	}

	// monthday 匹配不上
	baseTime = time.Date(2020, 8, 12, 16, 2, 0, 0, loc)
	result = time.Date(2020, 8, 2, 22, 45, 0, 0, loc)
	if !newMonthScheduler(1<<30|1<<45, 1<<10|1<<22, 1<<2).Prev(baseTime).Equal(result) {
		t.Fatalf("failed")
	}
	baseTime = time.Date(2020, 5, 12, 16, 2, 0, 0, loc)
	result = time.Date(2020, 3, 31, 22, 45, 0, 0, loc)
	if !newMonthScheduler(1<<30|1<<45, 1<<10|1<<22, 1<<31).Prev(baseTime).Equal(result) {
		t.Fatalf("failed")
	}
	baseTime = time.Date(2020, 5, 12, 16, 2, 0, 0, loc)
	result = time.Date(2020, 5, 10, 22, 45, 0, 0, loc)
	if !newMonthScheduler(1<<30|1<<45, 1<<10|1<<22, 1<<10|1<<31).Prev(baseTime).Equal(result) {
		t.Fatalf("failed")
	}

	// hour 匹配不上
	baseTime = time.Date(2020, 8, 2, 16, 2, 0, 0, loc)
	result = time.Date(2020, 8, 2, 10, 45, 0, 0, loc)
	if !newMonthScheduler(1<<30|1<<45, 1<<10|1<<22, 1<<2).Prev(baseTime).Equal(result) {
		t.Fatalf("failed")
	}
	baseTime = time.Date(2020, 8, 2, 23, 2, 0, 0, loc)
	result = time.Date(2020, 8, 2, 22, 45, 0, 0, loc)
	if !newMonthScheduler(1<<30|1<<45, 1<<10|1<<22, 1<<2).Prev(baseTime).Equal(result) {
		t.Fatalf("failed")
	}

	// minute 匹配不上
	baseTime = time.Date(2020, 8, 2, 10, 2, 0, 0, loc)
	result = time.Date(2020, 7, 2, 22, 45, 0, 0, loc)
	if !newMonthScheduler(1<<30|1<<45, 1<<10|1<<22, 1<<2).Prev(baseTime).Equal(result) {
		t.Fatalf("failed")
	}
	baseTime = time.Date(2020, 8, 2, 10, 32, 0, 0, loc)
	result = time.Date(2020, 8, 2, 10, 30, 0, 0, loc)
	if !newMonthScheduler(1<<30|1<<45, 1<<10|1<<22, 1<<2).Prev(baseTime).Equal(result) {
		t.Fatalf("failed")
	}
	baseTime = time.Date(2020, 8, 2, 10, 50, 0, 0, loc)
	result = time.Date(2020, 8, 2, 10, 45, 0, 0, loc)
	if !newMonthScheduler(1<<30|1<<45, 1<<10|1<<22, 1<<2).Prev(baseTime).Equal(result) {
		t.Fatalf("failed")
	}

	// 都能匹配的上
	baseTime = time.Date(2020, 8, 2, 10, 30, 0, 0, loc)
	result = time.Date(2020, 7, 2, 22, 45, 0, 0, loc)
	if !newMonthScheduler(1<<30|1<<45, 1<<10|1<<22, 1<<2).Prev(baseTime).Equal(result) {
		t.Fatalf("failed")
	}

	// 来个跨年的
	baseTime = time.Date(2020, 1, 1, 9, 2, 0, 0, loc)
	result = time.Date(2019, 12, 2, 22, 45, 0, 0, loc)
	if !newMonthScheduler(1<<30|1<<45, 1<<10|1<<22, 1<<2).Prev(baseTime).Equal(result) {
		t.Fatalf("failed")
	}
}

func TestMonthScheduler_Next(t *testing.T) {
	var (
		baseTime time.Time
		result   time.Time
		err      error
		loc      *time.Location
	)

	if loc, err = time.LoadLocation("Asia/Shanghai"); err != nil {
		t.Fatalf("load location failed, %v", err)
	}

	// monthday 匹配不上
	baseTime = time.Date(2020, 8, 12, 16, 2, 0, 0, loc)
	result = time.Date(2020, 9, 2, 10, 30, 0, 0, loc)
	if !newMonthScheduler(1<<30, 1<<10, 1<<2).Next(baseTime).Equal(result) {
		t.Fatalf("failed")
	}
	baseTime = time.Date(2020, 8, 12, 16, 2, 0, 0, loc)
	result = time.Date(2020, 8, 30, 10, 30, 0, 0, loc)
	if !newMonthScheduler(1<<30, 1<<10, 1<<2|1<<30).Next(baseTime).Equal(result) {
		t.Fatalf("failed")
	}
	baseTime = time.Date(2020, 6, 12, 16, 2, 0, 0, loc)
	result = time.Date(2020, 7, 31, 10, 30, 0, 0, loc)
	if !newMonthScheduler(1<<30, 1<<10, 1<<31).Next(baseTime).Equal(result) {
		t.Fatalf("failed")
	}

	// hour匹配不上
	baseTime = time.Date(2020, 8, 2, 16, 2, 0, 0, loc)
	result = time.Date(2020, 8, 3, 10, 30, 0, 0, loc)
	if !newMonthScheduler(1<<30, 1<<10, 1<<2|1<<3).Next(baseTime).Equal(result) {
		t.Fatalf("failed")
	}
	baseTime = time.Date(2020, 8, 2, 7, 2, 0, 0, loc)
	result = time.Date(2020, 8, 2, 10, 30, 0, 0, loc)
	if !newMonthScheduler(1<<30, 1<<10, 1<<2).Next(baseTime).Equal(result) {
		t.Fatalf("failed")
	}
	baseTime = time.Date(2020, 8, 2, 18, 39, 0, 0, loc)
	result = time.Date(2020, 9, 2, 10, 30, 0, 0, loc)
	if !newMonthScheduler(1<<30, 1<<10, 1<<2).Next(baseTime).Equal(result) {
		t.Fatalf("failed")
	}

	// minute 匹配不上
	baseTime = time.Date(2020, 8, 2, 10, 2, 0, 0, loc)
	result = time.Date(2020, 8, 2, 10, 30, 0, 0, loc)
	if !newMonthScheduler(1<<30, 1<<10, 1<<2).Next(baseTime).Equal(result) {
		t.Fatalf("failed")
	}
	baseTime = time.Date(2020, 8, 2, 10, 39, 0, 0, loc)
	result = time.Date(2020, 8, 2, 18, 30, 0, 0, loc)
	if !newMonthScheduler(1<<30, 1<<10|1<<18, 1<<2).Next(baseTime).Equal(result) {
		t.Fatalf("failed")
	}
	baseTime = time.Date(2020, 8, 2, 10, 39, 0, 0, loc)
	result = time.Date(2020, 8, 2, 10, 45, 0, 0, loc)
	if !newMonthScheduler(1<<30|1<<45, 1<<10|1<<18, 1<<2).Next(baseTime).Equal(result) {
		t.Fatalf("failed")
	}
	baseTime = time.Date(2020, 8, 2, 18, 50, 0, 0, loc)
	result = time.Date(2020, 9, 2, 10, 30, 0, 0, loc)
	if !newMonthScheduler(1<<30|1<<45, 1<<10|1<<18, 1<<2).Next(baseTime).Equal(result) {
		t.Fatalf("failed")
	}

	// 都能匹配的上
	baseTime = time.Date(2020, 8, 2, 18, 45, 0, 0, loc)
	result = time.Date(2020, 9, 2, 10, 30, 0, 0, loc)
	if !newMonthScheduler(1<<30|1<<45, 1<<10|1<<18, 1<<2).Next(baseTime).Equal(result) {
		t.Fatalf("failed")
	}

	// 来个跨年的
	baseTime = time.Date(2019, 12, 31, 10, 48, 0, 0, loc)
	result = time.Date(2020, 1, 3, 10, 30, 0, 0, loc)
	if !newMonthScheduler(1<<30, 1<<10, 1<<3).Next(baseTime).Equal(result) {
		t.Fatalf("failed")
	}
}
