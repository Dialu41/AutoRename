package main

import (
	"AutoRename/internal/config"
	"AutoRename/internal/ui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	// 创建一个新的 Fyne 应用实例
	a := app.New()
	// 为该应用创建一个新的窗口，窗口标题为 "AutoRename"
	w := a.NewWindow("AutoRename")

	// 创建一个新的用户配置实例
	cfg := config.NewUserConfig()
	// 从配置文件中读取用户配置信息
	cfg.ReadConfig(a)

	// 调整窗口大小
	w.Resize(fyne.NewSize(800, 600))
	// 将当前窗口设置为主窗口
	w.SetMaster()
	// 为窗口设置主菜单，菜单由 ui 包中的 MakeMenu 函数生成
	w.SetMainMenu(ui.MakeMenu(a, w, cfg))
	// 为窗口设置内容，内容由 ui 包中的 MakeTabs 函数生成
	w.SetContent(ui.MakeTabs(a, w, cfg))
	// 将窗口居中显示在屏幕上
	w.CenterOnScreen()

	// 显示窗口并开始运行应用的事件循环，阻塞直到窗口关闭
	w.ShowAndRun()
}
