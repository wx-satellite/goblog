package view

import (
	"goblog/pkg/auth"
	"goblog/pkg/logger"
	"goblog/pkg/route"
	"html/template"
	"io"
	"path/filepath"
	"strings"
)

// 填充到模板的数据
// {{ .Article.GetStringID }}
// 第一个 . 表示这个 D 数据，即 map 数据
// Article 表示取 map 键为 Article 的值，GetStringID 即调用这个值的 GetStringID 方法
// 因为值是 interface，所以这个应该使用的是反射
type D map[string]interface{}

// template 推荐阅读：https://www.cnblogs.com/52php/p/6059802.html

// go template 的自定义模板函数
// New的参数是模板名称，需要对应 ParseFiles() 中的文件名，否则会无法正确读取到模板，最终显示空白页面。
// Func() 方法的传参是 template.FuncMap 类型的 Map 对象。键为模板里调用的函数名称，值为当前上下文的函数名称
//tmpl, err := template.New("show.tmpl").Funcs(template.FuncMap{
//	"RouteNameToURL": route.NameToUrl,
//	"Uint64ToString": types.Uint64ToString,
//}).ParseFiles("resources/views/articles/show.tmpl")
//logger.Error(err)
//_ = tmpl.Execute(w, obj)

var (
	Dir = "resources/views/"
)

// Render 渲染通用视图
func Render(w io.Writer, data D, tmpFiles ...string) (err error) {
	return RenderTemplate(w, "myapp", data, tmpFiles...)
}

// RenderSimple 渲染简单的视图
func RenderSimple(w io.Writer, data D, tmpFiles ...string) (err error) {
	return RenderTemplate(w, "simple", data, tmpFiles...)
}

// RenderTemplate 渲染模板
func RenderTemplate(w io.Writer, templateName string, data D, tmpFiles ...string) (err error) {

	// 在所有模版中加入 isLogin 和 loginUser 变量
	data["isLogin"] = auth.Check()
	data["loginUser"] = auth.User()

	// 由于将模板划分成了几个布局文件共享，因此需要都加载这些文件
	// Glob 匹配所有符合规则的文件，用于获取这些布局文件
	files, err := filepath.Glob(Dir + "layouts/*.tmpl")
	if err != nil {
		logger.Error(err)
		return
	}

	for _, tmpFilePath := range tmpFiles {
		// articles.show --> articles/show
		tmpFilePath = strings.Replace(tmpFilePath, ".", "/", -1)
		files = append(files, Dir+tmpFilePath+".tmpl")
	}

	// 当使用了 ExecuteTemplate 时，name 值其实无所谓，go的模板会查找 ExecuteTemplate 指定的第二个参数对应的模板名称
	// 一般 ParseFiles 和 ExecuteTemplate 联用
	tmpl, err := template.New("").Funcs(template.FuncMap{
		"RouteName2URL": route.NameToUrl,
	}).ParseFiles(files...)

	if err != nil {
		logger.Error(err)
		return
	}

	// 模版中调用结构体的方法：{{ $article.Link }}，注意没有括号

	// {{define ... }} 是定义模板，而 {{template ...}} 是使用模板。
	// {{define ... }} 跟着的参数是模板的名称，而 {{template ...}} 有两个参数，第一个是模板，第二个是传给模板使用的数据。

	// 中间参数 name 是最终我们想要渲染的模板名称。注意这里是模板关键词 define 定义的模板名称，不是模板文件名称
	// 也就是说不是 app.tmpl 的 app 而是该文件内容 {{define "myapp"}} 中的 "myapp"
	return tmpl.ExecuteTemplate(w, templateName, data)

}
