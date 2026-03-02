package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/w1lam/mc-pacman/internal/app"
)

func getApp(cmd *cobra.Command) *app.App {
	v := cmd.Context().Value(appKey{})
	if v == nil {
		panic("app not found in context")
	}

	return v.(*app.App)
}

func newListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List available packages",
		RunE: func(cmd *cobra.Command, args []string) error {
			a := getApp(cmd)

			return a.UseCases.Lister.ListAllRemote(cmd.Context())
		},
	}
}

func newVerifyCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "verify",
		Short: "Verify installed packages",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("PLACEHOLDER: VERIFIER NEEDS REFACTOR")
			return nil
		},
	}
}

func newGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get [package-id]",
		Short: "Gets a package",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]
			a := getApp(cmd)

			return a.UseCases.Getter.Get(cmd.Context(), id)
		},
	}
}

func newInstallCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "install [package-id]",
		Short: "Installs a package",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]
			a := getApp(cmd)

			return a.UseCases.Installer.Install(cmd.Context(), id)
		},
	}
}

// TODO: IMPLEMENT UPDATE LOGIC
// newUpdateCmd PLACEHOLDER COMMAND
func newUpdateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "update [package-id]",
		Short: "Updates",
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]
			a := getApp(cmd)
			_ = a
			fmt.Println("PLACEHOLDER Update:", id)

			return nil
		},
	}
}

// TODO: IMPLEMENT UPGRADE LOGIC
func newUpgradeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "upgrade [package-id]",
		Short: "Upgrades",
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]
			if id == "" {
				id = "none"
			}
			a := getApp(cmd)
			_ = a
			fmt.Println("PLACEHOLDER Upgrade:", id)

			return nil
		},
	}
}
