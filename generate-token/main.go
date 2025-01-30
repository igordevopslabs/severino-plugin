package main

import (
	"fmt"

	"github.com/golang-jwt/jwt"
)

var (
	JWT_SECRET = "mysecret@12345"
)

func main() {
	token, err := GenerateTokenJwt()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("TOKEN: %s\n", token)

}

func GenerateTokenJwt() (string, error) {
	// Create a new token object, specifying signing method and the claims
	claims := jwt.MapClaims{
		"iss":   "my-issuer",
		"roles": []string{"admin"},
		//"exp":   time.Now().Add(time.Hour * 1).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(JWT_SECRET))

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return tokenString, nil
}
