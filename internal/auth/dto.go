package auth

import "github.com/golang-jwt/jwt"

type SignIn struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	UserId   string `json:"user_id"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

type GenerateToken struct {
	Username string
	UserId   string
	Role     string
}

type TokenData struct {
	Username *string
	UserId   *string
	Role     *string
}
