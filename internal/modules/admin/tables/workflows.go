package tables

import (
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
)

func GetWorkflowsTable(ctx *context.Context) table.Table {

	workflows := table.NewDefaultTable(ctx, table.DefaultConfigWithDriver("postgresql"))

	info := workflows.GetInfo().HideFilterArea()

	info.SetTable("public.workflows").SetTitle("Workflows").SetDescription("Workflows")

	formList := workflows.GetForm()

	formList.SetTable("public.workflows").SetTitle("Workflows").SetDescription("Workflows")

	return workflows
}
