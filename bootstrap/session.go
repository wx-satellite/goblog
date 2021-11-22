package bootstrap

import "goblog/pkg/session"

func SetUpStore() {
	session.Initialize()
}
