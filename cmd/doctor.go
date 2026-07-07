package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/variableway/innate/capture/internal/config"
	"github.com/variableway/innate/capture/internal/workspace"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check workspace directory structure",
	Long:  "Validate that the innate-works workspace root and required directories (ideas/, daily/) are properly configured and accessible.",
	RunE: func(cmd *cobra.Command, args []string) error {
		dataDir := getDataDir()
		cfg, err := config.Load(dataDir)
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		fmt.Printf("workspace root (config): %s\n", cfg.Workspace.Root)

		wsRoot := workspace.Root(cfg)
		if envRoot := os.Getenv("CAPTURE_WORKSPACE_ROOT"); envRoot != "" {
			fmt.Printf("workspace root (env):   %s\n", envRoot)
		}
		fmt.Printf("workspace root (used):  %s\n", wsRoot)

		fmt.Printf("ideas inbox dir:        %s\n", workspace.InboxDir(cfg))
		fmt.Printf("daily today path:       %s\n", workspace.DailyPath(cfg))

		if err := workspace.Validate(cfg); err != nil {
			fmt.Printf("\n✗ %v\n", err)
			return nil
		}

		fmt.Println("\n✓ Workspace is valid.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}
