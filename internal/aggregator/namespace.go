package aggregator

import (
	"github.com/charmbracelet/bubbles/list"
)

type NamespaceAggregator struct {
	contexts          []list.Item
	kubeconfig        string
	aggregatorChannel chan NamespaceResult
}

type NamespaceResult struct {
	ContextName string
	Namespaces  []string
	Err         error
}
