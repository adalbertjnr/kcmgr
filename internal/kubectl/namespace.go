package kubectl

func SetDefaultNamespace(namespace string) error {
	setNamespaceCommand := execCommand("kubectl", "config", "set-context", "--current", "--namespace", namespace)
	_, err := setNamespaceCommand.CombinedOutput()
	if err != nil {
		return err
	}

	return nil
}
