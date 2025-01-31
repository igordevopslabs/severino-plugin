package handler

import (
	"fmt"
	"net/http"

	"github.com/Kong/go-pdk"
	"github.com/golang-jwt/jwt/v5"
)

// Struct para representar os parametros de configuração do plugin no Kong
type Config struct {
	Issuer         string   `json:"issuer"`
	SecretKey      string   `json:"secret_key"`      // Usado para HS256
	ClaimName      string   `json:"claim_name"`      // ex: "roles"
	RequiredValues []string `json:"required_values"` // ex: ["admin", "superadmin"] etc.
	Algorithm      string   `json:"algorithm"`       // ex: "RS256", "HS256"
}

type SeverinoPlugin struct {
	Config Config
}

// Essa função recebe do Kong a config deserializada a partir do JSON e retorna uma instância do plugin
func New() interface{} {
	return &SeverinoPlugin{}
}

// Metodo Access (convenção) que será chamado pelo Kong
func (plugin *SeverinoPlugin) Access(kong *pdk.PDK) {
	//pegar o valor do header Authorization
	authHeader, err := kong.Request.GetHeader("Authorization")
	if err != nil || authHeader == "" {
		kong.Response.SetStatus(http.StatusUnauthorized)
		kong.Response.Exit(401, nil, nil)
		return
	}

	//valida se o token esta formato Bearer TOKEN
	tokenString, ok := plugin.getBearerToken(authHeader)
	if !ok {
		kong.Response.SetStatus(http.StatusUnauthorized)
		kong.Response.Exit(401, nil, nil)
		return
	}

	//valida o token JWT
	claims, validateErr := plugin.validateTokenJWT(tokenString)
	if validateErr != nil {
		kong.Response.SetStatus(http.StatusUnauthorized)
		kong.Response.Exit(401, nil, nil)
		return
	}
	//valida se possui escopo e permissão pra acessar a rota/upstream
	hasPermission := plugin.checkPermission(claims)
	if !hasPermission {
		kong.Response.SetStatus(http.StatusForbidden)
		kong.Response.Exit(403, nil, nil)
		return
	}
}

// valida se o header está no formato correto e extrai o valor do token
func (plugin *SeverinoPlugin) getBearerToken(authHeader string) (string, bool) {
	const prefix = "Bearer "
	if len(authHeader) > len(prefix) && authHeader[:len(prefix)] == prefix {
		return authHeader[len(prefix):], true
	}
	return "", false
}

func (plugin *SeverinoPlugin) validateTokenJWT(tokenString string) (jwt.MapClaims, error) {

	keyFunc := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(plugin.Config.SecretKey), nil
	}

	//parsea o token
	parsedToken, err := jwt.Parse(tokenString, keyFunc)
	if err != nil {
		return nil, fmt.Errorf("error parsing token: %v", err)
	}

	// Verifica se é válido
	if !parsedToken.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// valida se o issuer bate com o config
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	if claims["iss"] != plugin.Config.Issuer {
		return nil, fmt.Errorf("invalid issuer")
	}

	return claims, nil
}

func (plugin *SeverinoPlugin) checkPermission(claims jwt.MapClaims) bool {

	//validando se a claim name (roles) existe no token, se não, ja retorna false o que gerará um 401
	claimValue, ok := claims[plugin.Config.ClaimName]
	if !ok {
		return false
	}

	//validando se claimName é de fato um array, pois caso não, não é possivel associar diversos escopos.
	roles, ok := claimValue.([]interface{})
	if !ok {
		return false
	}

	//validação se todos os valores requiridos no array estão presentes
	required := make(map[string]bool)
	for _, val := range plugin.Config.RequiredValues {
		required[val] = false
	}

	for _, r := range roles {
		rStr, ok := r.(string)
		if !ok {
			continue
		}
		// Se esse valor está nas RequiredValues, marcamos como encontrado
		if _, found := required[rStr]; found {
			required[rStr] = true
		}
	}

	// Se algum valor não foi encontrado, retorna false
	for _, v := range required {
		if !v {
			return false
		}
	}

	return true
}
