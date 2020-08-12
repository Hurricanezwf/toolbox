// +build test

package mycrons

import (
	"testing"
	"time"
)

func TestDayOfWeekScheduler_Prev(t *testing.T) {
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
	if !NewWeekOfDayScheduler(1<<30|1<<45, 1<<10|1<<22, 1<<2).Prev(baseTime).Equal(result) {
		t.Fatalf("failed")
	}

	// hour 匹配不上
	baseTime = time.Date(2020, 8, 11, 16, 2, 0, 0, loc)
	result = time.Date(2020, 8, 11, 10, 45, 0, 0, loc)
	if !NewWeekOfDayScheduler(1<<30|1<<45, 1<<10|1<<22, 1<<2).Prev(baseTime).Equal(result) {
		t.Fatalf("failed")
	}
	baseTime = time.Date(2020, 8, 11, 23, 2, 0, 0, loc)
	result = time.Date(2020, 8, 11, 22, 45, 0, 0, loc)
	if !NewWeekOfDayScheduler(1<<30|1<<45, 1<<10|1<<22, 1<<2).Prev(baseTime).Equal(result) {
		t.Fatalf("failed")
	}

	// minute 匹配不上
	baseTime = time.Date(2020, 8, 11, 10, 2, 0, 0, loc)
	result = time.Date(2020, 8, 4, 22, 45, 0, 0, loc)
	if !NewWeekOfDayScheduler(1<<30|1<<45, 1<<10|1<<22, 1<<2).Prev(baseTime).Equal(result) {
		t.Fatalf("failed")
	}
	baseTime = time.Date(2020, 8, 11, 10, 32, 0, 0, loc)
	result = time.Date(2020, 8, 11, 10, 30, 0, 0, loc)
	if !NewWeekOfDayScheduler(1<<30|1<<45, 1<<10|1<<22, 1<<2).Prev(baseTime).Equal(result) {
		t.Fatalf("failed")
	}
	baseTime = time.Date(2020, 8, 11, 10, 50, 0, 0, loc)
	result = time.Date(2020, 8, 11, 10, 45, 0, 0, loc)
	if !NewWeekOfDayScheduler(1<<30|1<<45, 1<<10|1<<22, 1<<2).Prev(baseTime).Equal(result) {
		t.Fatalf("failed")
	}

	// 来个跨年的
	baseTime = time.Date(2020, 1, 1, 9, 2, 0, 0, loc)
	result = time.Date(2019, 12, 31, 22, 45, 0, 0, loc)
	if !NewWeekOfDayScheduler(1<<30|1<<45, 1<<10|1<<22, 1<<2).Prev(baseTime).Equal(result) {
		t.Fatalf("failed")
	}
}

func TestDayOfWeekScheduler_Next(t *testing.T) {
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
	if !NewWeekOfDayScheduler(1<<30, 1<<10, 1<<2).Next(baseTime).Equal(result) {
		t.Fatalf("failed")
	}

	// hour匹配不上
	baseTime = time.Date(2020, 8, 11, 16, 2, 0, 0, loc)
	result = time.Date(2020, 8, 12, 10, 30, 0, 0, loc)
	if !NewWeekOfDayScheduler(1<<30, 1<<10, 1<<2|1<<3).Next(baseTime).Equal(result) {
		t.Fatalf("failed")
	}
	baseTime = time.Date(2020, 8, 11, 7, 2, 0, 0, loc)
	result = time.Date(2020, 8, 11, 10, 30, 0, 0, loc)
	if !NewWeekOfDayScheduler(1<<30, 1<<10, 1<<2).Next(baseTime).Equal(result) {
		t.Fatalf("failed")
	}
	baseTime = time.Date(2020, 8, 11, 8, 39, 0, 0, loc)
	result = time.Date(2020, 8, 11, 10, 30, 0, 0, loc)
	if !NewWeekOfDayScheduler(1<<30, 1<<10, 1<<2).Next(baseTime).Equal(result) {
		t.Fatalf("failed")
	}
	baseTime = time.Date(2020, 8, 11, 18, 39, 0, 0, loc)
	result = time.Date(2020, 8, 18, 10, 30, 0, 0, loc)
	if !NewWeekOfDayScheduler(1<<30, 1<<10, 1<<2).Next(baseTime).Equal(result) {
		t.Fatalf("failed")
	}

	// minute 匹配不上
	baseTime = time.Date(2020, 8, 11, 10, 2, 0, 0, loc)
	result = time.Date(2020, 8, 11, 10, 30, 0, 0, loc)
	if !NewWeekOfDayScheduler(1<<30, 1<<10, 1<<2).Next(baseTime).Equal(result) {
		t.Fatalf("failed")
	}
	baseTime = time.Date(2020, 8, 11, 10, 39, 0, 0, loc)
	result = time.Date(2020, 8, 11, 18, 30, 0, 0, loc)
	if !NewWeekOfDayScheduler(1<<30, 1<<10|1<<18, 1<<2).Next(baseTime).Equal(result) {
		t.Fatalf("failed")
	}
	baseTime = time.Date(2020, 8, 11, 10, 39, 0, 0, loc)
	result = time.Date(2020, 8, 11, 10, 45, 0, 0, loc)
	if !NewWeekOfDayScheduler(1<<30|1<<45, 1<<10|1<<18, 1<<2).Next(baseTime).Equal(result) {
		t.Fatalf("failed")
	}
	baseTime = time.Date(2020, 8, 11, 10, 50, 0, 0, loc)
	result = time.Date(2020, 8, 11, 18, 30, 0, 0, loc)
	if !NewWeekOfDayScheduler(1<<30|1<<45, 1<<10|1<<18, 1<<2).Next(baseTime).Equal(result) {
		t.Fatalf("failed")
	}

	// 来个跨年的
	baseTime = time.Date(2019, 12, 31, 10, 48, 0, 0, loc)
	result = time.Date(2020, 1, 1, 10, 30, 0, 0, loc)
	if !NewWeekOfDayScheduler(1<<30, 1<<10, 1<<3).Next(baseTime).Equal(result) {
		t.Fatalf("failed")
	}
}
