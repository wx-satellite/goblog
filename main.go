package main

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "<h1>Hello, 欢迎来到 goblog！</h1>")
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "此博客是用以记录编程笔记，如您有反馈或建议，请联系 "+
		"<a href=\"mailto:summer@example.com\">summer@example.com</a>")
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "<h1>请求页面未找到 :(</h1><p>如有疑惑，请联系我们。</p>")
}

type Article struct {
	Title, Body string
	Id          int64
}

func (a *Article) Link() string {
	indexUrl, _ := router.Get("articles.show").URL("id", strconv.FormatInt(a.Id, 10))

	return indexUrl.String()
}

func articlesShowHandler(w http.ResponseWriter, r *http.Request) {
	id := getRouteVariable("id", r)

	article, err := getArticleById(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeTextToResponse(w, "服务器内部错误")
		return
	}
	if article.Id <= 0 {
		w.WriteHeader(http.StatusNotFound)
		writeTextToResponse(w, "文章不存在")
		return
	}
	tmpl, err := template.New("show.tmpl").Funcs(template.FuncMap{
		"RouteNameToURL": RouteNameToURL,
		"Int64ToString":  Int64ToString,
	}).ParseFiles("resources/views/articles/show.tmpl")

	if err != nil {
		panic(err)
	}
	_ = tmpl.Execute(w, article)
}

func Int64ToString(ori int64) string {
	return strconv.FormatInt(ori, 10)
}
func RouteNameToURL(routeName string, params ...string) string {
	link, _ := router.Get(routeName).URL(params...)
	return link.String()
}

func responseToInternalServerError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	writeTextToResponse(w, "服务器内部错误")
}
func articlesIndexHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT * FROM articles")
	if err != nil {
		responseToInternalServerError(w)
		return
	}
	defer func() {
		_ = rows.Close()
	}()
	articles := make([]Article, 0, 10)
	for rows.Next() {
		article := Article{}
		err = rows.Scan(&article.Id, &article.Title, &article.Body)
		if err != nil {
			responseToInternalServerError(w)
			return
		}
		articles = append(articles, article)
	}
	if err = rows.Err(); err != nil {
		responseToInternalServerError(w)
		return
	}

	// 加载模版
	tmpl, err := template.ParseFiles("resources/views/articles/index.tmpl")
	if err != nil {
		panic(err)
	}
	// 4. 渲染模板，将所有文章的数据传输进去
	_ = tmpl.Execute(w, articles)
}

type ArticlesFormData struct {
	Title, Body string
	URL         *url.URL
	Errors      map[string]string
}

func articlesStoreHandler(w http.ResponseWriter, r *http.Request) {
	//if err := r.ParseForm(); err != nil {
	//	writeTextToResponse(w, "请传递正确的数据")
	//	return
	//}

	title := r.PostFormValue("title")
	body := r.PostFormValue("body")

	errors := validateArticleFormData(title, body)

	if len(errors) > 0 {
		// 存在错误
		storeUrl, _ := router.Get("articles.store").URL()
		data := &ArticlesFormData{
			Title:  title,
			Body:   body,
			URL:    storeUrl,
			Errors: errors,
		}
		tmp, err := template.ParseFiles("resources/views/articles/create.tmpl")
		if err != nil {
			panic(err)
		}
		_ = tmp.Execute(w, data)
		return
	}
	// 没有错误，将数据存入数据库
	id, _ := saveArticleToDBMethodTwo(title, body)
	if id <= 0 {
		w.WriteHeader(http.StatusInternalServerError)
		writeTextToResponse(w, "服务器内部错误")
		return
	}
	writeTextToResponse(w, "成功插入，ID为："+strconv.FormatInt(id, 10))
	return
}
func saveArticleToDBMethodTwo(title, body string) (lastInsertId int64, err error) {
	res, err := db.Exec("INSERT INTO articles (title, body) VALUES(?,?)", title, body)
	if err != nil {
		return
	}
	return res.LastInsertId()
}

