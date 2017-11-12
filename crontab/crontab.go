package crontab

import "errors"

var (
	crond   *Crond
	tasksrd map[string]struct{}
)

var (
	ErrNotOpen = errors.New("not open")
	ErrExist   = errors.New("task has been existed")
)

func Open() {
	if crond == nil {
		crond = &Crond{}
	}
	if tasksrd == nil {
		tasksrd = make(map[string]struct{})
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
	if err := crond.Add(t); err != nil {
		return err
	}
	tasksrd[t.TaskName] = struct{}{}

	return nil
}

func Del(taskName string) (affectedNum int) {
	if crond == nil {
		return 0
	}
	if _, ok := tasksrd[taskName]; !ok {
		return 0
	}
	return crond.Del(taskName)
}
