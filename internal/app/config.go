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
)

type configKeySpec struct {
	valueType configValueType
	get       func(config.Config) string
	set       func(*config.Config, string) error
}

func allowedConfigKeys() map[string]configKeySpec {
	return map[string]configKeySpec{
		"provider":               {valueType: configTypeString, get: func(c config.Config) string { return c.Provider }, set: func(c *config.Config, v string) error { c.Provider = v; return nil }},
		"transcription.model":    {valueType: configTypeString, get: func(c config.Config) string { return c.Transcription.Model }, set: func(c *config.Config, v string) error { c.Transcription.Model = v; return nil }},
		"transcription.language": {valueType: configTypeString, get: func(c config.Config) string { return c.Transcription.Language }, set: func(c *config.Config, v string) error { c.Transcription.Language = v; return nil }},
		"cleanup.enabled": {valueType: configTypeBool, get: func(c config.Config) string { return strconv.FormatBool(c.Cleanup.Enabled) }, set: func(c *config.Config, v string) error {
			b, err := strconv.ParseBool(v)
			if err != nil {
				return fmt.Errorf("cleanup.enabled expects true/false")
			}
			c.Cleanup.Enabled = b
			return nil
		}},
		"cleanup.model": {valueType: configTypeString, get: func(c config.Config) string { return c.Cleanup.Model }, set: func(c *config.Config, v string) error { c.Cleanup.Model = v; return nil }},
		"output.mode":   {valueType: configTypeString, get: func(c config.Config) string { return c.Output.Mode }, set: func(c *config.Config, v string) error { c.Output.Mode = v; return nil }},
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
