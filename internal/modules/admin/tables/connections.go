package tables

import (
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
)

func GetConnectionsTable(ctx *context.Context) table.Table {

	connections := table.NewDefaultTable(ctx, table.DefaultConfigWithDriver("postgresql"))

	info := connections.GetInfo().HideFilterArea()

	info.SetTable("public.connections").SetTitle("Connections").SetDescription("Connections")

	formList := connections.GetForm()

	formList.SetTable("public.connections").SetTitle("Connections").SetDescription("Connections")

	return connections
}
