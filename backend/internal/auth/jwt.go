package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const AccessTokenDuration = 24 * time.Hour

type Claims struct {
	OrganizationID string `json:"organization_id"`
	Role           string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(
	userID string,
	organizationID string,
	role string,
	secret string,
) (string, error) {
	now := time.Now()

	claims := Claims{
		OrganizationID: organizationID,
		Role:           role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(AccessTokenDuration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("sign access token: %w", err)
	}

	return signedToken, nil
}
