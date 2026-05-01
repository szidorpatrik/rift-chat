package api

import (
	"fmt"
	"regexp"
)

const (
	usernameMinLength    = 3
	usernameMaxLength    = 30
	emailMinLength       = 2
	serverNameMinLength  = 3
	serverNameMaxLength  = 100
	serverDescMaxLength  = 255
	channelNameMinLength = 1
	channelNameMaxLength = 100
)

// Collapse multiple whitespace characters into a single space.
var replaceWhiteSpaceRegex = regexp.MustCompile(`\s+`)

// Allows letters, numbers, and underscores. Length: 3 to 30 characters.
var usernameRegex = regexp.MustCompile(fmt.Sprintf(`^[a-zA-Z0-9_]{%d,%d}$`, usernameMinLength, usernameMaxLength))

// Basic email format validation. Length: at least 2 characters.
var emailRegex = regexp.MustCompile(fmt.Sprintf(`^[a-zA-Z0-9._%%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{%d,}$`, emailMinLength))

// Allows any printable character, including Unicode/emojis.
// Length: 3 to 100 characters.
var serverNameRegex = regexp.MustCompile(fmt.Sprintf(`^[^\x00-\x1F\x7F]{%d,%d}$`, serverNameMinLength, serverNameMaxLength))

// Allows any printable character, including Unicode/emojis.
// Length: 0 to 255 characters.
var serverDescRegex = regexp.MustCompile(fmt.Sprintf(`^[^\x00-\x1F\x7F]{0,%d}$`, serverDescMaxLength))

// Allows any printable character, including Unicode/emojis.
// Length: 1 to 100 characters.
var channelNameRegex = regexp.MustCompile(fmt.Sprintf(`^[^\x00-\x1F\x7F]{%d,%d}$`, channelNameMinLength, channelNameMaxLength))
