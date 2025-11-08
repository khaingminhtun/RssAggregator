package users

import (
	"time"

	"github.com/khaingminhtun/rssagg/internal/db"
)

// UserResponse is what your API sends to client
type UserResponse struct {
	ID        int32     `json:"id"`
	Name      string    `json:"nme"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}

// SerializeUser maps DB model → API model
func SerializeUser(u *db.User) *UserResponse {
	return &UserResponse{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Password:  u.Password,
		CreatedAt: u.CreatedAt.Time,
	}
}

// SerializeUsers maps slice of users
func SerializeUsers(users []*db.User) []*UserResponse {
	out := make([]*UserResponse, len(users))
	for i, u := range users {
		out[i] = SerializeUser(u)
	}
	return out
}
