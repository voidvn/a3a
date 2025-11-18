package admin

import (
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/template/chartjs"
	"github.com/GoAdminGroup/go-admin/template/types"
	"html/template"
)

// GetDashboard возвращает дашборд админки с метриками
func GetDashboard(ctx *context.Context) (types.Panel, error) {
	// Получение статистики из БД
	conn := db.GetConnection(ctx.Request)

	// Общая статистика
	totalUsers, _ := conn.Table("users").Count()
	activeUsers, _ := conn.Table("users").Where("is_active", "=", true).Count()
	totalWorkflows, _ := conn.Table("workflows").Count()
	activeWorkflows, _ := conn.Table("workflows").Where("is_active", "=", true).Count()
	totalExecutions, _ := conn.Table("executions").Count()
	successExecutions, _ := conn.Table("executions").Where("status", "=", "success").Count()
	failedExecutions, _ := conn.Table("executions").Where("status", "=", "failed").Count()

	// Подписки
	freeSubscriptions, _ := conn.Table("subscriptions").Where("plan", "=", "free").Count()
	starterSubscriptions, _ := conn.Table("subscriptions").Where("plan", "=", "starter").Count()
	teamSubscriptions, _ := conn.Table("subscriptions").Where("plan", "=", "team").Count()

	// Создание дашборда
	col1 := types.NewCol().SetSize(types.SizeMD(3)).SetContent(template.HTML(`
		<div class="info-box">
			<span class="info-box-icon bg-aqua"><i class="fa fa-users"></i></span>
			<div class="info-box-content">
				<span class="info-box-text">Пользователи</span>
				<span class="info-box-number">` + totalUsers.String() + `</span>
				<small>` + activeUsers.String() + ` активных</small>
			</div>
		</div>
	`)).GetContent()

	col2 := types.NewCol().SetSize(types.SizeMD(3)).SetContent(template.HTML(`
		<div class="info-box">
			<span class="info-box-icon bg-green"><i class="fa fa-project-diagram"></i></span>
			<div class="info-box-content">
				<span class="info-box-text">Workflows</span>
				<span class="info-box-number">` + totalWorkflows.String() + `</span>
				<small>` + activeWorkflows.String() + ` активных</small>
			</div>
		</div>
	`)).GetContent()

	col3 := types.NewCol().SetSize(types.SizeMD(3)).SetContent(template.HTML(`
		<div class="info-box">
			<span class="info-box-icon bg-yellow"><i class="fa fa-play-circle"></i></span>
			<div class="info-box-content">
				<span class="info-box-text">Запуски</span>
				<span class="info-box-number">` + totalExecutions.String() + `</span>
				<small class="text-success">` + successExecutions.String() + ` успешных</small> / 
				<small class="text-danger">` + failedExecutions.String() + ` ошибок</small>
			</div>
		</div>
	`)).GetContent()

	col4 := types.NewCol().SetSize(types.SizeMD(3)).SetContent(template.HTML(`
		<div class="info-box">
			<span class="info-box-icon bg-red"><i class="fa fa-credit-card"></i></span>
			<div class="info-box-content">
				<span class="info-box-text">Подписки</span>
				<span class="info-box-number">Free: ` + freeSubscriptions.String() + `</span>
				<small>Starter: ` + starterSubscriptions.String() + ` | Team: ` + teamSubscriptions.String() + `</small>
			</div>
		</div>
	`)).GetContent()

	row1 := types.NewRow().SetContent(col1 + col2 + col3 + col4).GetContent()

	// График запусков по дням (последние 30 дней)
	executionsChart := getExecutionsChart(conn)

	// График по типам workflows
	workflowTypesChart := getWorkflowTypesChart(conn)

	row2 := types.NewRow().
		SetContent(
			types.NewCol().SetSize(types.SizeMD(6)).SetContent(executionsChart).GetContent() +
				types.NewCol().SetSize(types.SizeMD(6)).SetContent(workflowTypesChart).GetContent(),
		).GetContent()

	// Таблица последних активностей
	recentActivity := getRecentActivity(conn)

	row3 := types.NewRow().
		SetContent(
			types.NewCol().SetSize(types.SizeMD(12)).SetContent(recentActivity).GetContent(),
		).GetContent()

	return types.Panel{
		Content: row1 + row2 + row3,
		Title:   "Дашборд",
	}, nil
}

