package jwt

import (
	"sso/internal/domains/models"
	"time"

	"github.com/golang-jwt/jwt"
)

func New(user models.User, app models.App, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.Id
	claims["exp"] = time.Now().Add(duration).Unix()
	claims["AppId"] = app.Id

	tokenString, err := token.SignedString("SupposedToBeSecret") // todo: Somehow link app and its secret
	if err != nil {
		return "", err
	}

	return tokenString, nil
}