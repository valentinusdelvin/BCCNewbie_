package bootstrap

import (
	"fmt"
	"hackfest-uc/internal/app/user/interface/rest"
	"hackfest-uc/internal/app/user/repository"
	"hackfest-uc/internal/app/user/usecase"
	"hackfest-uc/internal/domain/entity"
	"hackfest-uc/internal/infra/env"
	"hackfest-uc/internal/infra/jwt"
	"hackfest-uc/internal/middleware"

	"log"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Start() error {
	config, err := env.New()
	if err != nil {
		log.Fatalf("Failed to load environment variables: %v", err)
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

	err = database.AutoMigrate(entity.User{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	app := fiber.New()

	validator := validator.New()

	jwt := jwt.NewJWT()
	middlewareService := middleware.NewMiddleware(jwt)
	middleware.CorsMiddleware(app)

	v1 := app.Group("/api/v1")

	userRepo := repository.NewUserMySQL(database)
	userUsecase := usecase.NewUserUsecase(userRepo, *jwt)
	rest.NewUserHandler(v1, userUsecase, *validator, middlewareService)

	return app.Listen(fmt.Sprintf(":%d", config.AppPort))
}
