package mycrons

import (
	"fmt"
	"strconv"
	"strings"
)

// Parser 用于解析表达式
type Parser interface {
	Parse(spec string) (Scheduler, error)
}

// parser 是对 Parser 接口的实现
type parser struct{}

func NewParser() *parser {
	return &parser{}
}

// Parse 解析表达式返回对应的调度器
//
// Week类型的表达式格式:
// <week> [minute] [hour] [weekday]
// 前缀    [0,59]   [0,23] [1,7]或者[1-7]
// 例子：
// (1) week 0,30   10,20   1,3,5
//
// Month类型的表达式格式:
// <month> [minute] [hour] [monthday]
// 前缀     [0,59]   [0,23] [1,18]或者[1-15]
//
func (p *parser) Parse(spec string) (Scheduler, error) {
	fields := strings.Split(spec, " ")
	if len(fields) != 4 {
		return nil, fmt.Errorf("bad spec format, fields count should be 4")
	}
	switch fields[0] {
	case "week":
		return p.parseWeek(fields[1], fields[2], fields[3])
	case "month":
		return p.parseMonth(fields[1], fields[2], fields[3])
	}
	return nil, fmt.Errorf("bad spec format, unknown prefix `%s` was found", fields[0])
}

func (p *parser) parseWeek(minute, hour, weekday string) (Scheduler, error) {
	var m uint64 = 0
	for _, v := range strings.Split(minute, ",") {
		intV, err := parseToInt(v)
		if err != nil {
			return nil, fmt.Errorf("bad minute field format, %v", err)
		}
		if intV < 0 || intV > 59 {
			return nil, fmt.Errorf("bad minute value `%d`, it should be in range [0,59]", intV)
		}
		m |= 1 << uint64(intV)
	}

	var h uint64 = 0
	for _, v := range strings.Split(hour, ",") {
		intV, err := parseToInt(v)
		if err != nil {
			return nil, fmt.Errorf("bad hour field format, %v", err)
		}
		if intV < 0 || intV > 23 {
			return nil, fmt.Errorf("bad hour value `%d`, it should be in range [0,23]", intV)
		}
		h |= 1 << uint64(intV)
	}

	var d uint64 = 0
	if strings.Contains(weekday, "-") {
		// 使用范围的方式解析
		fields := strings.Split(weekday, "-")
		if len(fields) != 2 {
			return nil, fmt.Errorf("bad weekday field format, too many '-' were found")
		}
		fromDay, err := parseToInt(fields[0])
		if err != nil {
			return nil, fmt.Errorf("bad weekday field format, %v", err)
		}
		toDay, err := parseToInt(fields[1])
		if err != nil {
			return nil, fmt.Errorf("bad weekday field format, %v", err)
		}
		if fromDay < 1 || fromDay > 7 || toDay < 1 || toDay > 7 {
			return nil, fmt.Errorf("bad weekday field, day value should be in range [1,7]")
		}
		if toDay < fromDay {
			return nil, fmt.Errorf("bad weekday field, day range `[%d-%d]` is invalid", fromDay, toDay)
		}
		for ; fromDay <= toDay; fromDay++ {
			d |= 1 << uint64(fromDay)
		}
	} else {
		// 使用逐个指定的方式解析
		for _, v := range strings.Split(weekday, ",") {
			intV, err := parseToInt(v)
			if err != nil {
				return nil, fmt.Errorf("bad weekday field format, %v", err)
			}
			if intV < 1 || intV > 7 {
				return nil, fmt.Errorf("bad weekday value `%d`, it should be in range [1,7]", intV)
			}
			d |= 1 << uint64(intV)
		}
	}

	return newWeekScheduler(m, h, d), nil
}

func (p *parser) parseMonth(minute, hour, monthday string) (Scheduler, error) {
	var m uint64 = 0
	for _, v := range strings.Split(minute, ",") {
		intV, err := parseToInt(v)
		if err != nil {
			return nil, fmt.Errorf("bad minute field format, %v", err)
		}
		if intV < 0 || intV > 59 {
			return nil, fmt.Errorf("bad minute value `%d`, it should be in range [0,59]", intV)
		}
		m |= 1 << uint64(intV)
	}

	var h uint64 = 0
	for _, v := range strings.Split(hour, ",") {
		intV, err := parseToInt(v)
		if err != nil {
			return nil, fmt.Errorf("bad hour field format, %v", err)
		}
		if intV < 0 || intV > 23 {
			return nil, fmt.Errorf("bad hour value `%d`, it should be in range [0,23]", intV)
		}
		h |= 1 << uint64(intV)
	}

	var d uint64 = 0
	if strings.Contains(monthday, "-") {
		// 按照范围的方式进行解析
		fields := strings.Split(monthday, "-")
		if len(fields) != 2 {
			return nil, fmt.Errorf("bad monthday field format, too many '-' were found")
		}
		fromDay, err := parseToInt(fields[0])
		if err != nil {
			return nil, fmt.Errorf("bad monthday field format, %v", err)
		}
		toDay, err := parseToInt(fields[1])
		if err != nil {
			return nil, fmt.Errorf("bad monthday field format, %v", err)
		}
		if fromDay < 1 || fromDay > 31 || toDay < 1 || toDay > 31 {
			return nil, fmt.Errorf("bad monthday field, day value should be in range [1,31]")
		}
		if toDay < fromDay {
			return nil, fmt.Errorf("bad monthday field, day range `[%d-%d]` is invalid", fromDay, toDay)
		}
		for ; fromDay <= toDay; fromDay++ {
			d |= 1 << uint64(fromDay)
		}
	} else {
		// 默认按照逗号分隔的方式解析
		for _, v := range strings.Split(monthday, ",") {
			intV, err := parseToInt(v)
			if err != nil {
				return nil, fmt.Errorf("bad monthday field format, %v", err)
			}
			if intV < 1 || intV > 31 {
				return nil, fmt.Errorf("bad monthday value `%d`, it should be in range [1,31]", intV)
			}
			d |= 1 << uint64(intV)
		}
	}
	return newMonthScheduler(m, h, d), nil
}

func parseToInt(s string) (int64, error) {
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return v, nil
}
