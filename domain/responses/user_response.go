package responses

import "github.com/FLUKKIES/Go-Hexagonal-Clean-Architecture-Recap.git/domain/entities"

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

// ToUserResponse แปลง User Entity เป็น UserResponse DTO
func ToUserResponse(user *entities.User) UserResponse {
	return UserResponse{
		ID:          user.ID.String(),
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		ProfileUrl:  user.ProfileUrl,
		Role:        string(user.Role),
		IsVerified:  user.VerifiedAt != nil,
	}
}