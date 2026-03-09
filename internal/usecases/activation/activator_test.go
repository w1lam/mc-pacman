package activation

import (
	"path/filepath"
	"testing"

	"github.com/w1lam/mc-pacman/internal/app/paths"
	"github.com/w1lam/mc-pacman/internal/core/packages"
)

func TestEntryDestDir(t *testing.T) {
	p := paths.New("/root", "/minecraft")
	a := &Activator{paths: p}

	tests := []struct {
		entryType packages.EntryTypeID
		want      string
		wantErr   bool
	}{
		{packages.EntryTypeMod, filepath.Join("/minecraft", "mods"), false},
		{packages.EntryTypeResourcepack, filepath.Join("/minecraft", "resourcepacks"), false},
		{packages.EntryTypeShaderpack, filepath.Join("/minecraft", "shaderpacks"), false},
		{"unknown", "", true},
	}

	for _, tt := range tests {
		t.Run(string(tt.entryType), func(t *testing.T) {
			got, err := a.entryDestDir(tt.entryType)
			if (err != nil) != tt.wantErr {
				t.Errorf("wantErr %v, got %v", tt.wantErr, err)
			}
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}
