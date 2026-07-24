package controller

import "strings"

const (
	minimumUsernameLength = 4
	minimumPasswordLength = 6
)

func validateRegistrationCredentials(username, password string) string {
	if len([]rune(strings.TrimSpace(username))) < minimumUsernameLength {
		return "账号至少需要4位"
	}
	return validatePassword(password)
}

func validatePassword(password string) string {
	if len([]rune(password)) < minimumPasswordLength {
		return "密码至少需要6位"
	}
	return ""
}
