package authService

import (
	"errors"
	"tosinjs/reminder-service/internal/entity/errorEntity"
	"tosinjs/reminder-service/internal/entity/tokenEntity"

	"github.com/golang-jwt/jwt/v4"
)

type authService struct {
	key string
}

type AuthService interface {
	ValidateJWT(token string) (*tokenEntity.IDTokenPayload, *errorEntity.LayerError)
}

type Claims struct {
	Id         string `json:"id"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	IsVerified bool   `json:"isVerified"`
	jwt.StandardClaims
}

func New(key string) AuthService {
	return authService{key: key}
}

func (a authService) ValidateJWT(jwtToken string) (*tokenEntity.IDTokenPayload, *errorEntity.LayerError) {
	token, err := jwt.ParseWithClaims(jwtToken, &Claims{}, func(t *jwt.Token) (any, error) {
		return []byte(a.key), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errorEntity.New(401, "service", err.Error(), err)
		}
		return nil, errorEntity.New(401, "service", "Unauthorized", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return &tokenEntity.IDTokenPayload{
			Id:         claims.Id,
			Username:   claims.Username,
			Email:      claims.Email,
			IsVerified: claims.IsVerified,
		}, nil
	}

	return nil, errorEntity.InternalServerError("service", nil)
}
