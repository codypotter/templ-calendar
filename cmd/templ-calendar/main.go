package main

import (
	"fmt"
	"os"
	"path/filepath"

	tc "github.com/codypotter/templ-calendar"
	"github.com/spf13/cobra"
)

var components = map[string]string{
	"calendar":  "calendar/calendar.templ",
	"navigator": "calendar/navigator.templ",
	"jumper":    "calendar/jumper.templ",
}

var rootCmd = &cobra.Command{
	Use:   "templ-calendar",
	Short: "Generate templ calendar components",
}

var addCmd = &cobra.Command{
	Use:   "add <component> [destination]",
	Short: "Copy a component into your project",
	Long: `Copy a templ-calendar component into your project.
Run templ generate yourself after adding.`,
	Args:      cobra.RangeArgs(1, 2),
	ValidArgs: []string{"calendar", "navigator", "jumper"},
	RunE: func(cmd *cobra.Command, args []string) error {
		component := args[0]
		dest := "."
		if len(args) == 2 {
			dest = args[1]
		}

		src, ok := components[component]
		if !ok {
			return fmt.Errorf("unknown component %q — run 'templ-calendar list' to see available components", component)
		}

		data, err := tc.Files.ReadFile(src)
		if err != nil {
			return fmt.Errorf("error reading component: %w", err)
		}

		if err := os.MkdirAll(dest, 0755); err != nil {
			return fmt.Errorf("error creating destination: %w", err)
		}

		destFile := filepath.Join(dest, filepath.Base(src))
		if err := os.WriteFile(destFile, data, 0644); err != nil {
			return fmt.Errorf("error writing file: %w", err)
		}

		fmt.Printf("added %s → %s\n", component, destFile)
		return nil
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available components",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Available components:")
		fmt.Println("  calendar   monthly calendar grid")
		fmt.Println("  navigator  prev/next month buttons")
		fmt.Println("  jumper     month/year jump form")
	},
}

func main() {
	rootCmd.AddCommand(addCmd, listCmd)
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
