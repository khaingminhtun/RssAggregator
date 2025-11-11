package dtos

import (
	"time"

	"github.com/khaingminhtun/rssagg/internal/db"
	"github.com/khaingminhtun/rssagg/utils"
)

// UserResponse is what your API sends to client
type UserResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"nme"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"update_at"`
}

// SerializeUser maps DB model → API model
func SerializeUser(u *db.User) *UserResponse {

	return &UserResponse{
		ID:        utils.UUIDToString(u.ID),
		Name:      u.Name,
		Email:     u.Email,
		CreatedAt: u.CreatedAt.Time,
		UpdatedAt: u.UpdatedAt.Time,
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
