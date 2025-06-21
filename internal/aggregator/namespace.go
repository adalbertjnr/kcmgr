package aggregator

import (
	"github.com/adalbertjnr/kcmgr/internal/client"
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

func New(aggregatorChannel chan NamespaceResult, kubeconfig string, contexts []list.Item) *NamespaceAggregator {
	return &NamespaceAggregator{
		kubeconfig:        kubeconfig,
		contexts:          contexts,
		aggregatorChannel: aggregatorChannel,
	}
}

func (n *NamespaceAggregator) Start() {
	for _, ctx := range n.contexts {
		ctxString := ctx.FilterValue()

		go func(ctxString string) {
			namespaces, err := client.GetNamespacesByContext(n.kubeconfig, ctxString)
			n.aggregatorChannel <- NamespaceResult{
				ContextName: ctxString,
				Namespaces:  namespaces,
				Err:         err,
			}
		}(ctxString)
	}
}
