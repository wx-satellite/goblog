package bootstrap

import (
	"github.com/gorilla/mux"
	"goblog/pkg/route"
	"goblog/routes"
)

// SetupRoute 路由初始化
func SetupRoute() *mux.Router {
	route.Initialize()

	routes.RegisterWebRoutes(route.Router)
	return route.Router
}
