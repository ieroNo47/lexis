// lexis app
package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/log"
)

func main() {
	if err := os.Remove("debug.log"); err != nil && !os.IsNotExist(err) {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	// create a logger that writes to a file
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	//nolint:errcheck
	defer f.Close()
	logger := log.NewWithOptions(f, log.Options{
		ReportTimestamp: true,
		Level:           log.DebugLevel,
	})
	// create a new bubbletea program with our model
	p := tea.NewProgram(newModel(logger))
	logger.Info("==== Starting lexis ====")
	// run the program
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
