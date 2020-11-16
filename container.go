package depcontainer

import (
	"errors"
	"fmt"
	"log"
)

type DepsMap map[string]interface{}
type Builder func(deps DepsMap) (interface{}, error)

type ComponentDef struct {
	Name        string
	Builder     Builder
	DependsOn   []string
	IsSingleton bool
}

func (c *ComponentDef) verify() error {
	if c.Name == "" {
		return errors.New("name must be set")
	}
	if c.Builder == nil {
		return errors.New("builder function must be set")
	}

	return nil
}

type Container struct {
	definitions map[string]*ComponentDef
	instances   map[string]interface{}
}

func NewContainer() *Container {
	return &Container{
		definitions: make(map[string]*ComponentDef),
		instances:   make(map[string]interface{}),
	}
}

func (c *Container) Add(def *ComponentDef) error {
	if err := def.verify(); err != nil {
		return err
	}

	c.definitions[def.Name] = def

	return nil
}

func (c *Container) MustAdd(def *ComponentDef) *Container {
	if err := c.Add(def); err != nil {
		log.Fatal(err)
	}

	return c
}

func (c *Container) Get(name string) (interface{}, error) {
	seen := make(map[string]struct{})
	return c.get(name, seen)
}

func (c *Container) get(name string, seen map[string]struct{}) (interface{}, error) {
	if _, found := seen[name]; found {
		return nil, fmt.Errorf("circular dependency for dependency %s", name)
	}

	if component, ok := c.instances[name]; ok {
		return component, nil
	}

	def, ok := c.definitions[name]
	if !ok {
		return nil, fmt.Errorf("no definition found for name %s", name)
	}
	seen[name] = struct{}{}
	dependencies := make(DepsMap)
	for _, depName := range def.DependsOn {
		dep, err := c.get(depName, seen)
		if err != nil {
			return nil, err
		}
		dependencies[depName] = dep
	}

	component, err := def.Builder(dependencies)
	if err != nil {
		return nil, err
	}
	if def.IsSingleton {
		c.instances[def.Name] = component
	}

	return component, nil
}

func (c *Container) MustGet(name string) interface{} {
	component, err := c.Get(name)
	if err != nil {
		log.Fatal(err)
	}

	return component
}
