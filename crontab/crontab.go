package crontab

import (
	"context"
	"errors"
	"sync"
)

var (
	ErrExist    = errors.New("task has been existed")
	ErrNotFound = errors.New("not found")
	ErrClosed   = errors.New("crontab had been closed")
)

// Crontab 抽象出了定时任务接口
type Crontab interface {
	// Close 关闭crontab
	Close()

	// Add 添加任务
	// 如果同名任务已经存在，则返回 ErrExist
	Add(t *Task) error

	// Del 根据任务名删除
	// 如果任务不存在则返回 ErrNotFound
	Del(taskName string) (affected int, err error)

	// List 查询所有任务
	List() ([]string, error)
}

// crontab 实现了 Crontab 接口
type crontab struct {
	once   sync.Once
	ctx    context.Context
	cancel context.CancelFunc

	mutex   sync.RWMutex
	crond   *Crond
	tasksrd map[string]*Task
}

func New() *crontab {
	cron := &crontab{}
	cron.tasksrd = make(map[string]*Task)
	cron.ctx, cron.cancel = context.WithCancel(context.Background())
	cron.crond = NewCrond(cron.ctx)
	go cron.crond.Run()
	return cron
}

func (c *crontab) Close() {
	c.once.Do(func() {
		if c.cancel != nil {
			c.cancel()
		}
	})
}

func (c *crontab) Add(t *Task) error {
	if c.checkClosed() {
		return ErrClosed
	}
	if t == nil {
		return errors.New("nil task")
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	if _, ok := c.tasksrd[t.TaskName]; ok {
		return ErrExist
	}
	if err := c.crond.Add(t, true); err != nil {
		return err
	}
	c.tasksrd[t.TaskName] = t

	return nil
}

func (c *crontab) Del(taskName string) (affected int, err error) {
	if c.checkClosed() {
		return 0, ErrClosed
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	if _, ok := c.tasksrd[taskName]; !ok {
		return 0, ErrNotFound
	}
	affected = c.crond.Del(taskName)
	if affected > 0 {
		delete(c.tasksrd, taskName)
	}
	return affected, nil
}

func (c *crontab) List() ([]string, error) {
	if c.checkClosed() {
		return nil, ErrClosed
	}

	c.mutex.RLock()
	defer c.mutex.RUnlock()

	res := make([]string, 0, len(c.tasksrd))
	for _, t := range c.tasksrd {
		res = append(res, t.String())
	}
	return res, nil
}

func (c *crontab) checkClosed() bool {
	select {
	case <-c.ctx.Done():
		return true
	default:
	}
	return false
}
