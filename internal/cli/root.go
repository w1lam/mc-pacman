package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/w1lam/mc-pacman/internal/app"
)

var rootCmd = &cobra.Command{
	Use:   "mcpac",
	Short: "Minecraft package manager",
	Long:  "mcpac manages Minecraft modpacks, resource packs, and shaders.",
}

func Execute(app *app.App) {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
