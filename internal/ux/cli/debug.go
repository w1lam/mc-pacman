package cli

import (
	"fmt"
	"strings"

	"github.com/w1lam/mc-pacman/internal/core/events"
)

type DebugView struct{}

func NewDebugView() *DebugView {
	return &DebugView{}
}

func (v *DebugView) Emit(e events.Event) {
	if e.SubScope == "perFile" {
		return
	}

	sep := ""
	if e.Op.ParentID != "" {
		sep = " ⮡ "
	} else {
		fmt.Println()
	}
	fmt.Printf("%s %s[%s:%s] (%s) -> %s\n",
		e.Timestamp.Format("15:04:05"),
		sep,
		strings.ToUpper(string(e.Op.Scope)),
		e.Op.ID,
		strings.ToUpper(string(e.Type)),
		e.Op.Intent,
	)
	if e.Message != "" {
		fmt.Printf("          ⮡ #INFO# %s\n", e.Message)
	}
	if e.Error != nil {
		fmt.Printf("          ⮡ [ERROR]: %v\n", e.Error)
	}
	if e.Payload != nil {
		p := e.Payload
		fmt.Print("          ⮡ [PAYLOAD]\n")
		if p.Package != nil {
			fmt.Printf("          ⮡ [PKG] -> %s (%s) - %s\n", p.Package.Name, p.Package.ID, p.Package.Type)
			fmt.Printf("                     * PKGVER: %s MCVER: %s\n", p.Package.PkgVersion, p.Package.McVersion)
		}
		if p.Packages != nil {
			for _, pkg := range p.Packages {
				fmt.Printf("          ⮡ [PKG] -> %s (%s) - %s\n", pkg.Name, pkg.ID, pkg.Type)
				fmt.Printf("                     * PKGVER: %s MCVER: %s\n", pkg.PkgVersion, pkg.McVersion)
			}
		}
	}
}
