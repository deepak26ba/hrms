package middleware

import (
	"fmt"
	"hrms/common/dto"
	"hrms/internals/config"
	"time"

	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(role string, id uuid.UUID) (string, error) {

	var jwtKey = config.GetKey()

	expirationTime := time.Now().Add(120 * time.Hour)

	claims := &dto.ClaimsJWT{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		Role:   role,
		UserId: id,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		errorResponse := &dto.Error{
			Message: "Failed generating the token",
			Error:   err.Error(),
		}
		return "", fmt.Errorf("%v", errorResponse)
	}

	return tokenString, nil
}
