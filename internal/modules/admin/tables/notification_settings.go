package tables

import (
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
)

func GetNotificationsettingsTable(ctx *context.Context) table.Table {

	notificationSettings := table.NewDefaultTable(ctx, table.DefaultConfigWithDriver("postgresql"))

	info := notificationSettings.GetInfo().HideFilterArea()

	info.SetTable("public.notification_settings").SetTitle("Notificationsettings").SetDescription("Notificationsettings")

	formList := notificationSettings.GetForm()

	formList.SetTable("public.notification_settings").SetTitle("Notificationsettings").SetDescription("Notificationsettings")

	return notificationSettings
}