func getExecutionsChart(conn db.Connection) template.HTML {
	// Получение данных по запускам за последние 30 дней
	execData, _ := conn.Query(`
		SELECT 
			DATE(started_at) as date,
			COUNT(*) as total,
			SUM(CASE WHEN status = 'success' THEN 1 ELSE 0 END) as success,
			SUM(CASE WHEN status = 'failed' THEN 1 ELSE 0 END) as failed
		FROM executions
		WHERE started_at >= NOW() - INTERVAL '30 days'
		GROUP BY DATE(started_at)
		ORDER BY date
	`)

	labels := "["
	totalData := "["
	successData := "["
	failedData := "["

	for i, row := range execData {
		if i > 0 {
			labels += ","
			totalData += ","
			successData += ","
			failedData += ","
		}
		labels += `"` + row["date"].(string) + `"`
		totalData += row["total"].(string)
		successData += row["success"].(string)
		failedData += row["failed"].(string)
	}

	labels += "]"
	totalData += "]"
	successData += "]"
	failedData += "]"

	chartHTML := chartjs.Line().
		SetID("executions-chart").
		SetTitle("Запуски за последние 30 дней").
		SetHeight(300).
		SetLabels(template.HTML(labels)).
		AddDataSet("Всего").
		DSData(template.HTML(totalData)).
		DSFill(false).
		DSBorderColor("rgb(75, 192, 192)").
		DSLineTension(0.1)

	return chartHTML.GetContent()
}

func getWorkflowTypesChart(conn db.Connection) template.HTML {
	// Статистика по типам триггеров
	triggerData, _ := conn.Query(`
		SELECT trigger_type, COUNT(*) as count
		FROM workflows
		GROUP BY trigger_type
	`)

	labels := "["
	data := "["
	colors := "["

	colorMap := map[string]string{
		"webhook":  "#3498db",
		"schedule": "#2ecc71",
		"polling":  "#f39c12",
	}

	for i, row := range triggerData {
		if i > 0 {
			labels += ","
			data += ","
			colors += ","
		}
		triggerType := row["trigger_type"].(string)
		labels += `"` + triggerType + `"`
		data += row["count"].(string)
		colors += `"` + colorMap[triggerType] + `"`
	}

	labels += "]"
	data += "]"
	colors += "]"

	chartHTML := chartjs.Pie().
		SetID("workflow-types-chart").
		SetTitle("Типы Workflows").
		SetHeight(300).
		SetLabels(template.HTML(labels)).
		AddDataSet("Workflows").
		DSData(template.HTML(data)).
		DSBackgroundColor(template.HTML(colors))

	return chartHTML.GetContent()
}

func getRecentActivity(conn db.Connection) template.HTML {
	// Последние 10 запусков
	recentExec, _ := conn.Query(`
		SELECT 
			e.id,
			w.name as workflow_name,
			u.full_name as user_name,
			e.status,
			e.started_at,
			e.duration_seconds
		FROM executions e
		JOIN workflows w ON e.workflow_id = w.id
		JOIN users u ON w.user_id = u.id
		ORDER BY e.started_at DESC
		LIMIT 10
	`)

	tableHTML := `
		<div class="box box-primary">
			<div class="box-header with-border">
				<h3 class="box-title">Последние запуски</h3>
			</div>
			<div class="box-body">
				<table class="table table-striped">
					<thead>
						<tr>
							<th>ID</th>
							<th>Workflow</th>
							<th>Пользователь</th>
							<th>Статус</th>
							<th>Время</th>
							<th>Длительность</th>
						</tr>
					</thead>
					<tbody>
	`

	for _, row := range recentExec {
		status := row["status"].(string)
		statusLabel := ""
		switch status {
		case "success":
			statusLabel = `<span class="label label-success">Успешно</span>`
		case "failed":
			statusLabel = `<span class="label label-danger">Ошибка</span>`
		case "running":
			statusLabel = `<span class="label label-info">Выполняется</span>`
		}

		tableHTML += `<tr>
			<td>` + row["id"].(string) + `</td>
			<td>` + row["workflow_name"].(string) + `</td>
			<td>` + row["user_name"].(string) + `</td>
			<td>` + statusLabel + `</td>
			<td>` + row["started_at"].(string) + `</td>
			<td>` + row["duration_seconds"].(string) + ` сек</td>
		</tr>`
	}

	tableHTML += `
					</tbody>
				</table>
			</div>
		</div>
	`

	return template.HTML(tableHTML)
}
