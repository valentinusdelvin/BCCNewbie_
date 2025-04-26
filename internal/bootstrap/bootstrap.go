package bootstrap

import (
	"fmt"
	UserHandler "hackfest-uc/internal/app/user/interface/rest"
	UserRepo "hackfest-uc/internal/app/user/repository"
	UserUsecase "hackfest-uc/internal/app/user/usecase"
	WasteDepositHandler "hackfest-uc/internal/app/waste_deposit/interface/rest"
	WasteDepositRepo "hackfest-uc/internal/app/waste_deposit/repository"
	WasteDepositUsecase "hackfest-uc/internal/app/waste_deposit/usecase"
	"hackfest-uc/internal/domain/entity"
	"hackfest-uc/internal/infra/env"
	"hackfest-uc/internal/infra/jwt"
	"hackfest-uc/internal/middleware"
	"hackfest-uc/internal/validation"

	"log"

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
		"%s%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
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

	err = database.AutoMigrate(entity.WasteDeposit{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	app := fiber.New()

	validator := validation.NewInputValidation()

	jwt := jwt.NewJWT()
	middlewareService := middleware.NewMiddleware(jwt)
	middleware.CorsMiddleware(app)

	v1 := app.Group("/api/v1")

	userRepo := UserRepo.NewUserMySQL(database)
	userUsecase := UserUsecase.NewUserUsecase(userRepo, *jwt, *validator)
	UserHandler.NewUserHandler(v1, *validator, userUsecase, middlewareService)

	wasteDepositRepo := WasteDepositRepo.NewWasteDepositMySQL(database)
	wasteDepositUsecase := WasteDepositUsecase.NewWasteDepositUsecase(wasteDepositRepo)
	WasteDepositHandler.NewWasteDepositHandler(v1, wasteDepositUsecase, middlewareService)

	return app.Listen(fmt.Sprintf(":%d", config.AppPort))
}
