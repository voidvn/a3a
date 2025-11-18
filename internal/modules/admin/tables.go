package admin

import (
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/form"
)

// Generators - список генераторов таблиц
var Generators = map[string]table.Generator{
	"users":         GetUsersTable,
	"workflows":     GetWorkflowsTable,
	"executions":    GetExecutionsTable,
	"connections":   GetConnectionsTable,
	"subscriptions": GetSubscriptionsTable,
}

// Tables - список таблиц для админки
var Tables = table.GeneratorList{
	"users":         GetUsersTable,
	"workflows":     GetWorkflowsTable,
	"executions":    GetExecutionsTable,
	"connections":   GetConnectionsTable,
	"subscriptions": GetSubscriptionsTable,
}

// GetUsersTable - таблица пользователей
func GetUsersTable(ctx *context.Context) table.Table {
	users := table.NewDefaultTable(table.DefaultConfigWithDriver("postgresql"))

	info := users.GetInfo().HideFilterArea()

	info.AddField("ID", "id", db.Int).
		FieldSortable().
		FieldFilterable()

	info.AddField("Полное имя", "full_name", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})

	info.AddField("Email", "email", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})

	info.AddField("Телефон", "phone", db.Varchar)

	info.AddField("Город", "city", db.Varchar)

	info.AddField("Роль", "role", db.Varchar).
		FieldDisplay(func(value types.FieldModel) interface{} {
			if value.Value == "admin" {
				return "<span class='label label-danger'>Админ</span>"
			}
			return "<span class='label label-info'>Пользователь</span>"
		})

	info.AddField("Активен", "is_active", db.Bool).
		FieldDisplay(func(value types.FieldModel) interface{} {
			if value.Value == "true" {
				return "<span class='label label-success'>Да</span>"
			}
			return "<span class='label label-danger'>Нет</span>"
		}).
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(types.FieldOptions{
			{Value: "true", Text: "Активен"},
			{Value: "false", Text: "Заблокирован"},
		})

	info.AddField("Email подтвержден", "email_verified", db.Bool).
		FieldDisplay(func(value types.FieldModel) interface{} {
			if value.Value == "true" {
				return "<span class='label label-success'>Да</span>"
			}
			return "<span class='label label-warning'>Нет</span>"
		})

	info.AddField("Тариф", "subscription_plan", db.Varchar).
		FieldDisplay(func(value types.FieldModel) interface{} {
			switch value.Value {
			case "free":
				return "<span class='label label-default'>Free</span>"
			case "starter":
				return "<span class='label label-primary'>Starter</span>"
			case "team":
				return "<span class='label label-success'>Team</span>"
			default:
				return value.Value
			}
		})

	info.AddField("Создан", "created_at", db.Timestamp).
		FieldSortable()

	info.AddField("Обновлен", "updated_at", db.Timestamp)

	info.SetTable("users").SetTitle("Пользователи").SetDescription("Управление пользователями")

	// Форма редактирования
	formList := users.GetForm()

	formList.AddField("ID", "id", db.Int, form.Default).FieldDisplayButCanNotEditWhenUpdate().FieldDisableWhenCreate()
	formList.AddField("Полное имя", "full_name", db.Varchar, form.Text).FieldMust()
	formList.AddField("Email", "email", db.Varchar, form.Email).FieldMust()
	formList.AddField("Телефон", "phone", db.Varchar, form.Text)
	formList.AddField("Город", "city", db.Varchar, form.Text)
	formList.AddField("Роль", "role", db.Varchar, form.SelectSingle).
		FieldOptions(types.FieldOptions{
			{Value: "user", Text: "Пользователь"},
			{Value: "admin", Text: "Администратор"},
		}).
		FieldDefault("user")

	formList.AddField("Активен", "is_active", db.Bool, form.Switch).
		FieldOptions(types.FieldOptions{
			{Value: "true", Text: "Да"},
			{Value: "false", Text: "Нет"},
		}).
		FieldDefault("true")

	formList.AddField("Пароль", "password_hash", db.Varchar, form.Password).
		FieldHide()

	formList.SetTable("users").SetTitle("Пользователи").SetDescription("Управление пользователями")

	return users
}

