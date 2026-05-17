package main

import (
	"log"

	"github.com/FLUKKIES/Go-Hexagonal-Clean-Architecture-Recap.git/configs"
	"github.com/FLUKKIES/Go-Hexagonal-Clean-Architecture-Recap.git/domain/ports"
	"github.com/FLUKKIES/Go-Hexagonal-Clean-Architecture-Recap.git/domain/usecases"
	"github.com/FLUKKIES/Go-Hexagonal-Clean-Architecture-Recap.git/internal/adapters/gormrepo"
	"github.com/FLUKKIES/Go-Hexagonal-Clean-Architecture-Recap.git/internal/adapters/rest"
	"github.com/FLUKKIES/Go-Hexagonal-Clean-Architecture-Recap.git/internal/infrastructure/database"
	"github.com/FLUKKIES/Go-Hexagonal-Clean-Architecture-Recap.git/internal/infrastructure/emailservice"
	"github.com/FLUKKIES/Go-Hexagonal-Clean-Architecture-Recap.git/internal/infrastructure/jwtservice"
	"github.com/FLUKKIES/Go-Hexagonal-Clean-Architecture-Recap.git/internal/infrastructure/oauthprovider"
	"github.com/FLUKKIES/Go-Hexagonal-Clean-Architecture-Recap.git/internal/infrastructure/router"
	"github.com/FLUKKIES/Go-Hexagonal-Clean-Architecture-Recap.git/internal/infrastructure/smsservice"
	"github.com/gofiber/fiber/v2"
)

func main() {
	// 1. โหลด Config จาก Environment Variables
	cfg := configs.LoadConfig()

	// 2. เชื่อมต่อ Database และ Auto Migrate
	db := database.ConnectDB(cfg)

	// 3. สร้าง Repositories (GORM Implementations)
	userRepo := gormrepo.NewUserRepositoryImpl(db)
	sessionRepo := gormrepo.NewSessionRepositoryImpl(db)
	oauthRepo := gormrepo.NewOAuthRepositoryImpl(db)
	resetTokenRepo := gormrepo.NewPasswordResetRepositoryImpl(db)
	phoneOTPRepo := gormrepo.NewPhoneOTPRepositoryImpl(db)
	eventRepo := gormrepo.NewEventRepositoryImpl(db)
	eventParticipantRepo := gormrepo.NewEventParticipantRepositoryImpl(db)

	// 4. สร้าง Infrastructure Services (External Services)
	jwtSvc := jwtservice.NewJWTService(cfg.JWTSecret)
	emailSvc := emailservice.NewSMTPEmailService(
		cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUsername, cfg.SMTPPassword, cfg.SMTPFrom,
	)
	smsSvc := smsservice.NewTwilioSMSService(
		cfg.TwilioAccountSID, cfg.TwilioAuthToken, cfg.TwilioFromNumber,
	)

	// 5. สร้าง OAuth Providers
	oauthProviders := []ports.IOAuthProvider{
		oauthprovider.NewGoogleOAuthProvider(
			cfg.GoogleClientID, cfg.GoogleClientSecret, cfg.GoogleRedirectURL,
		),
		oauthprovider.NewFacebookOAuthProvider(
			cfg.FacebookClientID, cfg.FacebookClientSecret, cfg.FacebookRedirectURL,
		),
	}

	// 6. สร้าง Usecase (Domain Logic)
	authUsecase := usecases.NewAuthUsecase(
		userRepo,
		sessionRepo,
		oauthRepo,
		resetTokenRepo,
		phoneOTPRepo,
		jwtSvc,
		emailSvc,
		smsSvc,
		oauthProviders,
		usecases.AuthConfig{
			AccessTokenExpiry:  cfg.AccessTokenExpiry,
			RefreshTokenExpiry: cfg.RefreshTokenExpiry,
			ResetTokenExpiry:   cfg.ResetTokenExpiry,
			OTPExpiry:          cfg.OTPExpiry,
			OTPRateLimit:       cfg.OTPRateLimit,
			AppBaseURL:         cfg.AppBaseURL,
		},
	)

	eventUsecase := usecases.NewEventUsecase(eventRepo, eventParticipantRepo, userRepo)

	// 7. สร้าง Controllers (REST Adapters)
	authCtrl := rest.NewAuthController(authUsecase, jwtSvc)
	oauthCtrl := rest.NewOAuthController(authUsecase)
	eventCtrl := rest.NewEventController(eventUsecase)

	// 8. ตั้งค่า Fiber App
	app := fiber.New(fiber.Config{
		AppName: "Hexagonal Clean Auth API",
	})

	// 9. ผูก Routes พร้อม Inject Services
	router.SetupRoutes(app, authCtrl, oauthCtrl, eventCtrl, jwtSvc)

	// 10. เริ่มรัน Server
	log.Printf("🚀 Server is running on port %s", cfg.AppPort)
	log.Fatal(app.Listen(":" + cfg.AppPort))
}
