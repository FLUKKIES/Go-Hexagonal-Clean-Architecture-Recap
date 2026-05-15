package router

import (
	"github.com/FLUKKIES/Go-Hexagonal-Clean-Architecture-Recap.git/domain/ports"
	"github.com/FLUKKIES/Go-Hexagonal-Clean-Architecture-Recap.git/internal/adapters/rest"
	"github.com/FLUKKIES/Go-Hexagonal-Clean-Architecture-Recap.git/internal/infrastructure/middleware"
	"github.com/gofiber/fiber/v2"
)

// SetupRoutes จัดการ Route ทั้งหมดของแอป แยกออกจาก main.go
func SetupRoutes(app *fiber.App, authCtrl *rest.AuthController, oauthCtrl *rest.OAuthController, jwtService ports.IJWTService) {
	api := app.Group("/api")

	// Middleware
	authMid := middleware.AuthMiddleware(jwtService)

	// Auth Routes
	auth := api.Group("/auth")
	auth.Post("/register", authCtrl.Register)
	auth.Post("/login", authCtrl.Login)
	auth.Post("/refresh", authCtrl.RefreshToken)
	auth.Post("/forgot-password", authCtrl.ForgotPassword)
	auth.Post("/reset-password", authCtrl.ResetPassword)

	// Protected Routes (ต้องมี JWT Token)
	protected := api.Group("/auth", authMid)
	protected.Post("/logout", authCtrl.Logout)
	protected.Post("/logout-all", authCtrl.LogoutAll)

	// Phone OTP Routes (ต้อง Login แล้ว)
	protected.Post("/phone/send-otp", authCtrl.SendPhoneOTP)
	protected.Post("/phone/verify-otp", authCtrl.VerifyPhoneOTP)

	// OAuth Routes
	auth.Get("/:provider", oauthCtrl.RedirectToProvider)
	auth.Get("/:provider/callback", oauthCtrl.HandleCallback)
}
