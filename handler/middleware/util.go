package middleware

import (
	"strings"
)

func jwtFromAuth(auth string) string {
	pieces := strings.Split(auth, " ")
	if len(pieces) != 2 {
		return ""
	}
	if strings.ToLower(pieces[0]) != "bearer" {
		return ""
	}
	return pieces[1]
}
