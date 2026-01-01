package users

import (
	"time"

	"github.com/khaingminhtun/rssagg/internal/adapters/database/repo"
)

// UserDTO is the "Layer" model for the frontend
type UserDTO struct {
    ID        int32  `json:"id"`
    Name      string `json:"name"`
    Email     string `json:"email"`
    AuthType  string `json:"auth_type"`
    CreatedAt string `json:"created_at"` // Standard string or time.Time
    UpdatedAt string `json:"updated_at"`
}

func ToUserDTO(user repo.User) UserDTO {
    return UserDTO{
        ID:       user.ID,
        Name:     user.Name,
        Email:    user.Email,
        AuthType: user.AuthType.String,
        // Convert pgtype.Timestamptz to a clean string or time.Time
        CreatedAt: user.CreatedAt.Time.Format(time.RFC3339),
        UpdatedAt: user.UpdatedAt.Time.Format(time.RFC3339),
    }
}