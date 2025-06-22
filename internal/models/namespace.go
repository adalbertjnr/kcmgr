package models

import "fmt"

type Namespace struct {
	Name string
	Age  string
}

func (c *Namespace) Title() string {
	return fmt.Sprintf("namespace: %s", c.Name)
}

func (c *Namespace) FilterValue() string {
	return c.Name
}

func (c *Namespace) Description() string {
	return fmt.Sprintf("created: %s", c.Age)
}
