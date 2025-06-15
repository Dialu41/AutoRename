package ui

import (
	"AutoRename/internal/config"
	"AutoRename/internal/service"
	"AutoRename/mywidget"
	"errors"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

var (
	tabs         *container.AppTabs
	RootTab      *container.TabItem
	NameOrderTab *container.TabItem
)

var renameData *service.RenameData

func MakeTabs(ap fyne.App, win fyne.Window, cfg *config.UserConfig) *container.AppTabs {
	renameData = service.NewRenameData()

	RootTab = container.NewTabItem("工作文件夹", makeRootTabContent(win))
	NameOrderTab = container.NewTabItem("命名顺序", makeNameOrderTabContent(cfg))

	tabs = container.NewAppTabs(
		RootTab,
		NameOrderTab,
	)

	//选项卡靠左
	tabs.SetTabLocation(container.TabLocationLeading)

	return tabs
}

func makeRootTabContent(win fyne.Window) *fyne.Container {

	rootPathEntry := mywidget.NewFolderOpenWithEntry(func(s string) {
		renameData.RootPath = s
	}, "各子文件夹的根目录", win)

	//点击跳转下一个选项卡
	rootPathNextButton := widget.NewButton("下一步", func() {
		if rootPathEntry.GetValid() {
			tabs.Select(NameOrderTab)
		} else {
			dialog.ShowError(errors.New("工作文件夹不存在或未填写"), win)
		}
	})
	rootPathNextButton.Importance = widget.HighImportance

	return container.NewVBox(
		widget.NewForm(
			widget.NewFormItem("工作文件夹", rootPathEntry)),
		//保持跳转按钮靠下
		layout.NewSpacer(),
		//保持跳转按钮居中
		container.NewHBox(
			layout.NewSpacer(),
			rootPathNextButton,
			layout.NewSpacer(),
		),
	)
}

func makeNameOrderTabContent(cfg *config.UserConfig) *fyne.Container {

	nameContainer := container.NewVBox()

	for _, name := range *cfg {
		temp := widget.NewEntry()
		temp.SetText(name)
		renameData.NameIndex = append(renameData.NameIndex, temp)
		nameContainer.Add(renameData.NameIndex[len(renameData.NameIndex)-1])
	}

	nameOrderNextButton := widget.NewButton("开始监测", func() {

	})
	nameOrderNextButton.Importance = widget.DangerImportance

	nameOrderBackButton := widget.NewButton("上一步", func() {
		tabs.Select(RootTab)
	})

	addNameButton := widget.NewButton("添加命名", func() {
		renameData.NameIndex = append(renameData.NameIndex, widget.NewEntry())
		nameContainer.Add(renameData.NameIndex[len(renameData.NameIndex)-1])
	})

	deleteNameButton := widget.NewButton("删除命名", func() {
		length := len(renameData.NameIndex)
		if length > 0 {
			nameContainer.Remove(renameData.NameIndex[length-1])
			renameData.NameIndex = renameData.NameIndex[:length-1]
		}
	})

	return container.NewVBox(
		nameContainer,
		container.NewHBox(
			layout.NewSpacer(),
			addNameButton,
			deleteNameButton,
			layout.NewSpacer(),
		),
		layout.NewSpacer(),
		container.NewHBox(
			layout.NewSpacer(),
			nameOrderBackButton,
			nameOrderNextButton,
			layout.NewSpacer(),
		),
	)
}
