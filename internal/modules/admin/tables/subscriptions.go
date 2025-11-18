package tables

import (
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
)

func GetSubscriptionsTable(ctx *context.Context) table.Table {

	subscriptions := table.NewDefaultTable(ctx, table.DefaultConfigWithDriver("postgresql"))

	info := subscriptions.GetInfo().HideFilterArea()

	info.SetTable("public.subscriptions").SetTitle("Subscriptions").SetDescription("Subscriptions")

	formList := subscriptions.GetForm()

	formList.SetTable("public.subscriptions").SetTitle("Subscriptions").SetDescription("Subscriptions")

	return subscriptions
}
