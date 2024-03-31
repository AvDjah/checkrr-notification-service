package Helpers

import (
	"fmt"
	"github.com/gofiber/contrib/websocket"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"log"
)

func ParseJWT(tokenString string) (jwt.MapClaims, bool) {

	secretKey := viper.GetString("SECRET_KEY")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secretKey), nil
	})

	if err != nil {
		log.Panic(err)
		return nil, false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		fmt.Println(claims["foo"], claims["nbf"])
		return claims, true
	} else {
		fmt.Println(err)
		return nil, false
	}
}

func GetUserIdFromJWTClaim(c *websocket.Conn) int64 {

	cookie := c.Cookies("Authorization", "")
	// Get Claims
	claims, result := ParseJWT(cookie)
	if result == false {
		return 0
	}

	return int64(claims["userId"].(float64))
}
