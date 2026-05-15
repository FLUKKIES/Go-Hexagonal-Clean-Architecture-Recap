package responses

type UserResponse struct {
	ID          string  `json:"id"`
	FirstName   string  `json:"first_name"`
	LastName    string  `json:"last_name"`
	Email       string  `json:"email"`
	PhoneNumber *string `json:"phone_number,omitempty"`
	ProfileUrl  *string `json:"profile_url,omitempty"`
	Role        string  `json:"role"`
	IsVerified  bool    `json:"is_verified"` // true ถ้า VerifiedAt ไม่ใช่ nil
}