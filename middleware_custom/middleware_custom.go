package middleware_custom

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func ValidarJWT(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Header 'Authorization' es requerido"})
		}

		splitBearer := strings.Split(authHeader, " ")
		if len(splitBearer) != 2 || strings.ToLower(splitBearer[0]) != "bearer" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Formato de autorización inválido"})
		}

		tokenString := strings.TrimSpace(splitBearer[1])
		miClave := []byte(os.Getenv("SECRET_JWT"))
		if len(miClave) == 0 {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Clave secreta no configurada"})
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Método de firma inesperado")
			}
			return miClave, nil
		})

		if err != nil || !token.Valid {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Token inválido"})
		}

		c.Set("user", token)
		return next(c)
	}
}
