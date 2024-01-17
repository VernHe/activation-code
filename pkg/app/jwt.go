package app

import (
	"fmt"
	"time"

	"configuration-management/global"
	"configuration-management/pkg/logger"

	"github.com/golang-jwt/jwt/v5"
)

const (
	TokenExp = time.Hour * 5
)

var (
	secretKey = []byte("secret-key") // TODO: 从环境变量中读取
)

type UserInfo struct {
	UserId   string   `json:"userId"`
	Username string   `json:"username"`
	MaxCnt   int      `json:"maxCnt"`
	Roles    []string `json:"roles"`
}

func (u *UserInfo) IsRoot() bool {
	for _, role := range u.Roles {
		if role == "root" {
			return true
		}
	}
	return false
}

type CreateTokenArgs struct {
	UserId   string
	Username string
	Roles    []string
}

func CreateToken(info UserInfo) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"userId":   info.UserId,
			"username": info.Username,
			"maxCnt":   info.MaxCnt,
			"roles":    info.Roles,
			"exp":      time.Now().Add(TokenExp).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}

func GetUserInfoFromToken(tokenString string) (UserInfo, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		global.Logger.WithFields(logger.Fields{
			"tokenString": tokenString,
		}).Error("parse token failed", err)
		return UserInfo{}, err
	}

	if !token.Valid {
		global.Logger.WithFields(logger.Fields{
			"tokenString": tokenString,
		}).Error("invalid token")
		return UserInfo{}, fmt.Errorf("invalid token")
	}

	var userInfo UserInfo
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if claims["roles"] == nil {
			global.Logger.WithFields(logger.Fields{
				"tokenString": tokenString,
			}).Error("invalid token claims")
			return UserInfo{}, fmt.Errorf("invalid token claims")
		}

		userInfo.UserId = claims["userId"].(string)
		userInfo.Username = claims["username"].(string)
		userInfo.MaxCnt = int(claims["maxCnt"].(float64))
		if roles, ok := claims["roles"].([]interface{}); ok {
			for _, role := range roles {
				if r, ok := role.(string); ok {
					userInfo.Roles = append(userInfo.Roles, r)
				}
			}
		}
		return userInfo, nil
	}

	global.Logger.WithFields(logger.Fields{
		"tokenString": tokenString,
	}).Error("invalid token claims")
	return UserInfo{}, fmt.Errorf("invalid token claims")
}
