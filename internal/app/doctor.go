package app

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"

	"whisprgo/internal/config"
	"whisprgo/internal/secrets"
)

func (a *App) newDoctorCommand() *cobra.Command {
	var strict bool

	cmd := &cobra.Command{
		Use:   "doctor",
		Short: "Check local runtime requirements",
		RunE: func(cmd *cobra.Command, args []string) error {
			issues := 0

			printCheck := func(name string, ok bool, detail string) {
				status := "OK"
				if !ok {
					status = "MISSING"
					issues++
				}
				fmt.Fprintf(cmd.OutOrStdout(), "%-18s %-8s %s\n", name, status, detail)
			}

			fmt.Fprintln(cmd.OutOrStdout(), "whisprgo doctor")
			fmt.Fprintln(cmd.OutOrStdout(), "----------------")

			_, ffmpegErr := exec.LookPath("ffmpeg")
			printCheck("ffmpeg", ffmpegErr == nil, "required for recording")

			_, xclipErr := exec.LookPath("xclip")
			_, xselErr := exec.LookPath("xsel")
			printCheck("clipboard", xclipErr == nil || xselErr == nil, "xclip or xsel")

			_, xdotoolErr := exec.LookPath("xdotool")
			printCheck("xdotool", xdotoolErr == nil, "optional paste")

			_, notifyErr := exec.LookPath("notify-send")
			printCheck("notify-send", notifyErr == nil, "desktop recording notifications")

			_, _, openAIErr := secrets.Get("openai")
			printCheck("OPENAI_API_KEY", openAIErr == nil, "env or keyring")

			cfgPath := config.Path()
			if _, err := os.Stat(cfgPath); err == nil {
				printCheck("config.yaml", true, cfgPath)
			} else {
				printCheck("config.yaml", false, cfgPath+" (optional, defaults are used)")
			}

			fmt.Fprintf(cmd.OutOrStdout(), "\nissues: %d\n", issues)
			if strict && issues > 0 {
				return fmt.Errorf("doctor found %d issue(s)", issues)
			}
			return nil
		},
	}

	cmd.Flags().BoolVar(&strict, "strict", false, "exit with non-zero status when any issue is found")
	return cmd
}
