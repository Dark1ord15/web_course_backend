package ds

import (
	"Road_services/internal/app/role"

	"github.com/golang-jwt/jwt"
)

type JWTClaims struct {
	jwt.StandardClaims      // все что точно необходимо по RFC
	UserID             uint `json:"user_id"` // наши данные - uuid этого пользователя в базе данных
	Role               role.Role
}
