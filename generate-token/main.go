package main

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func GenerateTokenJwt() (string, error) {
	secretKey := os.Getenv("JWT_SECRET_KEY") //deve ser setado antes de rodar o programa

	if secretKey == "" {
		return "", fmt.Errorf("JWT_SECRET_KEY not found or not set")
	}

	claims := jwt.MapClaims{
		"iss":   "my-issuer",
		"roles": []string{"admin", "superadmin"},
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	}

	// Criar token com header fixo
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token.Header["alg"] = "HS256" // Garantir que o header é idêntico
	token.Header["typ"] = "JWT"

	tokenString, err := token.SignedString([]byte(secretKey))
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