// GetWorkflowsTable - таблица workflows
func GetWorkflowsTable(ctx *context.Context) table.Table {
	workflows := table.NewDefaultTable(table.DefaultConfigWithDriver("postgresql"))

	info := workflows.GetInfo().HideFilterArea()

	info.AddField("ID", "id", db.Int).FieldSortable().FieldFilterable()

	info.AddField("Название", "name", db.Varchar).
		FieldFilterable(types.FilterType{Operator: types.FilterOperatorLike})

	info.AddField("Пользователь", "user_id", db.Int).
		FieldJoin(types.Join{
			Table:     "users",
			Field:     "id",
			JoinField: "user_id",
		}).
		FieldDisplay(func(value types.FieldModel) interface{} {
			return value.Row["users__full_name"]
		})

	info.AddField("Статус", "is_active", db.Bool).
		FieldDisplay(func(value types.FieldModel) interface{} {
			if value.Value == "true" {
				return "<span class='label label-success'>Активен</span>"
			}
			return "<span class='label label-default'>Неактивен</span>"
		}).
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(types.FieldOptions{
			{Value: "true", Text: "Активные"},
			{Value: "false", Text: "Неактивные"},
		})

	info.AddField("Тип триггера", "trigger_type", db.Varchar).
		FieldDisplay(func(value types.FieldModel) interface{} {
			switch value.Value {
			case "webhook":
				return "<span class='label label-primary'>Webhook</span>"
			case "schedule":
				return "<span class='label label-info'>Расписание</span>"
			case "polling":
				return "<span class='label label-warning'>Опрос</span>"
			default:
				return value.Value
			}
		})

	info.AddField("Всего запусков", "total_executions", db.Int)

	info.AddField("Успешных", "success_count", db.Int).
		FieldDisplay(func(value types.FieldModel) interface{} {
			return "<span class='badge bg-green'>" + value.Value + "</span>"
		})

	info.AddField("Ошибок", "error_count", db.Int).
		FieldDisplay(func(value types.FieldModel) interface{} {
			if value.Value != "0" {
				return "<span class='badge bg-red'>" + value.Value + "</span>"
			}
			return value.Value
		})

	info.AddField("Создан", "created_at", db.Timestamp).FieldSortable()
	info.AddField("Обновлен", "updated_at", db.Timestamp)

	info.SetTable("workflows").SetTitle("Workflows").SetDescription("Управление рабочими процессами")

	// Форма
	formList := workflows.GetForm()

	formList.AddField("ID", "id", db.Int, form.Default).FieldDisplayButCanNotEditWhenUpdate().FieldDisableWhenCreate()
	formList.AddField("Название", "name", db.Varchar, form.Text).FieldMust()
	formList.AddField("Активен", "is_active", db.Bool, form.Switch).FieldDefault("false")
	formList.AddField("JSON конфигурация", "json_config", db.Text, form.Code).
		FieldHelpMsg("Структура workflow в JSON формате")

	formList.SetTable("workflows").SetTitle("Workflows")

	return workflows
}

// GetExecutionsTable - таблица запусков
func GetExecutionsTable(ctx *context.Context) table.Table {
	executions := table.NewDefaultTable(table.DefaultConfigWithDriver("postgresql"))

	info := executions.GetInfo().HideFilterArea()

	info.AddField("ID", "id", db.Int).FieldSortable()

	info.AddField("Workflow", "workflow_id", db.Int).
		FieldJoin(types.Join{
			Table:     "workflows",
			Field:     "id",
			JoinField: "workflow_id",
		}).
		FieldDisplay(func(value types.FieldModel) interface{} {
			return value.Row["workflows__name"]
		})

	info.AddField("Статус", "status", db.Varchar).
		FieldDisplay(func(value types.FieldModel) interface{} {
			switch value.Value {
			case "success":
				return "<span class='label label-success'>Успешно</span>"
			case "failed":
				return "<span class='label label-danger'>Ошибка</span>"
			case "running":
				return "<span class='label label-info'>Выполняется</span>"
			case "pending":
				return "<span class='label label-warning'>Ожидание</span>"
			default:
				return value.Value
			}
		}).
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(types.FieldOptions{
			{Value: "success", Text: "Успешно"},
			{Value: "failed", Text: "Ошибка"},
			{Value: "running", Text: "Выполняется"},
			{Value: "pending", Text: "Ожидание"},
		})

	info.AddField("Время старта", "started_at", db.Timestamp).FieldSortable()
	info.AddField("Время завершения", "finished_at", db.Timestamp)

	info.AddField("Длительность (сек)", "duration_seconds", db.Int).
		FieldDisplay(func(value types.FieldModel) interface{} {
			if value.Value == "" || value.Value == "0" {
				return "-"
			}
			return value.Value + " сек"
		})

	info.AddField("Ошибка", "error_message", db.Text).
		FieldDisplay(func(value types.FieldModel) interface{} {
			if value.Value != "" {
				return "<span class='text-red'>" + value.Value + "</span>"
			}
			return "-"
		})

	info.SetTable("executions").SetTitle("Запуски").SetDescription("История выполнения workflows")

	// Только просмотр, без редактирования
	formList := executions.GetForm()
	formList.AddField("ID", "id", db.Int, form.Default).FieldDisplayButCanNotEditWhenUpdate()
	formList.SetTable("executions").SetTitle("Запуски").SetDescription("Детали запуска")

	return executions
}

