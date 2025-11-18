package tables

import "github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"

// The key of Generators is the prefix of table info url.
// The corresponding value is the Form and Table data.
//
// http://{{config.Domain}}:{{Port}}/{{config.Prefix}}/info/{{key}}
//
// example:
//
// "migrations" => http://localhost:9033/admin/info/migrations
// "users" => http://localhost:9033/admin/info/users
// "workflows" => http://localhost:9033/admin/info/workflows
// "executions" => http://localhost:9033/admin/info/executions
// "connections" => http://localhost:9033/admin/info/connections
// "subscriptions" => http://localhost:9033/admin/info/subscriptions
// "notification_settings" => http://localhost:9033/admin/info/notification_settings
//
// example end
var Generators = map[string]table.Generator{

	"users":                 GetUsersTable,
	"workflows":             GetWorkflowsTable,
	"executions":            GetExecutionsTable,
	"connections":           GetConnectionsTable,
	"subscriptions":         GetSubscriptionsTable,
	"notification_settings": GetNotificationsettingsTable,

	// generators end
}
