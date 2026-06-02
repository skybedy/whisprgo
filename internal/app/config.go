package app

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"whisprgo/internal/config"
)

type configValueType string

const (
	configTypeString configValueType = "string"
	configTypeBool   configValueType = "bool"
	configTypeInt    configValueType = "int"
)

type configKeySpec struct {
	valueType configValueType
	get       func(config.Config) string
	set       func(*config.Config, string) error
}

func allowedConfigKeys() map[string]configKeySpec {
	return map[string]configKeySpec{
		"provider":               {valueType: configTypeString, get: func(c config.Config) string { return c.Provider }, set: func(c *config.Config, v string) error { c.Provider = v; return nil }},
		"transcription.provider": {valueType: configTypeString, get: func(c config.Config) string { return c.Transcription.Provider }, set: func(c *config.Config, v string) error { c.Transcription.Provider = v; return nil }},
		"transcription.model":    {valueType: configTypeString, get: func(c config.Config) string { return c.Transcription.Model }, set: func(c *config.Config, v string) error { c.Transcription.Model = v; return nil }},
		"transcription.language": {valueType: configTypeString, get: func(c config.Config) string { return c.Transcription.Language }, set: func(c *config.Config, v string) error { c.Transcription.Language = v; return nil }},
		"transcription.sherpa_ws_url": {valueType: configTypeString, get: func(c config.Config) string { return c.Transcription.SherpaWSURL }, set: func(c *config.Config, v string) error {
			c.Transcription.SherpaWSURL = v
			c.Transcription.Parakeet.SherpaWSURL = v
			return nil
		}},
		"transcription.parakeet.mode": {valueType: configTypeString, get: func(c config.Config) string { return c.Transcription.Parakeet.Mode }, set: func(c *config.Config, v string) error {
			c.Transcription.Parakeet.Mode = v
			return nil
		}},
		"transcription.parakeet.sherpa_ws_url": {valueType: configTypeString, get: func(c config.Config) string { return c.Transcription.Parakeet.SherpaWSURL }, set: func(c *config.Config, v string) error {
			c.Transcription.Parakeet.SherpaWSURL = v
			c.Transcription.SherpaWSURL = v
			return nil
		}},
		"transcription.parakeet.binary": {valueType: configTypeString, get: func(c config.Config) string { return c.Transcription.Parakeet.Binary }, set: func(c *config.Config, v string) error {
			c.Transcription.Parakeet.Binary = v
			return nil
		}},
		"transcription.parakeet.model_dir": {valueType: configTypeString, get: func(c config.Config) string { return c.Transcription.Parakeet.ModelDir }, set: func(c *config.Config, v string) error {
			c.Transcription.Parakeet.ModelDir = v
			return nil
		}},
		"transcription.parakeet.host": {valueType: configTypeString, get: func(c config.Config) string { return c.Transcription.Parakeet.Host }, set: func(c *config.Config, v string) error {
			c.Transcription.Parakeet.Host = v
			return nil
		}},
		"transcription.parakeet.port": {valueType: configTypeInt, get: func(c config.Config) string { return strconv.Itoa(c.Transcription.Parakeet.Port) }, set: func(c *config.Config, v string) error {
			n, err := strconv.Atoi(v)
			if err != nil {
				return fmt.Errorf("transcription.parakeet.port expects integer")
			}
			c.Transcription.Parakeet.Port = n
			return nil
		}},
		"transcription.parakeet.num_threads": {valueType: configTypeInt, get: func(c config.Config) string { return strconv.Itoa(c.Transcription.Parakeet.NumThreads) }, set: func(c *config.Config, v string) error {
			n, err := strconv.Atoi(v)
			if err != nil {
				return fmt.Errorf("transcription.parakeet.num_threads expects integer")
			}
			c.Transcription.Parakeet.NumThreads = n
			return nil
		}},
		"transcription.parakeet.startup_timeout_seconds": {valueType: configTypeInt, get: func(c config.Config) string { return strconv.Itoa(c.Transcription.Parakeet.StartupTimeoutSeconds) }, set: func(c *config.Config, v string) error {
			n, err := strconv.Atoi(v)
			if err != nil {
				return fmt.Errorf("transcription.parakeet.startup_timeout_seconds expects integer")
			}
			c.Transcription.Parakeet.StartupTimeoutSeconds = n
			return nil
		}},
		"transcription.parakeet.request_timeout_seconds": {valueType: configTypeInt, get: func(c config.Config) string { return strconv.Itoa(c.Transcription.Parakeet.RequestTimeoutSeconds) }, set: func(c *config.Config, v string) error {
			n, err := strconv.Atoi(v)
			if err != nil {
				return fmt.Errorf("transcription.parakeet.request_timeout_seconds expects integer")
			}
			c.Transcription.Parakeet.RequestTimeoutSeconds = n
			return nil
		}},
		"cleanup.enabled": {valueType: configTypeBool, get: func(c config.Config) string { return strconv.FormatBool(c.Cleanup.Enabled) }, set: func(c *config.Config, v string) error {
			b, err := strconv.ParseBool(v)
			if err != nil {
				return fmt.Errorf("cleanup.enabled expects true/false")
			}
			c.Cleanup.Enabled = b
			return nil
		}},
		"cleanup.model": {valueType: configTypeString, get: func(c config.Config) string { return c.Cleanup.Model }, set: func(c *config.Config, v string) error { c.Cleanup.Model = v; return nil }},
		"cleanup.provider": {valueType: configTypeString, get: func(c config.Config) string { return c.Cleanup.Provider }, set: func(c *config.Config, v string) error {
			c.Cleanup.Provider = v
			return nil
		}},
		"cleanup.prompt": {valueType: configTypeString, get: func(c config.Config) string { return c.Cleanup.Prompt }, set: func(c *config.Config, v string) error {
			c.Cleanup.Prompt = v
			return nil
		}},
		"output.mode": {valueType: configTypeString, get: func(c config.Config) string { return c.Output.Mode }, set: func(c *config.Config, v string) error { c.Output.Mode = v; return nil }},
		"output.file_path": {valueType: configTypeString, get: func(c config.Config) string { return c.Output.FilePath }, set: func(c *config.Config, v string) error {
			c.Output.FilePath = v
			return nil
		}},
		"output.copy_to_clipboard": {valueType: configTypeBool, get: func(c config.Config) string { return strconv.FormatBool(c.Output.CopyToClipboard) }, set: func(c *config.Config, v string) error {
			b, err := strconv.ParseBool(v)
			if err != nil {
				return fmt.Errorf("output.copy_to_clipboard expects true/false")
			}
			c.Output.CopyToClipboard = b
			return nil
		}},
		"output.paste_to_active_window": {valueType: configTypeBool, get: func(c config.Config) string { return strconv.FormatBool(c.Output.PasteToActiveWindow) }, set: func(c *config.Config, v string) error {
			b, err := strconv.ParseBool(v)
			if err != nil {
				return fmt.Errorf("output.paste_to_active_window expects true/false")
			}
			c.Output.PasteToActiveWindow = b
			return nil
		}},
		"logging.include_text": {valueType: configTypeBool, get: func(c config.Config) string { return strconv.FormatBool(c.Logging.IncludeText) }, set: func(c *config.Config, v string) error {
			b, err := strconv.ParseBool(v)
			if err != nil {
				return fmt.Errorf("logging.include_text expects true/false")
			}
			c.Logging.IncludeText = b
			return nil
		}},
		"logging.file_path": {valueType: configTypeString, get: func(c config.Config) string { return c.Logging.FilePath }, set: func(c *config.Config, v string) error {
			c.Logging.FilePath = v
			return nil
		}},
	}
}

