package kubectl

import "fmt"

type Namespace struct {
	Name string
}

func (c *Namespace) Title() string {
	return fmt.Sprintf("Namespace: %s", c.Name)
}

func (c *Namespace) FilterValue() string {
	return c.Name
}

func (c *Namespace) Description() string {
	return c.Name
}

func SetDefaultNamespace(namespace string) error {
	setNamespaceCommand := execCommand("kubectl", "config", "set-context", "--current", "--namespace", namespace)
	_, err := setNamespaceCommand.CombinedOutput()
	if err != nil {
		return err
	}

	return nil
}
