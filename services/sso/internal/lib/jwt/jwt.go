package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/leyl1ne/LocalFarmersMarket/services/sso/internal/domain/models"
)

// TODO: написать тест
func NewToken(user models.User, app models.App, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodES256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.ID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(duration).Unix()
	claims["app_id"] = app.ID

	//подписываем токен, используя секретный ключ приложения
	tokenString, err := token.SignedString([]byte(app.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
