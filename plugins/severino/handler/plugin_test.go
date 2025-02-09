package handler

import (
	"reflect"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestCheckPermission(t *testing.T) {
	plugin := &SeverinoPlugin{
		Config: Config{
			ClaimName:      "roles",
			RequiredValues: []string{"admin", "user"},
		},
	}

	// Gera um token JWT com os claims adequados.
	claims := jwt.MapClaims{
		"iss":   "my-issuer",
		"roles": []string{"admin", "user"},
		"exp":   time.Now().Add(1 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Aqui, a assinatura não é relevante para o teste de checkPermission,
	// pois já estamos passando os claims diretamente.
	tokenString, err := token.SignedString([]byte("dummy"))
	if err != nil {
		t.Fatalf("Erro ao gerar token: %v", err)
	}

	// Imita a decodificação do token (ignorando a verificação de assinatura para o teste).
	parsedToken, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		t.Fatalf("Erro ao parsear token: %v", err)
	}
	parsedClaims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatalf("Erro ao converter claims")
	}

	// Verifica se checkPermission retorna true (pois "admin" está presente).
	if !plugin.checkPermission(parsedClaims) {
		t.Errorf("Esperado true para token com permissões válidas, mas obteve false")
	}
}

func TestValidateTokenJWT(t *testing.T) {
	plugin := &SeverinoPlugin{
		Config: Config{
			ClaimName:      "roles",
			RequiredValues: []string{"admin"},
			Issuer:         "my-issuer",
			Algorithm:      "HS256",
			SecretKey:      "dummy",
		},
	}

	// Gera um token JWT com os claims adequados.
	claims := jwt.MapClaims{
		"iss":   "my-issuer",
		"roles": []interface{}{"admin"},
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Aqui, a assinatura não é relevante para o teste de checkPermission,
	// pois já estamos passando os claims diretamente.
	tokenString, err := token.SignedString([]byte("dummy"))
	if err != nil {
		t.Fatalf("Erro ao gerar token: %v", err)
	}

	// Verifica se validateTokenJWT retorna os claims (pois o token é válido).
	parsedClaims, err := plugin.validateTokenJWT(tokenString)
	if err != nil {
		t.Fatalf("Erro ao validar token: %v", err)
	}

	// Verifica se os claims retornados são os mesmos que os passados.
	if parsedClaims["iss"] != claims["iss"] {
		t.Errorf("Esperado %v, mas obteve %v", claims["iss"], parsedClaims["iss"])
	}

	if !reflect.DeepEqual(parsedClaims["roles"], claims["roles"]) {
		t.Errorf("Esperado %v, mas obteve %v", claims["roles"], parsedClaims["roles"])
	}

	expectedExp := float64(claims["exp"].(int64))
	actualExp, ok := parsedClaims["exp"].(float64)

	if !ok {
		t.Fatalf("parsedClaims['exp'] não é float64")
	}
	if actualExp != expectedExp {
		t.Errorf("Esperado %v, mas obteve %v", expectedExp, actualExp)
	}
}
