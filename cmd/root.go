package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	dataDir string
)

var rootCmd = &cobra.Command{
	Use:   "capture",
	Short: "Capture - innate-works 前端（idea inbox + daily）",
	Long: `Capture 是 innate-works 的终端前端。
快速记录灵感到 inbox，并查看今日 daily 清单。`,
	SilenceUsage: true,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.capture/config.yaml)")
	rootCmd.PersistentFlags().StringVar(&dataDir, "data-dir", "", "data directory (default is $HOME/.capture)")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		viper.AddConfigPath(home + "/.capture")
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv()
	viper.SetEnvPrefix("CAPTURE")

	if err := viper.ReadInConfig(); err == nil {
		used := viper.ConfigFileUsed()
		if cfgFile != "" {
			fmt.Fprintln(os.Stderr, "Using config file:", used)
			return
		}

		// Default boot config always starts from ~/.capture/config.yaml.
		// If it redirects data dir, make that explicit to avoid confusion.
		if d := resolveConfiguredDataDir(); d != "" {
			fmt.Fprintf(os.Stderr, "Using bootstrap config file: %s (effective data-dir: %s)\n", used, d)
			return
		}
		fmt.Fprintln(os.Stderr, "Using config file:", used)
	}
}

func getDataDir() string {
	if dataDir != "" {
		return dataDir
	}
	if d := resolveConfiguredDataDir(); d != "" {
		return d
	}
	home, _ := os.UserHomeDir()
	return home + "/.capture"
}

func resolveConfiguredDataDir() string {
	// Backward compatibility: legacy flat key first.
	if d := viper.GetString("data_dir"); d != "" {
		return d
	}
	// Current config domain schema.
	if d := viper.GetString("app.data_dir"); d != "" {
		return d
	}
	return ""
}
