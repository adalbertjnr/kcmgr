package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/adalbertjnr/kcmgr/internal/bubble"
	"github.com/adalbertjnr/kcmgr/internal/kubectl"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "debug options")
	flag.Parse()

	if debug {
		file, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			log.Println("Error: ", err)
			return
		}
		defer file.Close()
	}

	currentContext, err := kubectl.CurrentContext()
	if err != nil {
		log.Println("Error getting current context: ", err)
		return
	}

	contextItems, err := kubectl.KubernetesContexts()
	if err != nil {
		log.Println("Error getting contexts: ", err)
		return
	}

	appModel := bubble.New(currentContext, contextItems)

	p := tea.NewProgram(appModel, tea.WithAltScreen())
	m, err := p.Run()
	if err != nil {
		fmt.Println("Error running program: ", err)
		return
	}

	if message, ok := m.(bubble.Model); ok && message.ContextMessage != "" {
		fmt.Println(message.ContextMessage)
	}

}
