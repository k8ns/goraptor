package container

import (
	"fmt"
	"sync"
)

type Container struct {
	instances map[string]interface{}
	once      map[string]*sync.Once
	mu        sync.Mutex
}

// NewServiceContainer creates a new container
func NewServiceContainer() *Container {
	return &Container{
		instances: make(map[string]interface{}),
		once:      make(map[string]*sync.Once),
	}
}

func (c *Container) Once(name string) *sync.Once {
	c.mu.Lock()
	defer c.mu.Unlock()

	once, exists := c.once[name]
	if !exists {
		once = &sync.Once{}
		c.once[name] = once
	}
	return once
}

func (c *Container) Register(name string, instance interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.instances[name] = instance
}

func (c *Container) Get(name string) interface{} {
	c.mu.Lock()
	defer c.mu.Unlock()

	instance, exists := c.instances[name]
	if !exists {
		panic(fmt.Sprintf("Service '%s' not registered", name))
	}

	return instance
}