func saveArticleToDB(title, body string) (lastInsertId int64, err error) {
	stmt, err := db.Prepare("INSERT INTO articles (title, body) VALUES(?,?)")
	if err != nil {
		return
	}
	// 关闭stmt防止占用连接
	defer func() {
		_ = stmt.Close()
	}()

	res, err := stmt.Exec(title, body)
	if err != nil {
		return
	}
	return res.LastInsertId()
}

func writeTextToResponse(w http.ResponseWriter, text string) {
	_, _ = fmt.Fprint(w, text)
}
func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// 设置请求头
func forceHTMLMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html;charset=utf-8")
		next.ServeHTTP(w, r)
	})
}

// 兼容去掉请求最后的反斜线
func removeTrailingSlash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			r.URL.Path = strings.TrimRight(r.URL.Path, "/")
		}
		next.ServeHTTP(w, r)
	})
}

var router = mux.NewRouter()

// 创建博文的表单
func articlesCreateHandler(w http.ResponseWriter, r *http.Request) {
	tmp, err := template.ParseFiles("resources/views/articles/create.tmpl")
	if err != nil {
		panic(err)
	}
	storeUrl, _ := router.Get("articles.store").URL()
	data := ArticlesFormData{
		Title:  "",
		Body:   "",
		URL:    storeUrl,
		Errors: nil,
	}

	_ = tmp.Execute(w, data)
}

var db *sql.DB

func initDb() (err error) {
	config := mysql.Config{
		User:                 "root",
		Passwd:               "admin123",
		Net:                  "tcp",
		Addr:                 "127.0.0.1:3306",
		DBName:               "goblog",
		AllowNativePasswords: true,
	}
	db, err = sql.Open("mysql", config.FormatDSN())
	if err != nil {
		panic(err)
		return
	}
	db.SetMaxIdleConns(25)
	db.SetMaxOpenConns(25)
	db.SetConnMaxLifetime(5 * time.Hour)

	if err = db.Ping(); err != nil {
		panic(err)
		return
	}
	return
}

func createTables() (err error) {
	createArticlesSQL := `CREATE TABLE IF NOT EXISTS articles(
    id bigint(20) PRIMARY KEY AUTO_INCREMENT NOT NULL,
    title varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
    body longtext COLLATE utf8mb4_unicode_ci
); `
	if _, err = db.Exec(createArticlesSQL); err != nil {
		panic(err)
	}
	return
}

func getRouteVariable(paramName string, r *http.Request) (value string) {
	vars := mux.Vars(r)
	return vars[paramName]
}

func getArticleById(id string) (res Article, err error) {
	err = db.QueryRow("SELECT * FROM articles WHERE id = ?", id).Scan(&res.Id, &res.Title, &res.Body)
	// 错误检查的时候需要区分：是未找到数据报错还是sql语句报错
	if err != nil && err != sql.ErrNoRows {
		return
	}
	err = nil
	return
}

func articlesEditHandler(w http.ResponseWriter, r *http.Request) {
	// 获取URL参数
	id := getRouteVariable("id", r)

	// 获取文章数据
	article, err := getArticleById(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeTextToResponse(w, "服务器内部错误")
		return
	}
	if article.Id <= 0 {
		w.WriteHeader(http.StatusNotFound)
		writeTextToResponse(w, "文章不存在")
		return
	}
	updateUrl, _ := router.Get("articles.update").URL("id", id)

	data := ArticlesFormData{
		Title:  article.Title,
		Body:   article.Body,
		URL:    updateUrl,
		Errors: nil,
	}

	tmpl, err := template.ParseFiles("resources/views/articles/edit.tmpl")
	if err != nil {
		panic(err)
	}

	_ = tmpl.Execute(w, data)
}

