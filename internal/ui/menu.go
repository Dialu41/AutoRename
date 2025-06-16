package ui

import (
	"AutoRename/constants"
	"AutoRename/internal/config"

	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// MakeMenu 创建并返回一个 Fyne 应用的主菜单。
// 参数 ap 是 Fyne 应用实例，用于调用应用相关的功能。
// 参数 win 是 Fyne 窗口实例，用于显示与窗口相关的对话框等。
// 参数 cfg 是用户配置的指针，可用于根据配置调整菜单内容，当前函数暂未使用该参数。
// 返回一个指向 fyne.MainMenu 的指针，代表创建好的主菜单。
func MakeMenu(ap fyne.App, win fyne.Window, cfg *config.UserConfig) *fyne.MainMenu {

	// 创建“使用说明”菜单项，点击该菜单项会调用 ShowHelp 函数打开使用说明页面
	helpItem := fyne.NewMenuItem("使用说明", func() {
		ShowHelp(ap)
	})
	// 创建“关于”菜单项，点击该菜单项会调用 showAbout 函数显示关于对话框
	aboutItem := fyne.NewMenuItem("关于", func() {
		showAbout(ap, win)
	})

	// 创建一个名为“选项”的菜单，并将“使用说明”和“关于”菜单项添加到该菜单中
	options := fyne.NewMenu("选项",
		helpItem,  // 使用说明菜单项
		aboutItem, // 关于菜单项
	)
	// 创建主菜单，并将“选项”菜单添加到主菜单中
	mainMenu := fyne.NewMainMenu(options)

	// 返回创建好的主菜单
	return mainMenu
}

func ShowHelp(ap fyne.App) {
	u, _ := url.Parse("https://www.huangoo.top/index.php/archives/169/")
	_ = ap.OpenURL(u)
}

// showAbout 显示关于对话框，包含软件版本号、版权信息、联系开发者方式、开源协议和鸣谢信息。
// 参数 ap 是 Fyne 应用实例，用于执行应用相关操作，如打开 URL、发送通知、操作剪贴板。
// 参数 win 是 Fyne 窗口实例，用于显示对话框。
func showAbout(ap fyne.App, win fyne.Window) {
	// 创建“我的邮箱”按钮，点击后将邮箱地址复制到剪贴板并发送通知
	contactButton := widget.NewButton("我的邮箱", func() {
		// 将邮箱地址复制到剪贴板
		ap.Clipboard().SetContent("1165011707@qq.com")
		ap.SendNotification(&fyne.Notification{
			Title:   "提示",
			Content: "已复制邮箱地址",
		})
	})
	// 创建“我的博客”按钮，点击后打开开发者博客页面
	blogButton := widget.NewButton("我的博客", func() {
		u, _ := url.Parse("https://www.huangoo.top")
		_ = ap.OpenURL(u)
	})
	// 创建“跳转仓库”按钮，点击后打开软件的 GitHub 仓库页面
	githubButton := widget.NewButton("跳转仓库", func() {
		u, _ := url.Parse(constants.GithubURL)
		_ = ap.OpenURL(u)
	})

	// 创建对话框内容，使用 2 列网格布局
	content := container.NewGridWithColumns(2,
		// 显示软件版本号
		widget.NewLabel("版本号"), widget.NewLabel(constants.AppVersion),
		// 显示版权信息
		widget.NewLabel("版权信息"), widget.NewLabel("Copyright © 2024 黄嚄嚄."),
		// 显示联系开发者方式，包含邮箱和博客按钮
		widget.NewLabel("联系开发者"), container.NewHBox(contactButton, blogButton),
		// 显示开源协议信息，包含协议说明和 GitHub 仓库按钮
		widget.NewLabel("开源协议"), container.NewHBox(widget.NewLabel("本软件使用MIT协议发行"), githubButton),
		// 显示鸣谢信息，包含对 Go 和 Fyne 的致谢按钮
		widget.NewLabel("鸣谢"), container.NewHBox(
			// 创建“Go”按钮，点击后打开 Go 语言的 GitHub 仓库
			widget.NewButton("Go", func() {
				u, _ := url.Parse(("https://github.com/golang/go"))
				_ = ap.OpenURL(u)
			}),
			// 创建“Fyne”按钮，点击后打开 Fyne 框架的 GitHub 仓库
			widget.NewButton("Fyne", func() {
				u, _ := url.Parse("https://github.com/fyne-io/fyne")
				_ = ap.OpenURL(u)
			}),
		),
	)
	// 显示自定义对话框，标题为“关于”，包含上述创建的内容
	dialog.ShowCustom("关于", "关闭", content, win)
}
