package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/w1lam/mc-pacman/internal/app"
)

func NewRootCmd(a *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mcpac",
		Short: "Minecraft package manager",
	}

	cmd.AddCommand(newInstallCmd(a))
	cmd.AddCommand(newListCmd())

	return cmd
}

func Execute(a *app.App) {
	cmd := NewRootCmd(a)

	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
