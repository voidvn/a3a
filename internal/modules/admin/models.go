package admin

import (
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	_ "github.com/GoAdminGroup/themes/adminlte"
	"your-project/internal/models"
)

func GetUserTable(ctx *table.Context) table.Table {
	userTable := table.NewDefaultTable(table.DefaultConfigWithDriver("postgresql"))

	info := userTable.GetInfo().HideFilterArea()
	info.AddField("ID", "id", db.Int).FieldFilterable()
	info.AddField("Имя", "name", db.Varchar)
	info.AddField("Email", "email", db.Varchar)
	info.AddField("Создано", "created_at", db.Timestamp).
		FieldSortable().
		FieldFilterable(table.RangeFilterType{})

	info.SetTable("users").SetTitle("Пользователи").SetDescription("Управление пользователями")

	form := userTable.GetForm()
	form.AddField("Имя", "name", db.Varchar, form.Text)
	form.AddField("Email", "email", db.Varchar, form.Email)

	form.SetTable("users").SetTitle("Пользователь").SetDescription("Редактирование")

	return userTable
}

// Повтори для других моделей: orders, products, etc.
