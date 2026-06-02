package app

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"whisprgo/internal/config"
	"whisprgo/internal/control"
	"whisprgo/internal/parakeet"
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

			transcriptionProvider := strings.TrimSpace(a.cfg.Transcription.Provider)
			if transcriptionProvider == "" {
				transcriptionProvider = strings.TrimSpace(a.cfg.Provider)
			}
			switch transcriptionProvider {
			case "parakeet":
				mode := strings.ToLower(strings.TrimSpace(a.cfg.Transcription.Parakeet.Mode))
				if mode == "" {
					mode = "external"
				}
				printCheck("parakeet.mode", true, mode)
				switch mode {
				case "external":
					wsURL := strings.TrimSpace(a.cfg.Transcription.Parakeet.SherpaWSURL)
					if wsURL == "" {
						wsURL = strings.TrimSpace(a.cfg.Transcription.SherpaWSURL)
					}
					printCheck("sherpa_ws_url", wsURL != "", "websocket endpoint for local transcription")
					if wsURL != "" {
						err := parakeet.CheckExternalEndpoint(wsURL, 750*time.Millisecond)
						printCheck("sherpa endpoint", err == nil, "external mode expects already running server")
					}
				case "managed":
					err := parakeet.ValidateManagedConfig(a.cfg.Transcription.Parakeet)
					printCheck("parakeet managed", err == nil, "managed mode does not require manual server start")
					printCheck("parakeet binary", fileExists(a.cfg.Transcription.Parakeet.Binary), a.cfg.Transcription.Parakeet.Binary)
					printCheck("parakeet model_dir", dirExists(a.cfg.Transcription.Parakeet.ModelDir), a.cfg.Transcription.Parakeet.ModelDir)
					for _, name := range []string{"tokens.txt", "encoder.int8.onnx", "decoder.int8.onnx", "joiner.int8.onnx"} {
						path := ""
						if a.cfg.Transcription.Parakeet.ModelDir != "" {
							path = a.cfg.Transcription.Parakeet.ModelDir + "/" + name
						}
						printCheck(name, fileExists(path), path)
					}
				case "serve":
					err := parakeet.ValidateManagedConfig(a.cfg.Transcription.Parakeet)
					printCheck("parakeet serve", err == nil, "serve mode keeps local backend loaded for fast dictation")
					printCheck("control socket", true, control.SocketPath())
					printCheck("serve running", control.IsServeReachable(cmd.Context()), "start with: whisprgo serve")
					printCheck("parakeet binary", fileExists(a.cfg.Transcription.Parakeet.Binary), a.cfg.Transcription.Parakeet.Binary)
					printCheck("parakeet model_dir", dirExists(a.cfg.Transcription.Parakeet.ModelDir), a.cfg.Transcription.Parakeet.ModelDir)
				default:
					printCheck("parakeet.mode", false, "supported values: external, managed, serve")
				}
			default:
				_, _, providerErr := secrets.GetForRole("transcription", transcriptionProvider)
				printCheck(strings.ToUpper(transcriptionProvider)+"_API_KEY", providerErr == nil, "env or keyring")
			}

			cfgPath := config.Path()
			if _, err := os.Stat(cfgPath); err == nil {
				printCheck("config.yaml", true, cfgPath)
			} else {
				printCheck("config.yaml", false, cfgPath+" (optional, defaults are used)")
			}

			a.ensureLogger(cmd.ErrOrStderr())
			if a.logger != nil {
				printCheck("log file", true, a.logger.Path())
			} else {
				printCheck("log file", false, "logger unavailable")
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

func fileExists(path string) bool {
	if strings.TrimSpace(path) == "" {
		return false
	}
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func dirExists(path string) bool {
	if strings.TrimSpace(path) == "" {
		return false
	}
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}
