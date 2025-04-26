package bootstrap

import (
	"fmt"
	"hackfest-uc/internal/domain/entity"
	"hackfest-uc/internal/infra/env"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Start() error {
	config, err := env.New()
	if err != nil {
		panic(err)
	}

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.DBUsername,
		config.DBPassword,
		config.DBHost,
		config.DBPort,
		config.DBName,
	)

	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = database.AutoMigrate(entity.User{}){
		if err != nil{
			Log.fatalf("Failed to migrate database: %v", err)
		}
	}

	app := fiber.New()

	jwt := jwt.NewJWT()
	middlewareService := middleware.NewMiddleware(jwt)

	v1 := app.Group("/api/v1")

	return app.Listen(fmt.Sprintf(":%d", config.AppPort))
}
