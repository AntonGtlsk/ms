package auth

import (
	"auth-ms/internal/entity/user"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JwtParser struct {
	jwtKey []byte
}

func NewJwtParser(jwtKey []byte) *JwtParser {
	return &JwtParser{jwtKey: jwtKey}
}

func (t *JwtParser) GenerateJWTToken(id float64, userId, username, avatarURL string, subscription user.Subscription) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = id
	claims["userid"] = userId
	claims["username"] = username
	claims["avatarURL"] = avatarURL
	claims["exp"] = time.Now().Add(time.Hour * 30).Unix()
	claims["subscription"] = subscription

	tokenString, err := token.SignedString(t.jwtKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (t *JwtParser) GenerateRefreshToken(id float64, userId, username, avatarURL string) (string, jwt.MapClaims, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = id
	claims["userid"] = userId
	claims["username"] = username
	claims["avatarURL"] = avatarURL
	claims["exp"] = time.Now().Add(time.Hour * 24 * 60).Unix()
	claims["createdAt"] = time.Now().Unix()

	tokenString, err := token.SignedString(t.jwtKey)

	if err != nil {
		return "", nil, err
	}

	return tokenString, claims, nil
}

// return uid, username and avatarURL
func (t *JwtParser) ParseToken(tokenString string) (float64, string, string, string, error) {

	splitString := strings.Split(tokenString, " ")

	if len(splitString) != 2 {
		return 0, "", "", "", fmt.Errorf("Failed to parse the token:")
	}

	token, err := jwt.Parse(splitString[1], func(token *jwt.Token) (interface{}, error) {
		return t.jwtKey, nil
	})

	if err != nil {
		return 0, "", "", "", fmt.Errorf("Failed to parse the token: %s", err)

	}

	if !token.Valid {
		return 0, "", "", "", fmt.Errorf("Invalid token")

	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, "", "", "", fmt.Errorf("Failed to parse token claims")
	}

	return claims["id"].(float64), claims["userid"].(string), claims["username"].(string), claims["avatarURL"].(string), nil
}

// return uid, exp and createdAt
func (t *JwtParser) ParseRefreshToken(tokenString string) (float64, string, string, string, float64, float64, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return t.jwtKey, nil
	})

	if err != nil {
		return 0, "", "", "", 0, 0, fmt.Errorf("Failed to parse the token: %s", err)
	}

	if !token.Valid {
		return 0, "", "", "", 0, 0, fmt.Errorf("Invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, "", "", "", 0, 0, fmt.Errorf("Failed to parse token claims")
	}

	return claims["id"].(float64), claims["userid"].(string), claims["username"].(string), claims["avatarURL"].(string), claims["exp"].(float64), claims["createdAt"].(float64), nil
}
