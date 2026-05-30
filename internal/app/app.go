package app

import (
	"fmt"

	"github.com/spf13/cobra"

	"whisprgo/internal/config"
)

type App struct {
	root *cobra.Command
	cfg  config.Config
}

func New() (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	a := &App{cfg: cfg}
	a.root = a.newRootCommand()
	a.registerCommands()
	return a, nil
}

func (a *App) Execute() error {
	return a.root.Execute()
}

func (a *App) newRootCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "whisprgo",
		Short: "Minimal Linux CLI dictation tool",
	}
}

func (a *App) registerCommands() {
	a.root.AddCommand(a.newToggleCommand())
	a.root.AddCommand(a.newStatusCommand())
	a.root.AddCommand(a.newCancelCommand())
	a.root.AddCommand(a.newTranscribeCommand())
	a.root.AddCommand(a.newConfigCommand())
}
