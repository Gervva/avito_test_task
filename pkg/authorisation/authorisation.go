package authorisation

import (
 	"strings"
)

type UserRole int

const (
	UserRoleAdmin    UserRole = 1
	UserRoleRegular  UserRole = 2
	UserUnauthorized UserRole = 3
)
   
func GetUser(token string) (bool, UserRole) {
	if strings.HasPrefix(token, "a") {
	 	return true, UserRoleAdmin
	} else if strings.HasPrefix(token, "u") {
	 	return true, UserRoleRegular
	}
   
	return false, UserUnauthorized
}
