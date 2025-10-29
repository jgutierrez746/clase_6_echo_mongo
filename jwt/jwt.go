package jwt

import (
	"os"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

func GenerarJWT(correo, nombre, id string) (string, error) {
	miClave := []byte(os.Getenv("SECRET_JWT"))

	/*if len(miClave) == 0 {
		return nil,
	}*/

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"correo":         correo,
		"nombre":         nombre,
		"generado_desde": "https://www.cesarcancino.com",
		"id":             id,
		"iat":            time.Now().Unix(),
		"exp":            time.Now().Add(time.Hour * 24).Unix(), // 24 horas
	})

	tokenString, err := token.SignedString(miClave)

	return tokenString, err
}
