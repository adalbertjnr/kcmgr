package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"

	"github.com/adalbertjnr/kcmgr/internal/bubble"
	"github.com/adalbertjnr/kcmgr/internal/client"
	"github.com/adalbertjnr/kcmgr/internal/kubectl"
	"github.com/adalbertjnr/kcmgr/internal/logger"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	contextsWindowTitle   = "Current + Available Contexts"
	namespacesWindowTitle = "Available Namespaces"
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "debug options")
	flag.Parse()

	cleanup, err := logger.Init(debug)
	if err != nil {
		log.Fatal("failed to initialize logger: ", err)
	}
	defer cleanup()

	currentContext, err := kubectl.CurrentContext()
	if err != nil {
		log.Fatal("Error getting current context: ", err)
	}

	contextItems, err := kubectl.KubernetesContexts()
	if err != nil {
		log.Fatal("Error getting contexts: ", err)
	}

	kubeconfig := client.GetKubeConfigFile()
	slog.Info("kubeconfig", "path", kubeconfig)

	model := bubble.New(
		contextsWindowTitle,
		namespacesWindowTitle,
		kubeconfig,
		currentContext,
		contextItems,
	)

	p := tea.NewProgram(model, tea.WithAltScreen())
	m, err := p.Run()
	if err != nil {
		log.Fatal("Error running program: ", err)
		return
	}

	if message, ok := m.(bubble.Model); ok && message.ContextMessage != "" {
		fmt.Println(message.ContextMessage)
	}

}
