package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/w1lam/mc-pacman/internal/services/remote"
)

func newListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List available packages",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Fetching packages...")

			remotePackageIndex, err := remote.GetAllAvailablePackages()
			if err != nil {
				return err
			}

			CliPackageListRenderer(remotePackageIndex)

			return nil
		},
	}
}
