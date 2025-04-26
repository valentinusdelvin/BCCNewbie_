package bootstrap

import (
	"fmt"
	MarketHandler "hackfest-uc/internal/app/market/interface/rest"
	MarketRepository "hackfest-uc/internal/app/market/repository"
	MarketUseCase "hackfest-uc/internal/app/market/usecase"
	UserHandler "hackfest-uc/internal/app/user/interface/rest"
	UserRepository "hackfest-uc/internal/app/user/repository"
	UserUsecase "hackfest-uc/internal/app/user/usecase"
	WasteDepositHandler "hackfest-uc/internal/app/waste_deposit/interface/rest"
	WasteDepositRepo "hackfest-uc/internal/app/waste_deposit/repository"
	WasteDepositUsecase "hackfest-uc/internal/app/waste_deposit/usecase"
	"hackfest-uc/internal/domain/entity"
	"hackfest-uc/internal/infra/env"
	"hackfest-uc/internal/infra/jwt"
	"hackfest-uc/internal/infra/midtrans"
	"hackfest-uc/internal/infra/supabase"
	"hackfest-uc/internal/middleware"
	"hackfest-uc/internal/validation"

	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Start() error {
	_ = godotenv.Load(".env")

	config, err := env.New()
	if err != nil {
		log.Fatalf("Failed to load environment variables: %v", err)
	}

	midtrans.MidtransInit()

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

	err = database.AutoMigrate(entity.Store{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	err = database.AutoMigrate(entity.Market{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	err = database.AutoMigrate(entity.WasteDeposit{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	app := fiber.New()

	sb := supabase.Init()

	validator := validation.NewInputValidation()

	jwt := jwt.NewJWT()
	middlewareService := middleware.NewMiddleware(jwt)
	middleware.CorsMiddleware(app)

	v1 := app.Group("/api/v1")

	userRepo := UserRepository.NewUserMySQL(database)
	userUsecase := UserUsecase.NewUserUsecase(userRepo, *jwt, *validator, sb)
	UserHandler.NewUserHandler(v1, *validator, userUsecase, middlewareService, sb)

	wasteDepositRepo := WasteDepositRepo.NewWasteDepositMySQL(database)
	wasteDepositUsecase := WasteDepositUsecase.NewWasteDepositUsecase(wasteDepositRepo)
	WasteDepositHandler.NewWasteDepositHandler(v1, wasteDepositUsecase, middlewareService, *validator)

	marketRepo := MarketRepository.NewMarketMySQL(database)
	marketUsecase := MarketUseCase.NewMarketUsecase(marketRepo, sb)
	MarketHandler.NewMarketHandler(v1, marketUsecase, middlewareService)

	return app.Listen(fmt.Sprintf(":%d", config.AppPort))
}
