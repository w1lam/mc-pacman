package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/w1lam/mc-pacman/internal/app"
)

type appKey struct{}

func newRootCmd() *cobra.Command {
	var mcDir string

	cmd := &cobra.Command{
		Use:   "mcpac",
		Short: "Minecraft package manager",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			cfg := app.Config{
				McDir: mcDir,
			}

			view := New()
			a, err := app.New(view, cfg)
			if err != nil {
				return err
			}

			ctx := context.WithValue(cmd.Context(), appKey{}, a)
			cmd.SetContext(ctx)
			return nil
		},
	}

	cmd.PersistentFlags().StringVar(
		&mcDir,
		"mc-dir",
		"",
		"Custom Minecraft directory",
	)

	cmd.AddCommand(newGetCmd())
	cmd.AddCommand(newListCmd())
	cmd.AddCommand(newVerifyCmd())
	cmd.AddCommand(newInstallCmd())
	cmd.AddCommand(newUpdateCmd())
	cmd.AddCommand(newUpgradeCmd())

	return cmd
}

func Run() {
	cmd := newRootCmd()

	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
