package config

import (
	"encoding/json"
	"io"
	"os"

	"fyne.io/fyne/v2"
)

type UserConfig []string

// NewUserConfig 创建一个新的 UserConfig 实例，并返回其指针。
// UserConfig 是一个字符串切片类型，用于存储用户配置项。
// 该函数返回一个指向空 UserConfig 切片的指针，可用于后续添加配置项。
func NewUserConfig() *UserConfig {
	// 使用 & 操作符创建一个新的空 UserConfig 切片，并返回其内存地址。
	return &UserConfig{}
}

// ReadConfig 从 config.json 文件中读取用户配置，并将其解析到 UserConfig 实例中。
// 如果文件不存在或读取、解析过程中出现错误，会通过 fyne 应用发送通知提示用户。
// 参数 ap 是一个 fyne.App 实例，用于发送通知。
func (config *UserConfig) ReadConfig(ap fyne.App) {
	// 尝试打开配置文件 config.json
	file, err := os.Open("config.json")
	if err != nil {
		if os.IsNotExist(err) {
			// 若文件不存在，发送提示通知给用户
			ap.SendNotification(&fyne.Notification{
				Title:   "提示",
				Content: "未找到配置文件，请先在设置中完成相关设置并保存",
			})
			return
		} else {
			// 若出现其他错误，发送错误通知给用户
			ap.SendNotification(&fyne.Notification{
				Title:   "错误",
				Content: "打开配置文件时出错，请联系开发者",
			})
			return
		}
	}
	// 确保文件在函数结束时关闭
	defer file.Close()

	// 读取配置文件的全部内容
	data, err := io.ReadAll(file)
	if err != nil {
		// 若读取文件出错，发送错误通知给用户
		ap.SendNotification(&fyne.Notification{
			Title:   "错误",
			Content: "读取配置文件时出错，请联系开发者",
		})
		return
	}

	// 将读取到的 JSON 数据解析到 UserConfig 实例中
	err = json.Unmarshal(data, config)
	if err != nil {
		// 若解析 JSON 数据出错，发送错误通知给用户
		ap.SendNotification(&fyne.Notification{
			Title:   "错误",
			Content: "解析JSON时出错，请联系开发者",
		})
		return
	}
}

// SaveConfigFile 将 UserConfig 实例中的配置信息保存到 config.json 文件中。
// 如果序列化配置信息或写入文件过程中出现错误，会通过 fyne 应用发送通知提示用户。
// 参数 ap 是一个 fyne.App 实例，用于发送通知。
func (config *UserConfig) SaveConfigFile(ap fyne.App) {
	// 使用 json.Marshal 将 UserConfig 实例序列化为 JSON 字节数据
	jsonData, err := json.Marshal(config)
	if err != nil {
		// 若序列化过程出错，发送错误通知给用户
		ap.SendNotification(&fyne.Notification{
			Title:   "错误",
			Content: "序列化配置文件时出错，请联系开发者",
		})
		return
	}
	// 使用 os.WriteFile 将序列化后的 JSON 数据写入 config.json 文件，文件权限设置为 0644
	err = os.WriteFile("config.json", jsonData, 0644)
	if err != nil {
		// 若写入文件过程出错，发送错误通知给用户
		ap.SendNotification(&fyne.Notification{
			Title:   "错误",
			Content: "保存配置文件时出错，请联系开发者",
		})
		return
	}
}
