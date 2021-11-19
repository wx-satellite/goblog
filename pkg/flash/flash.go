package flash

import (
	"encoding/gob"
	"goblog/pkg/session"
)

// 什么是 gob?
// 标准库 gob 是 Go 专属的编解码方式，是标准库自带的一个数据结构序列化的编码 / 解码工具。
// 类似于 JSON 或 XML，不过执行效率比他们更高。特别适合在 Go 语言程序间传递数据。
// 文档：https://blog.csdn.net/random_w/article/details/107482152

// Flashes 存入会话的数据
type Flashes map[string]interface{}

// flashKey 存入会话的键名
var flashKey = "_flashes"

func init() {
	// 在 gorilla/sessions 中存储 map 和 struct 数据需
	// 要提前注册 gob，方便后续 gob 序列化编码、解码
	gob.Register(Flashes{})
}

// Info 添加 Info 类型的消息提示
func Info(message string) {
	addFlash("info", message)
}

// Warning 添加 Warning 类型的消息提示
func Warning(message string) {
	addFlash("warning", message)
}

// Success 添加 Success 类型的消息提示
func Success(message string) {
	addFlash("success", message)
}

// Danger 添加 Danger 类型的消息提示
func Danger(message string) {
	addFlash("danger", message)
}

// addFlash 私有方法，新增一条提示
// 多次调用会覆盖的
func addFlash(key string, value string) {
	flashes := Flashes{}

	flashes[key] = value

	session.Put(flashKey, flashes)

	session.Save()
}

// All 获取所有消息
// 因为  Flash 消息是会话里的一次性数据，读取即销毁，需注意读取成功后将 flash 会话数据清空。
func All() (res Flashes) {
	val := session.Get(flashKey)
	res, ok := val.(Flashes)
	if !ok {
		return
	}
	session.Forget(flashKey)
	return
}
