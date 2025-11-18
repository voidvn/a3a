package admin

import (
	_ "github.com/GoAdminGroup/go-admin/adapter/gin"
	"github.com/GoAdminGroup/go-admin/engine"
	"github.com/GoAdminGroup/go-admin/modules/config"
	"github.com/GoAdminGroup/go-admin/modules/language"
	"github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/chartjs"
	_ "github.com/GoAdminGroup/themes/adminlte"
	"github.com/gin-gonic/gin"
)

func InitAdmin(r *gin.Engine, cfg *config.Config) *engine.Engine {
	eng := engine.Default()

	// Настройка конфигурации
	if err := eng.AddConfig(cfg).
		AddGenerators(Generators).
		Use(r); err != nil {
		panic(err)
	}

	// Регистрация плагина админки
	//adminPlugin := admin.NewAdmin(Tables)
	////adminPlugin.AddDisplayFilterXssJsFilter()
	//
	//// Добавление графиков
	template.AddComp(chartjs.NewChart())
	//eng.AddPlugins(adminPlugin)

	return eng
}

func GetAdminConfig(dbURL, appKey string) *config.Config {
	return &config.Config{
		Databases: config.DatabaseList{
			"default": {
				Driver:       config.DriverPostgresql,
				Host:         "postgres", //appConfig.GetString("DB_HOST", "localhost"),
				Port:         "5432",     //appConfig.GetString("DB_PORT", "5432"),
				User:         "postgres", //appConfig.GetString("DB_USER", "postgres"),
				Pwd:          "postgres", //appConfig.GetString("DB_PASSWORD", "postgres"),
				Name:         "s4s",      //appConfig.GetString("DB_NAME", "postgres"),
				MaxIdleConns: 50,
				MaxOpenConns: 150,
			},
		},
		UrlPrefix: "admin",
		LoginUrl:  "/admin/login",
		Store: config.Store{
			Path:   "./uploads",
			Prefix: "uploads",
		},
		Language:        language.EN,
		IndexUrl:        "/",
		Debug:           true,
		InfoLogOff:      false,
		ErrorLogOff:     false,
		AccessLogOff:    false,
		SqlLog:          true,
		ColorScheme:     "skin-blue",
		SessionLifeTime: 7200,
		//appKey:          appKey,
		Title:          "s4s Admin Panel",
		Logo:           "<b>s4s</b> Admin",
		MiniLogo:       "<b>s4s</b>",
		NoLimitLoginIP: true,
		AuthUserTable:  "goadmin_users",
		CustomHeadHtml: template.HTML(`<link rel="icon" type="image/png" sizes="32x32" href="/assets/img/favicon.ico">`),
		CustomFootHtml: template.HTML(`<div style="display:none;">Analytics code here</div>`),
		Animation: config.PageAnimation{
			Type:     "fadeInUp",
			Duration: 0.9,
		},
	}
}
