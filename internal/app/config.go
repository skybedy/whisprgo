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
		"provider.transcription": {
			valueType: configTypeString,
			get:       func(c config.Config) string { return c.Provider.Transcription },
			set: func(c *config.Config, value string) error {
				c.Provider.Transcription = value
				return nil
			},
		},
		"provider.transcription_model": {
			valueType: configTypeString,
			get:       func(c config.Config) string { return c.Provider.TranscriptionModel },
			set: func(c *config.Config, value string) error {
				c.Provider.TranscriptionModel = value
				return nil
			},
		},
		"cleanup.enabled": {
			valueType: configTypeBool,
			get:       func(c config.Config) string { return strconv.FormatBool(c.Cleanup.Enabled) },
			set: func(c *config.Config, value string) error {
				b, err := strconv.ParseBool(value)
				if err != nil {
					return fmt.Errorf("cleanup.enabled expects true/false")
				}
				c.Cleanup.Enabled = b
				return nil
			},
		},
		"cleanup.provider": {
			valueType: configTypeString,
			get:       func(c config.Config) string { return c.Cleanup.Provider },
			set: func(c *config.Config, value string) error {
				c.Cleanup.Provider = value
				return nil
			},
		},
		"cleanup.model": {
			valueType: configTypeString,
			get:       func(c config.Config) string { return c.Cleanup.Model },
			set: func(c *config.Config, value string) error {
				c.Cleanup.Model = value
				return nil
			},
		},
		"cleanup.prompt": {
			valueType: configTypeString,
			get:       func(c config.Config) string { return c.Cleanup.Prompt },
			set: func(c *config.Config, value string) error {
				c.Cleanup.Prompt = value
				return nil
			},
		},
		"output.clipboard": {
			valueType: configTypeBool,
			get:       func(c config.Config) string { return strconv.FormatBool(c.Output.Clipboard) },
			set: func(c *config.Config, value string) error {
				b, err := strconv.ParseBool(value)
				if err != nil {
					return fmt.Errorf("output.clipboard expects true/false")
				}
				c.Output.Clipboard = b
				return nil
			},
		},
		"output.paste": {
			valueType: configTypeBool,
			get:       func(c config.Config) string { return strconv.FormatBool(c.Output.Paste) },
			set: func(c *config.Config, value string) error {
				b, err := strconv.ParseBool(value)
				if err != nil {
					return fmt.Errorf("output.paste expects true/false")
				}
				c.Output.Paste = b
				return nil
			},
		},
		"output.paste_delay_ms": {
			valueType: configTypeInt,
			get:       func(c config.Config) string { return strconv.Itoa(c.Output.PasteDelayMs) },
			set: func(c *config.Config, value string) error {
				i, err := strconv.Atoi(value)
				if err != nil {
					return fmt.Errorf("output.paste_delay_ms expects integer")
				}
				c.Output.PasteDelayMs = i
				return nil
			},
		},
		"audio.input": {
			valueType: configTypeString,
			get:       func(c config.Config) string { return c.Audio.Input },
			set: func(c *config.Config, value string) error {
				c.Audio.Input = value
				return nil
			},
		},
		"audio.format": {
			valueType: configTypeString,
			get:       func(c config.Config) string { return c.Audio.Format },
			set: func(c *config.Config, value string) error {
				c.Audio.Format = value
				return nil
			},
		},
		"audio.recorder": {
			valueType: configTypeString,
			get:       func(c config.Config) string { return c.Audio.Recorder },
			set: func(c *config.Config, value string) error {
				c.Audio.Recorder = value
				return nil
			},
		},
	}
}

func (a *App) newConfigCommand() *cobra.Command {
	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Config helpers",
	}

	configCmd.AddCommand(&cobra.Command{
		Use:   "path",
		Short: "Print config path",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintln(cmd.OutOrStdout(), config.DisplayPath())
		},
	})

	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Create default config if missing",
		RunE: func(cmd *cobra.Command, args []string) error {
			force, _ := cmd.Flags().GetBool("force")
			if config.Exists() && !force {
				fmt.Fprintf(cmd.OutOrStdout(), "config exists at %s. overwrite? [y/N]: ", config.DisplayPath())
				var answer string
				if _, err := fmt.Fscanln(cmd.InOrStdin(), &answer); err != nil {
					return errors.New("aborted: config exists (use --force to overwrite)")
				}
				answer = strings.ToLower(strings.TrimSpace(answer))
				if answer != "y" && answer != "yes" {
					return errors.New("aborted: config exists")
				}
			}

			cfg := config.Default()
			if err := config.Save(cfg); err != nil {
				return fmt.Errorf("failed to write config: %w", err)
			}
			a.cfg = cfg
			fmt.Fprintf(cmd.OutOrStdout(), "config initialized at %s\n", config.DisplayPath())
			return nil
		},
	}
	initCmd.Flags().Bool("force", false, "overwrite existing config without prompt")
	configCmd.AddCommand(initCmd)

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

			if strings.Contains(strings.ToLower(key), "api_key") {
				return errors.New("API keys are not accepted in config; use auth commands or .env")
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
