package utils

import (
	"errors"
	"net/url"
)

func ValidateURL(raw string) (*url.URL, error) {
	u, err := url.ParseRequestURI(raw)
	if err != nil {
		return nil, err
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return nil, errors.New("invalid scheme")
	}
	return u, nil
}

// resolveURL resolves a possibly relative feed URL against the site URL.
// Examples:
//
//	siteURL: https://techcrunch.com
//	feedURL: /feed
//	result:  https://techcrunch.com/feed
func ResolveURL(siteURL string, feedURL string) string {
	base, err := url.Parse(siteURL)
	if err != nil {
		return feedURL
	}

	ref, err := url.Parse(feedURL)
	if err != nil {
		return feedURL
	}

	// If feedURL is already absolute, ResolveReference returns it unchanged
	return base.ResolveReference(ref).String()
}
