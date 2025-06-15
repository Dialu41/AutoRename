package service

import (
	"fyne.io/fyne/v2/widget"
)

type RenameData struct {
	RootPath  string
	NameIndex []*widget.Entry
}

func NewRenameData() *RenameData {
	return &RenameData{}
}
