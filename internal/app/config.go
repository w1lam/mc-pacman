package app

import (
	"github.com/w1lam/mc-pacman/internal/app/paths"
	"github.com/w1lam/mc-pacman/internal/core/state"
)

type Config struct {
	McDir string
}

func resolveMincraftDir(st *state.State, custom string) string {
	if custom != "" {
		return custom
	}

	if st.McDir != "" {
		return st.McDir
	}

	return paths.DefaultMinecraftDir()
}
