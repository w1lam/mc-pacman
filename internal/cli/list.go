package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/w1lam/mc-pacman/internal/services/repository"
	ux "github.com/w1lam/mc-pacman/internal/ux/cli"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available packages",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Fetching packages...")

		remotePackageIndex, err := repository.GetAllAvailablePackages()
		if err != nil {
			return err
		}

		ux.CliPackageListRenderer(remotePackageIndex)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
