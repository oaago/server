package utils

import (
	"github.com/pterm/pterm"
)

func AppStartPrint() {
	pterm.DefaultCenter.Print(pterm.DefaultHeader.WithBackgroundStyle(pterm.NewStyle(pterm.BgLightBlue)).WithMargin(30).Sprint("Oaago Start"))
}

func AppDownPrint() {
	pterm.DefaultCenter.Print(pterm.DefaultHeader.WithBackgroundStyle(pterm.NewStyle(pterm.BgLightBlue)).WithMargin(30).Sprint("Oaago Server ShutDown"))
}
