package ui

import (
	"AutoRename/internal/config"
	"AutoRename/internal/service"
	"AutoRename/mywidget"
	"errors"
	"os"
	"path/filepath"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

var (
	renameData *service.RenameData

	tabs         *container.AppTabs
	RootTab      *container.TabItem
	DirsTab      *container.TabItem
	NameOrderTab *container.TabItem

	rootPathEntry *mywidget.FolderOpenWithEntry
)

// MakeTabs 创建一个包含多个选项卡的应用选项卡组件。
// 参数 ap 是 Fyne 应用实例，用于与应用层面的功能交互。
// 参数 win 是 Fyne 窗口实例，用于处理窗口相关的操作。
// 参数 cfg 是用户配置的指针，包含用户自定义的配置信息。
// 返回一个指向 container.AppTabs 的指针。
func MakeTabs(ap fyne.App, win fyne.Window, cfg *config.UserConfig) *container.AppTabs {
	// 初始化重命名数据，调用 service 包中的 NewRenameData 函数
	renameData = service.NewRenameData()

	// 创建“工作文件夹”选项卡，调用 makeRootTabContent 函数生成选项卡内容
	RootTab = container.NewTabItem("工作文件夹", makeRootTabContent(win))
	DirsTab = container.NewTabItem("子文件夹", makeDirsTabContent())
	// 创建“命名顺序”选项卡，调用 makeNameOrderTabContent 函数生成选项卡内容
	NameOrderTab = container.NewTabItem("命名顺序", makeNameOrderTabContent(ap, win, cfg))

	// 创建应用选项卡组件，将“工作文件夹”和“命名顺序”选项卡添加进去
	tabs = container.NewAppTabs(
		RootTab,
		DirsTab,
		NameOrderTab,
	)

	// 设置选项卡的位置为左侧
	tabs.SetTabLocation(container.TabLocationLeading)

	// 返回创建好的应用选项卡组件
	return tabs
}

// makeRootTabContent 创建“工作文件夹”选项卡的内容容器。
// 参数 win 是 Fyne 窗口实例，用于处理窗口相关的操作，如显示对话框等。
// 返回一个指向 fyne.Container 的指针，包含“工作文件夹”选项卡的内容。
func makeRootTabContent(win fyne.Window) *fyne.Container {
	// 创建一个自定义的文件夹选择输入框组件，当选择文件夹后，将所选文件夹路径赋值给 renameData.RootPath
	// 第二个参数 "各子文件夹的根目录" 为输入框的提示信息，第三个参数 win 为关联的窗口
	rootPathEntry = mywidget.NewFolderOpenWithEntry(func(s string) {
		renameData.RootPath = s
	}, "各子文件夹的根目录", win)

	// 点击该按钮会跳转到下一个选项卡
	rootPathNextButton := widget.NewButton("下一步", func() {
		// 检查根目录输入框中的路径是否有效
		if rootPathEntry.GetValid() {
			// 若路径有效，切换到“命名顺序”选项卡
			tabs.Select(DirsTab)
		} else {
			// 若路径无效，显示错误对话框提示用户
			dialog.ShowError(errors.New("工作文件夹不存在或未填写"), win)
		}
	})
	// 设置按钮的重要性为高，通常会有不同的视觉样式
	rootPathNextButton.Importance = widget.HighImportance

	// 返回一个垂直布局的容器，包含表单、间隔器和按钮组件
	return container.NewVBox(
		// 创建一个表单，包含一个表单项，表单项的标签为“工作文件夹”，内容为根目录输入框
		widget.NewForm(
			widget.NewFormItem("工作文件夹", rootPathEntry)),
		// 添加间隔器，使跳转按钮靠下显示
		layout.NewSpacer(),
		// 创建一个水平布局的容器，使跳转按钮居中显示
		container.NewHBox(
			// 添加间隔器，使按钮居中
			layout.NewSpacer(),
			rootPathNextButton,
			layout.NewSpacer(),
		),
	)
}

// makeDirsTabContent 创建“子文件夹”选项卡的内容容器。
// 返回一个指向 fyne.Container 的指针，该容器包含添加、删除子文件夹输入框以及导航按钮等组件。
func makeDirsTabContent() *fyne.Container {
	// 创建一个垂直布局的容器，用于存放子文件夹输入框
	DirsContainer := container.NewVBox()

	// 创建“上一步”按钮，点击后切换到“工作文件夹”选项卡
	dirsBackButton := widget.NewButton("上一步", func() {
		tabs.Select(RootTab)
	})

	// 创建“下一步”按钮，点击后切换到“命名顺序”选项卡
	dirsNextButton := widget.NewButton("下一步", func() {
		tabs.Select(NameOrderTab)
	})

	// 创建“添加子文件夹”按钮，点击后添加一个新的输入框用于输入子文件夹信息
	// 新的输入框会被添加到 DirsContainer 中，同时更新 renameData 中的 DirsEntryIndex
	addDirButton := widget.NewButton("添加子文件夹", func() {
		renameData.DirsEntryIndex = append(renameData.DirsEntryIndex, widget.NewEntry())
		DirsContainer.Add(renameData.DirsEntryIndex[len(renameData.DirsEntryIndex)-1])
	})

	// 创建“删除子文件夹”按钮，点击后删除最后一个子文件夹输入框
	// 若存在输入框，则从 DirsContainer 中移除，并更新 renameData 中的 DirsEntryIndex
	deleteDirButton := widget.NewButton("删除子文件夹", func() {
		length := len(renameData.DirsEntryIndex)
		if length > 0 {
			DirsContainer.Remove(renameData.DirsEntryIndex[length-1])
			renameData.DirsEntryIndex = renameData.DirsEntryIndex[:length-1]
		}
	})

	// 返回一个垂直布局的容器，包含子文件夹输入框容器、添加删除按钮容器、间隔器以及导航按钮容器
	return container.NewVBox(
		DirsContainer,
		// 水平布局容器，包含添加和删除按钮，两侧使用间隔器使按钮居中
		container.NewHBox(
			layout.NewSpacer(),
			addDirButton,
			deleteDirButton,
			layout.NewSpacer(),
		),
		// 添加间隔器，使导航按钮靠下显示
		layout.NewSpacer(),
		// 水平布局容器，包含上一步和下一步按钮，两侧使用间隔器使按钮居中
		container.NewHBox(
			layout.NewSpacer(),
			dirsBackButton,
			dirsNextButton,
			layout.NewSpacer(),
		),
	)
}

// makeNameOrderTabContent 创建“命名顺序”选项卡的内容容器。
// 参数 ap 是 Fyne 应用实例，用于应用级别的操作，如保存配置文件。
// 参数 win 是 Fyne 窗口实例，用于显示对话框等窗口相关操作。
// 参数 cfg 是用户配置的指针，包含用户自定义的命名顺序配置。
// 返回一个指向 fyne.Container 的指针，该容器包含“命名顺序”选项卡的所有内容。
func makeNameOrderTabContent(ap fyne.App, win fyne.Window, cfg *config.UserConfig) *fyne.Container {
	// 创建一个垂直布局的容器，用于存放命名顺序输入框
	nameContainer := container.NewVBox()

	// 遍历用户配置中的命名顺序，为每个命名创建一个输入框
	for _, name := range *cfg {
		temp := widget.NewEntry()
		temp.SetText(name)
		renameData.NameEntryIndex = append(renameData.NameEntryIndex, temp)
		nameContainer.Add(renameData.NameEntryIndex[len(renameData.NameEntryIndex)-1])
	}

	// 创建“开始监测”按钮，点击后开始验证输入并执行重命名相关操作
	nameOrderNextButton := widget.NewButton("开始监测", func() {
		// 检查工作文件夹路径是否有效
		if !rootPathEntry.GetValid() {
			dialog.ShowError(errors.New("工作文件夹不存在或未填写"), win)
			return
		}

		// 遍历子文件夹输入框，检查是否有未填写的情况
		for _, dIndex := range renameData.DirsEntryIndex {
			if dIndex.Text == "" {
				dialog.ShowError(errors.New("子文件夹填写错误"), win)
				return
			}
		}
		// 遍历命名顺序输入框，检查是否有未填写的情况
		for _, nIndex := range renameData.NameEntryIndex {
			if nIndex.Text == "" {
				dialog.ShowError(errors.New("命名顺序填写错误"), win)
				return
			}
		}

		// 清空用户配置中的命名顺序
		*cfg = (*cfg)[:0]

		// 将当前命名顺序输入框中的文本添加到用户配置中
		for _, nIndex := range renameData.NameEntryIndex {
			*cfg = append(*cfg, nIndex.Text)
		}
		// 保存更新后的用户配置到文件
		cfg.SaveConfigFile(ap)

		// 清空工作文件夹下的所有文件和子文件夹
		files, _ := filepath.Glob(filepath.Join(renameData.RootPath, "*"))
		for _, file := range files {
			os.RemoveAll(file)
		}

		// 根据子文件夹输入框的内容，在工作文件夹下创建子文件夹
		for _, dIndex := range renameData.DirsEntryIndex {
			dirPath := filepath.Join(renameData.RootPath, dIndex.Text)
			os.MkdirAll(dirPath, 0755)
		}

		// 创建活动指示器，用于显示任务正在执行
		act := widget.NewActivity()
		text := widget.NewLabel("自动重命名中...")
		// 创建垂直布局的容器，包含活动指示器和状态标签
		content := container.NewVBox(
			// 水平布局，将活动指示器居中显示
			container.NewHBox(
				layout.NewSpacer(),
				act,
				layout.NewSpacer(),
			),
			// 水平布局，将状态标签居中显示
			container.NewHBox(
				layout.NewSpacer(),
				text,
				layout.NewSpacer(),
			),
		)

		// 创建自定义对话框，显示任务状态
		resultDialog := dialog.NewCustom("任务状态", "关闭", content, win)
		resultDialog.Resize(fyne.NewSize(200, 150))
		resultDialog.Show()

		// 启动活动指示器
		act.Start()

		// 创建一个 WaitGroup 用于等待所有文件夹监测任务完成
		var wg sync.WaitGroup
		numDirs := len(renameData.DirsEntryIndex)
		wg.Add(numDirs)
		// 为每个子文件夹启动一个 goroutine 进行监测
		for i := 0; i < numDirs; i++ {
			go func(dirID int) {
				defer wg.Done()
				// 调用 MonitorDirs 方法监测指定子文件夹的变化并进行重命名
				renameData.MonitorDirs(renameData.DirsEntryIndex[dirID].Text)
			}(i)
		}

		// 启动一个 goroutine 等待所有重命名任务完成
		go func() {
			wg.Wait()
			// 停止活动指示器
			act.Stop()
			// 更新状态标签文本为任务完成
			text.SetText("完成！")
			// 刷新内容容器，使文本更新生效
			content.Refresh()
		}()

	})
	nameOrderNextButton.Importance = widget.DangerImportance

	// 创建“上一步”按钮，点击后切换到“子文件夹”选项卡
	nameOrderBackButton := widget.NewButton("上一步", func() {
		tabs.Select(DirsTab)
	})

	// 创建“添加命名”按钮，点击后添加一个新的命名输入框
	addNameButton := widget.NewButton("添加命名", func() {
		renameData.NameEntryIndex = append(renameData.NameEntryIndex, widget.NewEntry())
		nameContainer.Add(renameData.NameEntryIndex[len(renameData.NameEntryIndex)-1])
	})

	// 创建“删除命名”按钮，点击后删除最后一个命名输入框
	deleteNameButton := widget.NewButton("删除命名", func() {
		length := len(renameData.NameEntryIndex)
		// 若存在输入框，则删除
		if length > 0 {
			nameContainer.Remove(renameData.NameEntryIndex[length-1])
			renameData.NameEntryIndex = renameData.NameEntryIndex[:length-1]
		}
	})

	// 返回一个垂直布局的容器，包含命名输入框、操作按钮和导航按钮
	return container.NewVBox(
		nameContainer,
		// 水平布局，包含添加和删除命名按钮，两侧用间隔器居中显示
		container.NewHBox(
			layout.NewSpacer(),
			addNameButton,
			deleteNameButton,
			layout.NewSpacer(),
		),
		// 添加间隔器，使导航按钮靠下显示
		layout.NewSpacer(),
		// 水平布局，包含上一步和开始监测按钮，两侧用间隔器居中显示
		container.NewHBox(
			layout.NewSpacer(),
			nameOrderBackButton,
			nameOrderNextButton,
			layout.NewSpacer(),
		),
	)
}
