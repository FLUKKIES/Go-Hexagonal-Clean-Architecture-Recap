package ports

// OAuthUserProfile ข้อมูล User ที่ได้จาก OAuth Provider
type OAuthUserProfile struct {
	ProviderID string
	Email      string
	FirstName  string
	LastName   string
	ProfileUrl string
}

// IOAuthProvider — Port สำหรับ OAuth Provider แต่ละเจ้า
// Implement แยกกันสำหรับ Google และ Facebook
type IOAuthProvider interface {
	// GetAuthURL สร้าง URL สำหรับ Redirect ผู้ใช้ไป Consent Screen
	GetAuthURL(state string) string
	// GetUserProfile แลก Authorization Code เป็นข้อมูล User
	GetUserProfile(code string) (*OAuthUserProfile, error)
	// ProviderName คืนชื่อ Provider เช่น "google", "facebook"
	ProviderName() string
}
