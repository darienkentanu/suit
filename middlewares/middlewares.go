package middlewares

import (
	"fmt"
	"time"
	"github.com/darienkentanu/suit/constants"

	jwt "github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

var IsLoggedIn = echoMiddleware.JWTWithConfig(echoMiddleware.JWTConfig{
	SigningKey: []byte(constants.JWT_SECRET),
})

func CreateToken(id int, role string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["id"] = id
	claims["role"] = role
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix() // Token expires after 1 hour

	tokenString, err := token.SignedString([]byte(constants.JWT_SECRET))
	if err != nil {
		fmt.Printf("Something Went Wrong: %s", err.Error())
		return "", err
	}
	return tokenString, nil
}

func IsStaff(next echo.HandlerFunc) echo.HandlerFunc {
	return func(e echo.Context) error {
		token := e.Get("user").(*jwt.Token)
		if token.Valid {
			claims := token.Claims.(jwt.MapClaims)
			role := claims["role"].(string)
			if role != "staff" {
				return echo.ErrUnauthorized
			}
		}
		return next(e)
	}
}

func IsUser(next echo.HandlerFunc) echo.HandlerFunc {
	return func(e echo.Context) error {
		token := e.Get("user").(*jwt.Token)
		if token.Valid {
			claims := token.Claims.(jwt.MapClaims)
			role := claims["role"].(string)
			if role != "user" {
				return echo.ErrUnauthorized
			}
		}
		return next(e)
	}
}

func CurrentLoginUser(e echo.Context) int {
	token := e.Get("user").(*jwt.Token)
	if token != nil && token.Valid {
		claims := token.Claims.(jwt.MapClaims)
		id := claims["id"]
		switch id.(type) {
		case float64:
			return int(id.(float64))
		default:
			return id.(int)
		}
	}
	return -1 // invalid user
}

func CurrentRoleLoginUser(e echo.Context) string {
	token := e.Get("user").(*jwt.Token)
	if token != nil && token.Valid {
		claims := token.Claims.(jwt.MapClaims)
		role := claims["role"]
		if role == "staff" {
			return "staff"
		} else {
			return "user"
		}
	}
	return "" // invalid
}
