package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/variableway/innate/capture/internal/config"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
}

var configGetCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "Get a config value",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := normalizeConfigKey(args[0])
		val, err := config.Get(getDataDir(), key)
		if err != nil {
			return err
		}
		fmt.Println(val)
		return nil
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a config value",
	Long: `Set a config value in config.yaml.

Windows (especially Nushell) notes:
- Backslashes inside double quotes are treated as escapes.
- Prefer single quotes, unquoted paths, or forward slashes.

Examples:
  capture config set workspace.root .
  capture config set workspace.root D:\innate-works\innate-works\innate-capture
  capture config set workspace.root 'D:\innate-works\innate-works\innate-capture'
  capture config set workspace.root D:/innate-works/innate-works/innate-capture`,
	Example: `  # If your shell reports "unrecognized escape after '\'", try:
  capture config set workspace.root 'D:\innate-works\innate-works\innate-capture'
  capture config set workspace.root D:/innate-works/innate-works/innate-capture`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := normalizeConfigKey(args[0])
		if err := config.Set(getDataDir(), key, args[1]); err != nil {
			return err
		}
		fmt.Printf("Set %s = %s\n", key, args[1])
		return nil
	},
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show effective full configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(getDataDir())
		if err != nil {
			return err
		}

		out, err := yaml.Marshal(cfg)
		if err != nil {
			return fmt.Errorf("marshal config: %w", err)
		}
		fmt.Print(string(out))
		return nil
	},
}

func init() {
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configShowCmd)
	rootCmd.AddCommand(configCmd)
}

func normalizeConfigKey(key string) string {
	switch key {
	case "data_dir", "app.data_dir":
		return "workspace.root"
	default:
		return key
	}
}
