package tables

import (
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
)

func GetExecutionsTable(ctx *context.Context) table.Table {

	executions := table.NewDefaultTable(ctx, table.DefaultConfigWithDriver("postgresql"))

	info := executions.GetInfo().HideFilterArea()

	info.SetTable("public.executions").SetTitle("Executions").SetDescription("Executions")

	formList := executions.GetForm()

	formList.SetTable("public.executions").SetTitle("Executions").SetDescription("Executions")

	return executions
}
