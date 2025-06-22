package models

import "fmt"

type Context struct {
	Name    string `json:"name"`
	Context struct {
		Cluster string `json:"cluster"`
		User    string `json:"user"`
	} `json:"context"`
}

func (c *Context) Title() string {
	return fmt.Sprintf("Name: %s", c.Name)
}

func (c *Context) Description() string {
	return fmt.Sprintf("Cluster: %s", c.Context.Cluster)
}

func (c *Context) FilterValue() string {
	return c.Name
}
