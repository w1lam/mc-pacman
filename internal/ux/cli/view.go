package cli

import (
	ansi "github.com/w1lam/Packages/tui"
	"github.com/w1lam/mc-pacman/internal/core/events"
)

// TODO: FIX LISTER RENDERERING TO WORK WITH NEW PACKAGES INDEX STRUCTURE

type View struct {
	ansi bool
}

func NewCLIView() *View {
	return &View{
		ansi: ansi.SupportsANSI(),
	}
}

func (v *View) Emit(e events.Event) {
	switch e.Op.Scope {
	case events.ScopeDownloader:
		downloaderRenderer(e, v.ansi)
	case events.ScopeGetter:
		getterRenderer(e, v.ansi)
	case events.ScopeInstaller:
		installerRenderer(e, v.ansi)
	case events.ScopeUninstaller:
		uninstallRenderer(e, v.ansi)
	case events.ScopeUpdater:
		updaterRenderer(e, v.ansi)
	case events.ScopeResolver:
		resolverRenderer(e, v.ansi)
	case events.ScopeList:
		listRenderer(e, v.ansi)
	case events.ScopeVerifier:
		verifierRenderer(e, v.ansi)
	case events.ScopeBackup:
		backupRenderer(e, v.ansi)
	}
}
