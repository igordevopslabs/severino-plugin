package main

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func GenerateTokenJwt() (string, error) {
	SECRET_KEY := "secretkey@123!"

	claims := jwt.MapClaims{
		"iss":   "my-issuer",
		"roles": []string{"admin", "superan"},
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	}

	// Criar token com header fixo
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token.Header["alg"] = "HS256" // Garantir que o header é idêntico
	token.Header["typ"] = "JWT"

	tokenString, err := token.SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func main() {
	token, err := GenerateTokenJwt()
	if err != nil {
		fmt.Println("Erro ao gerar token:", err)
		return
	}
	fmt.Println(token)
}
