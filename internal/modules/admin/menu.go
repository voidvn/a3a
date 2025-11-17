package admin

import "github.com/GoAdminGroup/go-admin/modules/config"

func GetMenu() config.Menu {
	return config.Menu{
		{Title: "Главная", Header: true},
		{Name: "dashboard", Title: "Дашборд", Url: "/admin"},
		{Name: "users", Title: "Пользователи", Url: "/admin/info/users"},
		{Name: "orders", Title: "Заказы", Url: "/admin/info/orders"},
		// добавляй сколько угодно
	}
}
