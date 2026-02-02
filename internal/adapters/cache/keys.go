package cache

import "fmt"

// Website URL → RSS URL
func RSSURLKey(websiteURL string) string {
	return fmt.Sprintf("rss:resolved:%s", websiteURL)
}


// ---------------------- Tokens Keys ----------------------

// Key for a refresh token associated with a user
func RefreshTokenKey(userID int64) string {
	return fmt.Sprintf("token:refresh:%d", userID)
}

// Key for an access token associated with a user
func AccessTokenKey(userID int64) string {
	return fmt.Sprintf("token:access:%d", userID)
}

// ---------------------- Optional: other caches ----------------------

// Example: cache for user profile
func UserProfileKey(userID int64) string {
	return fmt.Sprintf("user:profile:%d", userID)
}

// Example: cache for some setting
func SettingKey(name string) string {
	return fmt.Sprintf("setting:%s", name)
}
