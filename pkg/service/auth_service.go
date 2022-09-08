package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"time"

	todo "github.com/RomanKlyachev/go-rest-api"
	"github.com/RomanKlyachev/go-rest-api/pkg/repository"
	"github.com/golang-jwt/jwt"
)

const (
	salt      = "jsgkaKFHQ3FKkf9am*jsf2"
	tokenTTL  = 12 * time.Hour
	signinKey = "gda7HASF9(2h^%hsdfj"
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user todo.User) (int, error) {
	user.Password = generatePasswordHash(user.Password)
	return s.repo.CreateUser(user)
}

func (s *AuthService) GenerateToken(username string, password string) (string, error) {
	user, err := s.repo.GetUser(username, generatePasswordHash(password))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{

			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
		},
		user.Id,
	})

	return token.SignedString([]byte(signinKey))
}

func (s *AuthService) ParseToten(accessToken string) (int, error) {

	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Invalid signed method")
		}

		return []byte(signinKey), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *tokenClaims")
	}

	return claims.UserId, nil

}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
