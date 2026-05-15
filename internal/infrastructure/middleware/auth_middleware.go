package middleware

import (
	"strings"

	"github.com/FLUKKIES/Go-Hexagonal-Clean-Architecture-Recap.git/domain/ports"
	"github.com/gofiber/fiber/v2"
)

// AuthMiddleware ใช้สำหรับป้องกัน Route ที่ต้อง Login ก่อนถึงจะเข้าได้
func AuthMiddleware(jwtService ports.IJWTService) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// 1. ดึง Token จาก Cookie ก่อน
		tokenString := ctx.Cookies("access_token")
		
		// 1.1 ถ้าไม่มีใน Cookie ให้ลองดึงจาก Header "Authorization: Bearer <token>" (เผื่อ Mobile App ใช้)
		if tokenString == "" {
			authHeader := ctx.Get("Authorization")
			if authHeader != "" {
				parts := strings.Split(authHeader, " ")
				if len(parts) == 2 && parts[0] == "Bearer" {
					tokenString = parts[1]
				}
			}
		}

		if tokenString == "" {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "missing or invalid authorization token",
			})
		}

		// 2. ตรวจสอบ JWT ผ่าน Port
		claims, err := jwtService.ValidateAccessToken(tokenString)
		if err != nil {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid or expired token",
			})
		}

		// 3. แนบข้อมูล User (ID, Role) ลงไปใน Context เพื่อให้ Controller นำไปใช้ต่อได้ง่ายๆ
		ctx.Locals("user_id", claims.UserID.String())
		ctx.Locals("user_role", claims.Role)

		// 4. ไปยัง Handler ถัดไป (เช่น Controller)
		return ctx.Next()
	}
}

// RoleMiddleware ใช้สำหรับจำกัดสิทธิ์ (เช่น ต้องเป็น admin เท่านั้น)
// หมายเหตุ: ต้องเรียกต่อจาก AuthMiddleware เสมอ
func RoleMiddleware(requiredRoles ...string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		userRole, ok := ctx.Locals("user_role").(string)
		if !ok {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "unauthorized",
			})
		}

		for _, role := range requiredRoles {
			if userRole == role {
				return ctx.Next() // เจอ Role ที่ตรงกัน ผ่านได้
			}
		}

		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "forbidden: insufficient permissions",
		})
	}
}
