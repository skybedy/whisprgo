package app

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/term"

	"whispergo/internal/config"
)

func (a *App) newAuthCommand() *cobra.Command {
	authCmd := &cobra.Command{
		Use:   "auth",
		Short: "Authentication helpers",
	}

	authCmd.AddCommand(&cobra.Command{
		Use:   "status",
		Short: "Show OPENAI_API_KEY availability source",
		Run: func(cmd *cobra.Command, args []string) {
			_, source := config.ResolveSecretWithSource("OPENAI_API_KEY")
			fmt.Fprintf(cmd.OutOrStdout(), "OPENAI_API_KEY: %s\n", source)
		},
	})

	authCmd.AddCommand(&cobra.Command{
		Use:   "set-openai-key",
		Short: "Set OPENAI_API_KEY in ~/.config/whispergo/.env",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprint(cmd.OutOrStdout(), "Enter OPENAI_API_KEY: ")
			key, err := readHiddenInput()
			fmt.Fprintln(cmd.OutOrStdout())
			if err != nil {
				return fmt.Errorf("failed to read key: %w", err)
			}
			key = strings.TrimSpace(key)
			if key == "" {
				return fmt.Errorf("key cannot be empty")
			}
			if err := config.UpsertKeyInEnvFile(config.EnvPath(), "OPENAI_API_KEY", key); err != nil {
				return fmt.Errorf("failed to save key: %w", err)
			}
			fmt.Fprintln(cmd.OutOrStdout(), "OPENAI_API_KEY saved to ~/.config/whispergo/.env")
			return nil
		},
	})

	authCmd.AddCommand(&cobra.Command{
		Use:   "clear-openai-key",
		Short: "Remove OPENAI_API_KEY from ~/.config/whispergo/.env",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := config.RemoveKeyFromEnvFile(config.EnvPath(), "OPENAI_API_KEY"); err != nil {
				return fmt.Errorf("failed to clear key: %w", err)
			}
			fmt.Fprintln(cmd.OutOrStdout(), "OPENAI_API_KEY removed from ~/.config/whispergo/.env")
			return nil
		},
	})

	return authCmd
}

func readHiddenInput() (string, error) {
	fd := int(os.Stdin.Fd())
	if term.IsTerminal(fd) {
		b, err := term.ReadPassword(fd)
		if err != nil {
			return "", err
		}
		return string(b), nil
	}

	var v string
	if _, err := fmt.Fscanln(os.Stdin, &v); err != nil {
		return "", err
	}
	return v, nil
}
