package middleware

import (
	"errors"
	"event-booking-api/app/constant"
	"event-booking-api/app/pkg"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Auth(c *gin.Context) {
	defer pkg.PanicHandler(c)

	authHeader := c.Request.Header.Get("Authorization")
	if authHeader == "" {
		pkg.PanicException(constant.Unauthorized)
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		pkg.PanicException(constant.Unauthorized)
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("unexpected signing method")
		}

		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})
	if err != nil {
		pkg.PanicException(constant.Unauthorized)
	}

	if !parsedToken.Valid {
		pkg.PanicException(constant.Unauthorized)
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		pkg.PanicException(constant.Unauthorized)
	}

	userId := int(claims["user_id"].(float64))
	roleId := int(claims["role_id"].(float64))

	c.Set("userId", userId)
	c.Set("roleId", roleId)

	c.Next()
}
