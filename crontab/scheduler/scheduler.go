package scheduler

import (
	"github.com/Hurricanezwf/toolbox/crontab/scheduler/src/basic"
	"github.com/Hurricanezwf/toolbox/crontab/scheduler/src/pro"
	"github.com/Hurricanezwf/toolbox/crontab/scheduler/types"
)

// GetParser 获取指定调度器的解析器
func GetParser(t types.ParserType) types.Parser {
	switch t {
	case types.ParserBasic:
		return basic.NewParser()
	case types.ParserPro:
		return pro.NewParser()
	}
	return nil
}