func isSecretKeyPath(key string) bool {
	k := strings.ToLower(strings.TrimSpace(key))
	if strings.Contains(k, "api_key") || strings.Contains(k, "secrets") {
		return true
	}
	return k == "openai_api_key" || k == "gemini_api_key" || k == "api_key"
}

func (a *App) newConfigCommand() *cobra.Command {
	configCmd := &cobra.Command{Use: "config", Short: "Config helpers"}

	configCmd.AddCommand(&cobra.Command{
		Use:   "path",
		Short: "Print config path",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintln(cmd.OutOrStdout(), config.DisplayPath())
		},
	})

	configCmd.AddCommand(&cobra.Command{
		Use:   "show",
		Short: "Show current config",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}
			raw, err := yaml.Marshal(&cfg)
			if err != nil {
				return fmt.Errorf("failed to marshal config: %w", err)
			}
			fmt.Fprint(cmd.OutOrStdout(), string(raw))
			return nil
		},
	})

	configCmd.AddCommand(&cobra.Command{
		Use:   "get <key>",
		Short: "Get one config value",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			key := strings.TrimSpace(args[0])
			spec, ok := allowedConfigKeys()[key]
			if !ok {
				return fmt.Errorf("unknown key: %s", key)
			}
			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}
			fmt.Fprintln(cmd.OutOrStdout(), spec.get(cfg))
			return nil
		},
	})

	configCmd.AddCommand(&cobra.Command{
		Use:   "set <key> <value>",
		Short: "Set one config value",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			key := strings.TrimSpace(args[0])
			value := args[1]

			if isSecretKeyPath(key) {
				return errors.New("API keys are secrets. Use: whisprgo auth set openai")
			}
			spec, ok := allowedConfigKeys()[key]
			if !ok {
				return fmt.Errorf("unknown key: %s", key)
			}

			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}
			if spec.valueType == configTypeString {
				value = strings.TrimSpace(value)
				if value == "" {
					return fmt.Errorf("%s expects non-empty string", key)
				}
			}
			if err := spec.set(&cfg, value); err != nil {
				return err
			}
			if err := config.Save(cfg); err != nil {
				return fmt.Errorf("failed to save config: %w", err)
			}
			a.cfg = cfg
			fmt.Fprintf(cmd.OutOrStdout(), "updated %s\n", key)
			return nil
		},
	})

	return configCmd
}
