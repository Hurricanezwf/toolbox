package crontab

import "errors"

var (
	crond   *Crond
	tasksrd map[string]*Task
)

var (
	ErrNotOpen = errors.New("not open")
	ErrExist   = errors.New("task has been existed")
)

func Open() {
	if crond == nil {
		crond = NewCrond()
	}
	if tasksrd == nil {
		tasksrd = make(map[string]*Task)
	}

	go crond.Run()
}

func Close() {
	crond.Close()
	crond = nil
	tasksrd = nil
}

func Add(t *Task) error {
	if crond == nil {
		return ErrNotOpen
	}
	if t == nil {
		return errors.New("nil task")
	}

	if _, ok := tasksrd[t.TaskName]; ok {
		return ErrExist
	}
	if err := crond.Add(t, true); err != nil {
		return err
	}
	tasksrd[t.TaskName] = t

	return nil
}

func Del(taskName string) (affected int) {
	if crond == nil {
		return 0
	}
	if _, ok := tasksrd[taskName]; !ok {
		return 0
	}
	affected = crond.Del(taskName)
	if affected > 0 {
		delete(tasksrd, taskName)
	}
	return affected
}

func List() []string {
	res := make([]string, 0, len(tasksrd))
	for _, t := range tasksrd {
		res = append(res, t.String())
	}
	return res
}
