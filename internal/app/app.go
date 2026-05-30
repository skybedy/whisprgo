package app

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"whisprgo/internal/config"
	"whisprgo/internal/logx"
	"whisprgo/internal/notifier"
	"whisprgo/internal/recorder"
)

type App struct {
	root     *cobra.Command
	cfg      config.Config
	recorder recorder.Recorder
	notifier *notifier.Notifier
	logger   *logx.Logger
	verbose  bool
}

func New() (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	a := &App{
		cfg:      cfg,
		recorder: recorder.NewFFmpegRecorder(),
		notifier: notifier.New(),
	}
	a.root = a.newRootCommand()
	a.registerCommands()
	return a, nil
}

func (a *App) Execute() error {
	return a.root.Execute()
}

func (a *App) newRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "whisprgo",
		Short: "Minimal Linux CLI dictation tool",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			a.ensureLogger(cmd.ErrOrStderr())
		},
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	cmd.PersistentFlags().BoolVar(&a.verbose, "verbose", false, "print diagnostic logs to stderr")
	return cmd
}

func (a *App) registerCommands() {
	a.root.AddCommand(a.newToggleCommand())
	a.root.AddCommand(a.newStatusCommand())
	a.root.AddCommand(a.newCancelCommand())
	a.root.AddCommand(a.newTranscribeCommand())
	a.root.AddCommand(a.newConfigCommand())
	a.root.AddCommand(a.newVersionCommand())
	a.root.AddCommand(a.newDoctorCommand())
}

func (a *App) ensureLogger(stderr io.Writer) {
	if a.logger != nil {
		return
	}

	logger, err := logx.New()
	if err != nil {
		if a.verbose {
			fmt.Fprintf(stderr, "warning: failed to initialize logger: %v\n", err)
		}
		return
	}
	a.logger = logger
}

func (a *App) logInfof(stderr io.Writer, format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	if a.logger != nil {
		a.logger.Info(msg)
	}
	if a.verbose && stderr != nil {
		fmt.Fprintf(stderr, "info: %s\n", msg)
	}
}

func (a *App) logErrorf(stderr io.Writer, format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	if a.logger != nil {
		a.logger.Error(msg)
	}
	if a.verbose && stderr != nil {
		fmt.Fprintf(stderr, "error: %s\n", msg)
	}
}

func (a *App) notify(summary string, body string, icon string, stderr io.Writer) {
	if a.notifier == nil {
		return
	}
	if err := a.notifier.Notify(summary, body, icon); err != nil {
		a.logErrorf(stderr, "notification failed: %v", err)
	}
}
