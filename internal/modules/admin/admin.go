package admin

import (
	"github.com/GoAdminGroup/go-admin/engine"
	"github.com/GoAdminGroup/go-admin/plugins/admin"
	"github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/chartjs"
	"github.com/GoAdminGroup/themes/adminlte"
	"github.com/gin-gonic/gin"
)

func Init(r *gin.Engine) {
	eng := engine.Default()

	// Настройка темы и шаблонов
	cfg := engine.Config{
		Databases: engine.Databases{{
			Host:       "postgres", // имя сервиса в docker-compose
			Port:       "5432",
			User:       "postgres",
			Pass:       "postgres",
			Name:       "s4s",
			MaxIdleCon: 50,
			MaxOpenCon: 150,
			Driver:     engine.DriverPostgresql,
		}},
		UrlPrefix: "admin",
		Store: engine.Store{
			Path:   "./uploads",
			Prefix: "uploads",
		},
		Language:    "ru", // или "en"
		IndexUrl:    "/",
		Debug:       true,
		ColorScheme: adminlte.ColorschemeSkinBlack,
	}

	template.AddComp(chartjs.NewChart())

	if err := eng.AddConfig(&cfg).
		AddGenerators(GeneratorList).
		Use(r); err != nil {
		panic(err)
	}

	// Страница входа
	r.GET("/admin", func(c *gin.Context) {
		engine.Content(c, func(ctx *gin.Context) (engine.Panel, error) {
			return Panel{Title: "Добро пожаловать", Description: "Войдите в систему"}, nil
		})
	})

	eng.HTML("GET", "/admin", PageIndex)
}
