package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/w1lam/mc-pacman/internal/core/packages"
	"github.com/w1lam/mc-pacman/internal/services/repository"
)

var installCmd = &cobra.Command{
	Use:   "install [package id]",
	Short: "Install a package",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]

		pkg, err := repository.GetSpecificPackage(packages.PkgID(id))
		if err != nil {
			return err
		}

		if err := actions.InstallPackage(pkg); err != nil {
			return err
		}

		fmt.Println("Installed:", id)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
