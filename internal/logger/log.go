package logger

import (
	"fmt"
	"io"
	"log/slog"

	tea "github.com/charmbracelet/bubbletea"
)

func Init(debug bool) (func(), error) {
	var loggerWritter io.Writer
	var cleanup func() = func() {}

	if debug {
		file, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Println("Error: ", err)
			return nil, err
		}
		loggerWritter = file
		cleanup = func() { file.Close() }
	} else {
		loggerWritter = io.Discard
	}

	logger := slog.New(slog.NewTextHandler(loggerWritter, &slog.HandlerOptions{}))
	slog.SetDefault(logger)
	return cleanup, nil
}
