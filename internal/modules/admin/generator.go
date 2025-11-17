package admin

import "github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"

var GeneratorList = map[string]table.Generator{
	"users": GetUserTable,
	// "orders": GetOrderTable,
	// "products": GetProductTable,
}
