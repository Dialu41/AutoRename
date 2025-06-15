package main

import (
	"AutoRename/internal/config"
	"AutoRename/internal/ui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	a := app.New()
	w := a.NewWindow("AutoRename")

	cfg := config.NewUserConfig()

	w.Resize(fyne.NewSize(800, 600))
	w.SetMaster()
	w.SetMainMenu(ui.MakeMenu(a, w, cfg))
	w.SetContent(ui.MakeTabs(a, w, cfg))
	w.CenterOnScreen()

	w.ShowAndRun()
}
