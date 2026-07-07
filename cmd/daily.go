package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"

	"github.com/variableway/innate/capture/internal/config"
	"github.com/variableway/innate/capture/internal/daily"
	"github.com/variableway/innate/capture/internal/workspace"
)

var (
	dailyOpen    bool
	dailySection string
	dailyReset   bool
)

var dailyCmd = &cobra.Command{
	Use:   "daily",
	Short: "Show today's daily checklist from innate-works",
	Long: `Print the contents of innate-works daily/today.md.

If today.md does not exist yet, it is bootstrapped from
daily/_template/day.md with __DATE__ replaced by today's date.

Use --section to print only one focus section, or --open to launch
defaults.editor on today.md.

Use --reset to regenerate today.md from daily/_template/day.md
before printing or opening. This is automation-friendly and
non-interactive.`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(getDataDir())
		if err != nil {
			return err
		}

		if dailyReset {
			if err := daily.BootstrapFromTemplate(cfg); err != nil {
				return err
			}
		}

		if dailyOpen {
			return openDaily(cfg)
		}

		if dailySection != "" {
			section, err := daily.PrintSection(cfg, dailySection)
			if err != nil {
				return err
			}
			fmt.Println(section)
			return nil
		}

		content, err := daily.Read(cfg)
		if err != nil {
			return err
		}
		fmt.Print(content)
		if !strings.HasSuffix(content, "\n") {
			fmt.Println()
		}
		return nil
	},
}

// openDaily ensures today.md exists (bootstrapping if needed) and opens
// it in the configured editor. Checkbox state is left to the user.
func openDaily(cfg *config.Config) error {
	if err := workspace.Validate(cfg); err != nil {
		return err
	}

	path := workspace.DailyPath(cfg)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := daily.BootstrapFromTemplate(cfg); err != nil {
			return err
		}
	}

	editor := cfg.Defaults.Editor
	if editor == "" {
		return fmt.Errorf("defaults.editor is not configured; set it in config.yaml")
	}

	c := exec.Command(editor, path)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	if err := c.Run(); err != nil {
		return fmt.Errorf("launch editor %q: %w", editor, err)
	}
	return nil
}

func init() {
	dailyCmd.Flags().BoolVar(&dailyOpen, "open", false, "Open today.md in defaults.editor")
	dailyCmd.Flags().StringVar(&dailySection, "section", "", "Print only a section (input|output|ideas)")
	dailyCmd.Flags().BoolVar(&dailyReset, "reset", false, "Regenerate today.md from template before reading")
	rootCmd.AddCommand(dailyCmd)
}
