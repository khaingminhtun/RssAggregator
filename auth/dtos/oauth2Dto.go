package dtos

import "strings"

// GoogleUser maps to the JSON response from Google's user info endpoint.
type GoogleUser struct {
	Email     string `json:"email"`
	FirstName string `json:"given_name"`
	LastName  string `json:"family_name"`
	AvatarURL string `json:"picture"`
}

// Method to get full name
func (u *GoogleUser) FullName() string {
	return strings.TrimSpace(u.FirstName + " " + u.LastName)
}

// GithubUser maps to the JSON response from GitHub's user endpoint.
type GithubUser struct {
	Email     string `json:"email"` // May be empty if not public
	Name      string `json:"name"`  // Full name
	AvatarURL string `json:"avatar_url"`
}
