package service

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"fyne.io/fyne/v2/widget"
	"github.com/fsnotify/fsnotify"
)

type RenameData struct {
	RootPath       string
	NameEntryIndex []*widget.Entry
	DirsEntryIndex []*widget.Entry
}

var imageExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
}

// NewRenameData 创建并返回一个指向 RenameData 结构体的指针。
// 该函数会初始化一个 RenameData 结构体实例，
// 返回值：指向初始化后的 RenameData 结构体的指针。
func NewRenameData() *RenameData {
	// 返回一个使用零值初始化的 RenameData 结构体实例的指针
	return &RenameData{}
}

// isImageFile 判断给定文件路径对应的文件是否为图片文件。
// 参数 filePath 是要判断的文件的完整路径。
// 返回一个布尔值，若文件为支持的图片类型则返回 true，否则返回 false。
func isImageFile(filePath string) bool {
	// 从文件路径中提取文件扩展名，例如 ".jpg"、".png" 等
	ext := filepath.Ext(filePath)
	// 检查提取的扩展名是否存在于预定义的图片扩展名映射中
	// 若存在，映射会返回 true，否则返回 false
	return imageExtensions[ext]
}

// MonitorDirs 监测指定目录下的文件创建事件，当有图片文件创建时，按顺序对其进行重命名。
// 参数 dirPath 是要监测的目录的相对路径，该路径会与 RenameData 中的 RootPath 拼接成完整路径。
func (r *RenameData) MonitorDirs(dirPath string) {
	// 将传入的相对路径与根路径拼接成完整的目录路径
	dirPath = filepath.Join(r.RootPath, dirPath)

	// 创建一个新的文件系统监视器，忽略可能产生的错误（实际使用中建议处理错误）
	watcher, _ := fsnotify.NewWatcher()
	defer watcher.Close()

	// 将指定目录添加到监视器中，开始监测该目录下的文件变化
	watcher.Add(dirPath)

	// 已使用的命名索引，用于记录已经重命名了多少个文件
	numUsedNames := 0
	usedNames := make(map[string]bool)
	numName := len(r.NameEntryIndex)

	// 当还有可用的命名时，持续监测文件创建事件
	for numUsedNames < numName {
		select {
		// 监听文件系统监视器的事件通道
		case event, ok := <-watcher.Events:
			// 如果通道关闭，停止监测
			if !ok {
				return
			}
			// 检查事件类型是否为文件创建事件
			if event.Op&fsnotify.Create == fsnotify.Create {
				// 检查新创建的文件是否为图片文件
				if isImageFile(event.Name) {
					//跳过已重命名过的文件
					if usedNames[event.Name] {
						continue
					}
					// 再次检查是否还有可用的命名，如果没有则跳出当前处理逻辑
					if numUsedNames >= numName {
						break
					}
					// 拼接新的文件名，使用预定义的命名列表中的名称
					newName := filepath.Join(dirPath, r.NameEntryIndex[numUsedNames].Text+filepath.Ext(event.Name))
					// 等待 500 毫秒，确保文件写入完成
					time.Sleep(500 * time.Millisecond)
					// 检查文件是否存在，如果不存在则跳过本次处理
					if _, err := os.Stat(event.Name); os.IsNotExist(err) {
						continue
					}
					// 尝试对文件进行重命名
					if err := os.Rename(event.Name, newName); err != nil {
						fmt.Printf("重命名文件 %s 失败: %v\n", event.Name, err)
					} else {
						fmt.Printf("文件 %s 重命名为 %s\n", event.Name, newName)
						numUsedNames++
						usedNames[newName] = true
					}
				}
			}
		// 监听文件系统监视器的错误通道
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			fmt.Printf("监测文件夹 %s 时出错: %v\n", dirPath, err)
		}
	}
}
