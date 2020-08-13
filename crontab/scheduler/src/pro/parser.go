package pro

import (
	"fmt"
	"strings"

	"github.com/Hurricanezwf/toolbox/crontab/scheduler/types"
)

type parser struct{}

func NewParser() *parser {
	return &parser{}
}

func (p *parser) Parse(spec string) (types.Scheduler, error) {
	if len(spec) > 0 && spec[0] == '@' {
		return p.parseSpec(spec)
	}
	// Split on whitespace.  We require 5 or 6 fields.
	// (second) (minute) (hour) (day of month) (month) (day of week, optional)
	fields := strings.Fields(spec)
	if len(fields) != 5 && len(fields) != 6 {
		return nil, fmt.Errorf("Expected 5 or 6 fields, found %d: %s", len(fields), spec)
	}

	// If a sixth field is not provided (DayOfWeek), then it is equivalent to star.
	if len(fields) == 5 {
		fields = append(fields, "*")
	}

	schedule := &Schedule{
		Second: getField(fields[0], seconds),
		Minute: getField(fields[1], minutes),
		Hour:   getField(fields[2], hours),
		Day:    getField(fields[3], days),
		Month:  getField(fields[4], months),
		Week:   getField(fields[5], weeks),
	}

	return schedule, nil
}

func (p *parser) parseSpec(spec string) (types.Scheduler, error) {
	switch spec {
	case "@yearly", "@annually":
		return &Schedule{
			Second: 1 << seconds.min,
			Minute: 1 << minutes.min,
			Hour:   1 << hours.min,
			Day:    1 << days.min,
			Month:  1 << months.min,
			Week:   all(weeks),
		}, nil

	case "@monthly":
		return &Schedule{
			Second: 1 << seconds.min,
			Minute: 1 << minutes.min,
			Hour:   1 << hours.min,
			Day:    1 << days.min,
			Month:  all(months),
			Week:   all(weeks),
		}, nil

	case "@weekly":
		return &Schedule{
			Second: 1 << seconds.min,
			Minute: 1 << minutes.min,
			Hour:   1 << hours.min,
			Day:    all(days),
			Month:  all(months),
			Week:   1 << weeks.min,
		}, nil

	case "@daily", "@midnight":
		return &Schedule{
			Second: 1 << seconds.min,
			Minute: 1 << minutes.min,
			Hour:   1 << hours.min,
			Day:    all(days),
			Month:  all(months),
			Week:   all(weeks),
		}, nil

	case "@hourly":
		return &Schedule{
			Second: 1 << seconds.min,
			Minute: 1 << minutes.min,
			Hour:   all(hours),
			Day:    all(days),
			Month:  all(months),
			Week:   all(weeks),
		}, nil
	}
	return nil, fmt.Errorf("Unrecognized descriptor: %s", spec)
}
