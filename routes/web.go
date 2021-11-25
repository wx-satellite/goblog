package routes

import (
	"github.com/gorilla/mux"
	"goblog/app/http/controllers"
	"goblog/app/middlewares"
	"net/http"
)

// RegisterWebRoutes 注册网页相关的路由
func RegisterWebRoutes(router *mux.Router) {

	// ForceHtml 中间件：强制内容类型为 HTML
	//router.Use(middlewares.ForceHtml)

	// template.Execute 在渲染模版的时候会正确设置 content-type
	// http.FileServer 文件目录处理器也会根据文件后缀设置正确的 content-type

	// 全局中间
	router.Use(middlewares.StartSession)

	// 静态页面处理
	pc := new(controllers.PagesController)
	router.NotFoundHandler = http.HandlerFunc(pc.NotFound)
	router.HandleFunc("/about", pc.About).Methods("GET").Name("about")

	// 文章相关页面
	ac := new(controllers.ArticlesController)
	router.HandleFunc("/", ac.Index).Methods("GET").Name("home")
	router.HandleFunc("/articles/{id:[0-9]+}", ac.Show).Methods("GET").Name("articles.show")
	router.HandleFunc("/articles", ac.Index).Methods("GET").Name("articles.index")
	router.HandleFunc("/articles/create", middlewares.Auth(ac.Create)).Methods("GET").Name("articles.create")
	router.HandleFunc("/articles", middlewares.Auth(ac.Store)).Methods("POST").Name("articles.store")
	router.HandleFunc("/articles/{id:[0-9]+}/edit", middlewares.Auth(ac.Edit)).Methods("GET").Name("articles.edit")
	router.HandleFunc("/articles/{id:[0-9]+}", middlewares.Auth(ac.Update)).Methods("POST").Name("articles.update")
	router.HandleFunc("/articles/{id:[0-9]+}/delete", middlewares.Auth(ac.Delete)).Methods("POST").Name("articles.delete")

	// 用户认证
	auc := new(controllers.AuthController)
	router.HandleFunc("/auth/register", middlewares.Guest(auc.Register)).Methods("GET").Name("auth.register")
	router.HandleFunc("/auth/do-register", middlewares.Guest(auc.DoRegister)).Methods("POST").Name("auth.doregister")
	router.HandleFunc("/auth/login", middlewares.Guest(auc.Login)).Methods("GET").Name("auth.login")
	// 路由设置成 dologin 也是可以的
	router.HandleFunc("/auth/do-login", middlewares.Guest(auc.DoLogin)).Methods("POST").Name("auth.dologin")
	router.HandleFunc("/auth/logout", middlewares.Auth(auc.Logout)).Methods("POST").Name("auth.logout")

	// 用户相关
	uc := new(controllers.UserController)
	router.HandleFunc("/users/{id:[0-9]+}", uc.Show).Methods("GET").Name("users.show")

	// 分类相关
	cc := new(controllers.CategoryController)
	router.HandleFunc("/categories/create", middlewares.Auth(cc.Create)).Methods("GET").Name("categories.create")
	router.HandleFunc("/categories", middlewares.Auth(cc.Store)).Methods("POST").Name("categories.store")

	// 静态资源库
	router.PathPrefix("/css/").Handler(http.FileServer(http.Dir("./public")))
	router.PathPrefix("/js/").Handler(http.FileServer(http.Dir("./public")))
}
