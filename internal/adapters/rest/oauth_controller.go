package rest

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/FLUKKIES/Go-Hexagonal-Clean-Architecture-Recap.git/domain/usecases"
	"github.com/gofiber/fiber/v2"
)

type OAuthController struct {
	usecase usecases.IAuthUsecase
}

func NewOAuthController(uc usecases.IAuthUsecase) *OAuthController {
	return &OAuthController{usecase: uc}
}

// GET /api/auth/:provider — Redirect ไปหน้า Consent Screen
// provider = "google" | "facebook"
func (c *OAuthController) RedirectToProvider(ctx *fiber.Ctx) error {
	provider := ctx.Params("provider")

	// สร้าง State แบบ Cryptographically Random สำหรับ CSRF Protection
	state, err := generateOAuthState()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}

	// เก็บ State ไว้ใน Cookie เพื่อเปรียบเทียบตอน Callback
	ctx.Cookie(&fiber.Cookie{
		Name:     "oauth_state",
		Value:    state,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
		MaxAge:   300, // 5 นาที
	})

	authURL, err := c.usecase.GetOAuthURL(provider, state)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("unsupported provider: %s", provider)})
	}

	return ctx.Redirect(authURL, fiber.StatusTemporaryRedirect)
}

// GET /api/auth/:provider/callback — รับ Code กลับจาก Provider
func (c *OAuthController) HandleCallback(ctx *fiber.Ctx) error {
	provider := ctx.Params("provider")

	// ตรวจสอบ State ป้องกัน CSRF Attack
	cookieState := ctx.Cookies("oauth_state")
	queryState := ctx.Query("state")
	if cookieState == "" || cookieState != queryState {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid oauth state"})
	}

	// ล้าง State Cookie
	ctx.ClearCookie("oauth_state")

	code := ctx.Query("code")
	if code == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "authorization code not found"})
	}

	resp, err := c.usecase.HandleOAuthCallback(provider, code)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "oauth authentication failed"})
	}

	return ctx.JSON(resp)
}

func generateOAuthState() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
