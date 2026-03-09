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

const (
	name    = "mc-pacman"
	version = "0.1"
)

// UserAgent returns the user agent string for http requests
func UserAgent() string {
	return name + "-v" + version
}
