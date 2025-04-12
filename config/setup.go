package config

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"log"
	"sharing-vision-id/api/routes"
	"sharing-vision-id/db"
)

type SetupDatabase struct {
	DBMain *gorm.DB
}

func SetDatabase() *SetupDatabase {
	configDb := db.DatabaseConfig{
		MaxIdleConns: Config.Database.MAX_IDLE_CONNECTIONS,
		Host:         Config.Database.Host,
		Password:     Config.Database.Password,
		Username:     Config.Database.Username,
		Driver:       Config.Database.DRIVER,
		Port:         Config.Database.Port,
		Dbname:       Config.Database.Database,
		MaxLifetime:  Config.Database.MaxLifeTimeSeconds,
		MaxOpenConns: Config.Database.MAX_OPEN_CONNECTIONS,
	}

	configDb.ConnectionDataBaseMain()

	return &SetupDatabase{
		DBMain: db.ConnDB,
	}
}

func Start() {
	AppConfig()
	conf := GetConfig()
	SetDatabase()

	isProduction := conf.Server.Mode == "production"

	routes.Router = fiber.New(fiber.Config{
		Prefork:               false,
		CaseSensitive:         true,
		StrictRouting:         true,
		ServerHeader:          "Fiber",
		AppName:               conf.Server.AppName,
		DisableStartupMessage: isProduction,
	})

	routes.InitialRoute()
	routes.Route()
	db.ConnectRabbitMQ()
	go db.StartConsumerPayment()

	if !fiber.IsChild() {
		log.Printf("INFO: SERVICE RUNNING ON PORT %s", conf.Server.Port)
	}

	err := routes.Router.Listen(":" + conf.Server.Port)
	if err != nil {
		log.Fatalf("ERROR: cannot start server on port %s: %s", conf.Server.Port, err.Error())
		return
	}
}
