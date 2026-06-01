package app

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/term"

	"whisprgo/internal/secrets"
)

func (a *App) newAuthCommand() *cobra.Command {
	authCmd := &cobra.Command{Use: "auth", Short: "Authentication helpers"}

	authCmd.AddCommand(&cobra.Command{
		Use:   "status",
		Short: "Show secrets configuration status",
		Run: func(cmd *cobra.Command, args []string) {
			transcriptionProvider := strings.TrimSpace(a.cfg.Transcription.Provider)
			if transcriptionProvider == "" {
				transcriptionProvider = strings.TrimSpace(a.cfg.Provider)
			}
			if transcriptionProvider == "" {
				transcriptionProvider = "openai"
			}

			cleanupProvider := strings.TrimSpace(a.cfg.Cleanup.Provider)
			if cleanupProvider == "" {
				cleanupProvider = strings.TrimSpace(a.cfg.Provider)
			}
			if cleanupProvider == "" {
				cleanupProvider = "openai"
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Transcription (%s): %s\n", transcriptionProvider, secrets.StatusForRole("transcription", transcriptionProvider))
			fmt.Fprintf(cmd.OutOrStdout(), "Cleanup (%s): %s\n", cleanupProvider, secrets.StatusForRole("cleanup", cleanupProvider))
		},
	})

	authCmd.AddCommand(&cobra.Command{
		Use:   "set <provider>",
		Short: "Set provider API key in keyring",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			provider := strings.ToLower(strings.TrimSpace(args[0]))
			if provider != "openai" && provider != "gemini" && provider != "mistral" && provider != "deepseek" && provider != "transcription" && provider != "cleanup" {
				return fmt.Errorf("supported targets: transcription, cleanup, openai, gemini, mistral, deepseek")
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Enter %s API key: ", strings.ToUpper(provider))
			key, err := readHiddenInput()
			fmt.Fprintln(cmd.OutOrStdout())
			if err != nil {
				return fmt.Errorf("failed to read key: %w", err)
			}
			if strings.TrimSpace(key) == "" {
				return fmt.Errorf("key cannot be empty")
			}
			if provider == "transcription" || provider == "cleanup" {
				if err := secrets.SetRole(provider, key); err != nil {
					return err
				}
			} else {
				if err := secrets.Set(provider, key); err != nil {
					return err
				}
			}
			fmt.Fprintf(cmd.OutOrStdout(), "%s key stored in keyring\n", strings.Title(provider))
			return nil
		},
	})

	authCmd.AddCommand(&cobra.Command{
		Use:   "delete <provider>",
		Short: "Delete provider API key from keyring",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			provider := strings.ToLower(strings.TrimSpace(args[0]))
			if provider != "openai" && provider != "gemini" && provider != "mistral" && provider != "deepseek" && provider != "transcription" && provider != "cleanup" {
				return fmt.Errorf("supported targets: transcription, cleanup, openai, gemini, mistral, deepseek")
			}
			if provider == "transcription" || provider == "cleanup" {
				if err := secrets.DeleteRole(provider); err != nil {
					return err
				}
			} else {
				if err := secrets.Delete(provider); err != nil {
					return err
				}
			}
			fmt.Fprintf(cmd.OutOrStdout(), "%s key deleted from keyring\n", strings.Title(provider))
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