// GetConnectionsTable - таблица интеграций
func GetConnectionsTable(ctx *context.Context) table.Table {
	connections := table.NewDefaultTable(table.DefaultConfigWithDriver("postgresql"))

	info := connections.GetInfo()

	info.AddField("ID", "id", db.Int).FieldSortable()

	info.AddField("Пользователь", "user_id", db.Int).
		FieldJoin(types.Join{
			Table:     "users",
			Field:     "id",
			JoinField: "user_id",
		}).
		FieldDisplay(func(value types.FieldModel) interface{} {
			return value.Row["users__full_name"]
		})

	info.AddField("Сервис", "service_name", db.Varchar).
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(types.FieldOptions{
			{Value: "gmail", Text: "Gmail"},
			{Value: "slack", Text: "Slack"},
			{Value: "google_sheets", Text: "Google Sheets"},
			{Value: "telegram", Text: "Telegram"},
			{Value: "pipedrive", Text: "Pipedrive"},
		})

	info.AddField("Название", "connection_name", db.Varchar)

	info.AddField("Активна", "is_active", db.Bool).
		FieldDisplay(func(value types.FieldModel) interface{} {
			if value.Value == "true" {
				return "<span class='label label-success'>Да</span>"
			}
			return "<span class='label label-default'>Нет</span>"
		})

	info.AddField("Создана", "created_at", db.Timestamp).FieldSortable()

	info.SetTable("connections").SetTitle("Интеграции").SetDescription("Подключения к сервисам")

	formList := connections.GetForm()
	formList.AddField("ID", "id", db.Int, form.Default).FieldDisplayButCanNotEditWhenUpdate().FieldDisableWhenCreate()
	formList.AddField("Сервис", "service_name", db.Varchar, form.SelectSingle).
		FieldOptions(types.FieldOptions{
			{Value: "gmail", Text: "Gmail"},
			{Value: "slack", Text: "Slack"},
			{Value: "google_sheets", Text: "Google Sheets"},
			{Value: "telegram", Text: "Telegram"},
			{Value: "pipedrive", Text: "Pipedrive"},
		}).FieldMust()
	formList.AddField("Название", "connection_name", db.Varchar, form.Text).FieldMust()
	formList.AddField("Активна", "is_active", db.Bool, form.Switch).FieldDefault("true")

	formList.SetTable("connections").SetTitle("Интеграции")

	return connections
}

// GetSubscriptionsTable - таблица подписок
func GetSubscriptionsTable(ctx *context.Context) table.Table {
	subscriptions := table.NewDefaultTable(table.DefaultConfigWithDriver("postgresql"))

	info := subscriptions.GetInfo()

	info.AddField("ID", "id", db.Int).FieldSortable()

	info.AddField("Пользователь", "user_id", db.Int).
		FieldJoin(types.Join{
			Table:     "users",
			Field:     "id",
			JoinField: "user_id",
		}).
		FieldDisplay(func(value types.FieldModel) interface{} {
			return value.Row["users__full_name"]
		})

	info.AddField("Тариф", "plan", db.Varchar).
		FieldDisplay(func(value types.FieldModel) interface{} {
			switch value.Value {
			case "free":
				return "<span class='label label-default'>Free</span>"
			case "starter":
				return "<span class='label label-primary'>Starter ($19)</span>"
			case "team":
				return "<span class='label label-success'>Team ($99)</span>"
			default:
				return value.Value
			}
		}).
		FieldFilterable(types.FilterType{FormType: form.SelectSingle}).
		FieldFilterOptions(types.FieldOptions{
			{Value: "free", Text: "Free"},
			{Value: "starter", Text: "Starter"},
			{Value: "team", Text: "Team"},
		})

	info.AddField("Статус", "status", db.Varchar).
		FieldDisplay(func(value types.FieldModel) interface{} {
			switch value.Value {
			case "active":
				return "<span class='label label-success'>Активна</span>"
			case "canceled":
				return "<span class='label label-danger'>Отменена</span>"
			case "expired":
				return "<span class='label label-warning'>Истекла</span>"
			default:
				return value.Value
			}
		})

	info.AddField("Workflows (лимит)", "workflows_limit", db.Int)
	info.AddField("Запуски/мес (лимит)", "executions_limit", db.Int)

	info.AddField("Начало", "started_at", db.Timestamp).FieldSortable()
	info.AddField("Окончание", "expires_at", db.Timestamp)

	info.AddField("Stripe ID", "stripe_subscription_id", db.Varchar)

	info.SetTable("subscriptions").SetTitle("Подписки").SetDescription("Управление подписками")

	formList := subscriptions.GetForm()
	formList.AddField("ID", "id", db.Int, form.Default).FieldDisplayButCanNotEditWhenUpdate().FieldDisableWhenCreate()
	formList.AddField("Тариф", "plan", db.Varchar, form.SelectSingle).
		FieldOptions(types.FieldOptions{
			{Value: "free", Text: "Free"},
			{Value: "starter", Text: "Starter"},
			{Value: "team", Text: "Team"},
		})
	formList.AddField("Статус", "status", db.Varchar, form.SelectSingle).
		FieldOptions(types.FieldOptions{
			{Value: "active", Text: "Активна"},
			{Value: "canceled", Text: "Отменена"},
			{Value: "expired", Text: "Истекла"},
		})

	formList.SetTable("subscriptions").SetTitle("Подписки")

	return subscriptions
}
