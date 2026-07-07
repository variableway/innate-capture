package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/variableway/innate/capture/internal/config"
	"github.com/variableway/innate/capture/internal/idea"
)

var (
	ideaDesc    string
	ideaContext string
)

var ideaCmd = &cobra.Command{
	Use:   "idea",
	Short: "Capture ideas to innate-works inbox",
}

var ideaAddCmd = &cobra.Command{
	Use:   "add <title>",
	Short: "Add an idea to inbox",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		title := strings.Join(args, " ")
		cfg, err := config.Load(getDataDir())
		if err != nil {
			return err
		}

		path, err := idea.Write(cfg, title, ideaDesc, ideaContext, idea.SourceCLI)
		if err != nil {
			return err
		}

		fmt.Printf("Created: %s\n", path)
		return nil
	},
}

var ideaListCmd = &cobra.Command{
	Use:   "list",
	Short: "List ideas in inbox",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(getDataDir())
		if err != nil {
			return err
		}

		entries, err := idea.List(cfg)
		if err != nil {
			return err
		}

		if len(entries) == 0 {
			fmt.Println("No ideas in inbox.")
			return nil
		}

		fmt.Printf("%-12s %-30s %s\n", "DATE", "FILE", "TITLE")
		fmt.Println(strings.Repeat("-", 80))
		for _, e := range entries {
			fmt.Printf("%-12s %-30s %s\n", e.Date, e.Filename, truncate(e.Title, 40))
		}
		fmt.Printf("\nTotal: %d idea(s)\n", len(entries))
		return nil
	},
}

func init() {
	ideaAddCmd.Flags().StringVarP(&ideaDesc, "description", "d", "", "One-line description")
	ideaAddCmd.Flags().StringVarP(&ideaContext, "context", "c", "", "Optional raw context")
	ideaCmd.AddCommand(ideaAddCmd)
	ideaCmd.AddCommand(ideaListCmd)
	rootCmd.AddCommand(ideaCmd)
}

func truncate(s string, max int) string {
	if max <= 0 || len(s) <= max {
		return s
	}
	if max <= 3 {
		return s[:max]
	}
	return s[:max-3] + "..."
}
