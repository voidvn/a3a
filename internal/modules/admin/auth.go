package admin

import (
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/auth"
	"html/template"
)

var Panel = engine.Panel{
	Title:       "S4S Backend",
	Description: "Администрирование",
}

func PageIndex(ctx *context.Context) {
	auth.ShowLoginPage(ctx, func(user auth.User) template.HTML {
		return "Добро пожаловать, " + user.Name
	})
}
