package admin

import (
	_ "github.com/GoAdminGroup/go-admin/adapter/gin"
	"github.com/GoAdminGroup/go-admin/engine"
	"github.com/GoAdminGroup/go-admin/modules/config"
	"github.com/GoAdminGroup/go-admin/modules/language"
	"github.com/GoAdminGroup/go-admin/plugins/admin"
	"github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/chartjs"
	_ "github.com/GoAdminGroup/themes/adminlte"
	"github.com/gin-gonic/gin"
)

// InitAdmin инициализирует GoAdmin панель
func InitAdmin(r *gin.Engine, cfg *config.Config) *engine.Engine {
	eng := engine.Default()

	// Настройка конфигурации
	if err := eng.AddConfig(cfg).
		AddGenerators(Generators).
		Use(r); err != nil {
		panic(err)
	}

	// Регистрация плагина админки
	adminPlugin := admin.NewAdmin(Tables)
	adminPlugin.AddDisplayFilterXssJsFilter()

	// Добавление графиков
	template.AddComp(chartjs.NewChart())

	eng.AddPlugins(adminPlugin)

	return eng
}

// GetAdminConfig возвращает конфигурацию GoAdmin
func GetAdminConfig(dbURL, appKey string) *config.Config {
	return &config.Config{
		Databases: config.DatabaseList{
			"default": {
				Driver: config.DriverPostgresql,
				Host:   "localhost",
				Port:   "5432",
				User:   "user",
				Pwd:    "password",
				Name:   "s4s",
			},
		},
		UrlPrefix: "admin",
		Store: config.Store{
			Path:   "./uploads",
			Prefix: "uploads",
		},
		Language:        language.EN,
		IndexUrl:        "/",
		Debug:           false,
		ColorScheme:     "skin-blue",
		SessionLifeTime: 7200,
		//appKey:          appKey,
		Title:    "s4s Admin Panel",
		Logo:     "<b>s4s</b> Admin",
		MiniLogo: "<b>s4s</b>",
	}
}
