// Package cli holds the CLI view/renderers/commands
package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/w1lam/mc-pacman/internal/app"
	"github.com/w1lam/mc-pacman/internal/core/packages"
	"github.com/w1lam/mc-pacman/internal/core/state"
)

func newRootCmd(a *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mcpac",
		Short: "Minecraft package manager",
	}

	// add commands
	for _, subCmd := range commands(a) {
		cmd.AddCommand(subCmd)
	}

	return cmd
}

func Run(a *app.App) {
	cmd := newRootCmd(a)

	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func commands(a *app.App) []*cobra.Command {
	return []*cobra.Command{
		// COMMAND: list lists installed packages
		{
			Use:   "list [package-id]",
			Short: "List installed packages",
			Args:  cobra.MaximumNArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				if len(args) == 1 {
					return a.Lister.ListPkg(cmd.Context(), packages.PkgID(args[0]))
				}
				return a.Lister.ListAll(cmd.Context())
			},
		},

		// COMMAND: search seasrches for remote packages
		{
			Use:   "search [package-id]",
			Short: "Search for available packages",
			Args:  cobra.MaximumNArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				if len(args) == 1 {
					return a.Lister.SearchPkg(cmd.Context(), packages.PkgID(args[0]))
				}
				return a.Lister.SearchAll(cmd.Context())
			},
		},

		// COMMAND: verify verifies packages
		{
			Use:   "verify",
			Short: "Verify installed packages",
			RunE: func(cmd *cobra.Command, args []string) error {
				fmt.Println("PLACEHOLDER: VERIFIER NEEDS REFACTOR")
				return nil
			},
		},

		// COMMAND: get gets given package
		{
			Use:   "get [package-id]",
			Short: "Gets a package",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				id := args[0]

				return a.Getter.Get(cmd.Context(), id)
			},
		},

		// COMMAND: install installs given package
		{
			Use:   "install [package-id]",
			Short: "Installs a package",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				id := args[0]

				return a.Installer.Install(cmd.Context(), id)
			},
		},

		// COMMAND: update updates given package
		{
			Use:   "update [package-id]",
			Short: "Updates",
			RunE: func(cmd *cobra.Command, args []string) error {
				id := args[0]
				fmt.Println("PLACEHOLDER Update:", id)

				return nil
			},
		},

		// COMMAND: upgrade upgrades specified package
		{
			Use:   "upgrade [package-id]",
			Short: "Upgrades",
			RunE: func(cmd *cobra.Command, args []string) error {
				id := args[0]
				if id == "" {
					id = "none"
				}
				fmt.Println("PLACEHOLDER Upgrade:", id)

				return nil
			},
		},

		// COMMAND: mcdir sets/shows minecraft directory
		{
			Use:   "mcdir [path]",
			Short: "Show or set Minecraft directory",
			RunE: func(cmd *cobra.Command, args []string) error {
				if len(args) == 0 {
					current := a.Paths.McDir()
					if current == "" {
						fmt.Println("No Minecraft directory configured")
						fmt.Println("Set with: mcpac mcdir /path/to/.minecraft")
					} else {
						fmt.Println("Current Minecraft directory:", current)
					}
					return nil
				}

				newPath := args[0]
				if _, err := os.Stat(newPath); err != nil {
					return fmt.Errorf("directory does not exist: %s", newPath)
				}

				return a.StateRepo.Update(func(s *state.State) error {
					s.McDir = newPath
					return nil
				})
			},
		},
	}
}
