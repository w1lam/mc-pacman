package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/w1lam/mc-pacman/internal/app"
	"github.com/w1lam/mc-pacman/internal/core/packages"
)

func NewRootCmd(a *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mcpac",
		Short: "Minecraft package manager",
	}

	cmd.AddCommand(newInstallCmd(a))
	cmd.AddCommand(newListCmd(a))
	cmd.AddCommand(newVerifyCmd(a))

	return cmd
}

func Run(a *app.App) {
	cmd := NewRootCmd(a)

	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func newInstallCmd(a *app.App) *cobra.Command {
	return &cobra.Command{
		Use:   "install [package-id]",
		Short: "Install a package",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]

			ctx := cmd.Context()

			pkg, err := a.Services.Installer.RemoteRepo.GetByID(ctx, packages.PkgID(id))
			if err != nil {
				return err
			}
			return a.Services.Installer.Install(cmd.Context(), pkg)
		},
	}
}

func newListCmd(a *app.App) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List available packages",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			return a.Services.Lister.ListAllRemote(ctx)
		},
	}
}

func newVerifyCmd(a *app.App) *cobra.Command {
	return &cobra.Command{
		Use:   "verify",
		Short: "Verify installed packages",
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.Services.Verifier.Verify()
		},
	}
}
