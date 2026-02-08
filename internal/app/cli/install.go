package app

import (
	"fmt"

	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install [package id]",
	Short: "Install a package",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]

		pkg := packages.Pkg{
			Name: id,
			Type: "",
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
