package cli

import (
	"github.com/spf13/cobra"
	"github.com/w1lam/mc-pacman/internal/app"
	"github.com/w1lam/mc-pacman/internal/core/packages"
	"github.com/w1lam/mc-pacman/internal/services/remote"
)

func newInstallCmd(a *app.App) *cobra.Command {
	return &cobra.Command{
		Use:   "install [package-id]",
		Short: "Install a package",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]

			pkg, err := remote.GetSpecificPackage(packages.PkgID(id))
			if err != nil {
				return err
			}
			return a.Services.Installer.Install(cmd.Context(), pkg)
		},
	}
}
