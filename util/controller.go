package util

import (
	"sync"
)

type Controller interface {
	Add(int)
	Wait()
	Lock()
	Unlock()
	Done()
}

type controller struct {
	mu sync.Mutex
	wg sync.WaitGroup
}

func NewController() (c Controller) {
	c = &controller{}
	return
}

func (c *controller) Add(n int) {
	c.wg.Add(n)
}

func (c *controller) Done() {
	c.wg.Done()
}

func (c *controller) Wait() {
	c.wg.Wait()
}

func (c *controller) Lock() {
	c.mu.Lock()
}

func (c *controller) Unlock() {
	c.mu.Unlock()
}

func GoFunction[T any](c Controller, f func(T), params T) {
	defer c.Done()
	c.Lock()
	defer c.Unlock()
	f(params)
}