func validateArticleFormData(title, body string) (errors map[string]string) {
	errors = make(map[string]string)
	// 验证标题
	if title == "" {
		errors["title"] = "标题不能为空"
	} else if utf8.RuneCountInString(title) < 3 || utf8.RuneCountInString(title) > 40 {
		errors["title"] = "标题长度需介于 3-40"
	}
	// 验证内容
	if body == "" {
		errors["body"] = "内容不能为空"
	} else if utf8.RuneCountInString(body) < 10 {
		errors["body"] = "内容长度需大于或等于 10 个字节"
	}
	return
}

func articlesUpdateHandler(w http.ResponseWriter, r *http.Request) {
	id := getRouteVariable("id", r)
	article, err := getArticleById(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeTextToResponse(w, "服务器内部错误")
		return
	}
	if article.Id <= 0 {
		w.WriteHeader(http.StatusNotFound)
		writeTextToResponse(w, "文章不存在")
		return
	}

	// 表单验证
	title := r.PostFormValue("title")
	body := r.PostFormValue("body")
	errors := validateArticleFormData(title, body)
	if len(errors) > 0 {
		updateSql, _ := router.Get("articles.update").URL("id", id)
		data := ArticlesFormData{
			Title:  article.Title,
			Body:   article.Body,
			URL:    updateSql,
			Errors: errors,
		}
		var tmpl *template.Template
		tmpl, err = template.ParseFiles("resources/views/articles/edit.tmpl")
		if err != nil {
			panic(err)
		}
		_ = tmpl.Execute(w, data)
		return
	}

	// 没有错误，更新数据
	res, err := db.Exec("UPDATE articles SET title = ?,body=? WHERE id = ?", title, body, article.Id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeTextToResponse(w, "服务器内部错误")
		return
	}
	rows, _ := res.RowsAffected()
	if rows <= 0 {
		fmt.Println("没有做任何修改")
	}
	showUrl, _ := router.Get("articles.show").URL("id", id)
	fmt.Println("重定向的地址：", showUrl.String())
	http.Redirect(w, r, showUrl.String(), http.StatusFound)
	return
}

func articlesDeleteHandler(w http.ResponseWriter, r *http.Request) {
	id := getRouteVariable("id", r)
	article, err := getArticleById(id)
	if err != nil {
		responseToInternalServerError(w)
		return
	}
	if article.Id <= 0 {
		w.WriteHeader(http.StatusNotFound)
		writeTextToResponse(w, "文章不存在")
		return
	}

	// 上面查询过文章，因此这个sql只要保证不错误就一定可以正确删除文章
	_, err = db.Exec("DELETE FROM articles WHERE id = ?", strconv.FormatInt(article.Id, 10))
	if err != nil {
		responseToInternalServerError(w)
		return
	}

	indexUrl, _ := router.Get("articles.index").URL()
	http.Redirect(w, r, indexUrl.String(), http.StatusFound)
	return
}

func main() {
	initDb()
	createTables()
	router.HandleFunc("/", homeHandler).Methods("GET").Name("home")
	router.HandleFunc("/about", aboutHandler).Methods("GET").Name("about")

	router.HandleFunc("/articles/{id:[0-9]+}", articlesShowHandler).Methods("GET").Name("articles.show")
	router.HandleFunc("/articles/{id:[0-9]+}/edit", articlesEditHandler).Methods("GET").Name("articles.edit")
	router.HandleFunc("/articles/{id:[0-9]+}", articlesUpdateHandler).Methods("POST").Name("articles.update")
	router.HandleFunc("/articles", articlesIndexHandler).Methods("GET").Name("articles.index")
	router.HandleFunc("/articles", articlesStoreHandler).Methods("POST").Name("articles.store")
	router.HandleFunc("/articles/create", articlesCreateHandler).Methods("GET").Name("articles.create")
	router.HandleFunc("/articles/{id:[0-9]+}/delete", articlesDeleteHandler).Methods("POST").Name("articles.delete")
	// 自定义 404 页面
	router.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	// 中间件，设置content-type
	router.Use(forceHTMLMiddleware)

	http.ListenAndServe(":3000", removeTrailingSlash(router))
}
