package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Kong/go-pdk"
	"github.com/golang-jwt/jwt/v5"
	"github.com/igordevopslabs/severino-plugin/pkg"
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
		kong.Response.Exit(401, []byte("token header empty or not valid"), nil)
		pkg.LogError("token header empty or not valid", fmt.Errorf("token header empty or not valid"))
		return
	}

	//valida se o token esta formato Bearer TOKEN
	tokenString, ok := plugin.getBearerToken(authHeader)
	if !ok {
		kong.Response.SetStatus(http.StatusUnauthorized)
		kong.Response.Exit(401, []byte("token header is incorrect"), nil)
		pkg.LogError("token header is incorrect", fmt.Errorf("token header is incorrect"))
		return
	}

	//valida o token JWT
	claims, validateErr := plugin.validateTokenJWT(tokenString)
	if validateErr != nil {
		kong.Response.SetStatus(http.StatusUnauthorized)
		kong.Response.Exit(401, []byte("token header is not valid or expired"), nil)
		pkg.LogError("token header is not valid or expired", validateErr)
		return
	}
	//valida se possui escopo e permissão pra acessar a rota/upstream
	hasPermission := plugin.checkPermission(claims)
	if !hasPermission {
		kong.Response.SetStatus(http.StatusForbidden)
		kong.Response.Exit(403, []byte("token has no permission"), nil)
		pkg.LogError("token has no permission", fmt.Errorf("token has no permission"))
		return
	}
}

// valida se o header está no formato correto e extrai o valor do token
func (plugin *SeverinoPlugin) getBearerToken(authHeader string) (string, bool) {
	const prefix = "Bearer "
	if len(authHeader) > len(prefix) && authHeader[:len(prefix)] == prefix {
		pkg.LogInfo("token found")
		return authHeader[len(prefix):], true
	}
	pkg.LogError("token not found", fmt.Errorf("token not found"))
	return "", false
}

func (plugin *SeverinoPlugin) validateTokenJWT(tokenString string) (jwt.MapClaims, error) {

	keyFunc := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			err := fmt.Errorf("unexpected signing method")
			pkg.LogError("unexpected signing method", err)
			return nil, err
		}
		pkg.LogInfo("token validated")
		return []byte(plugin.Config.SecretKey), nil
	}

	//parsea o token
	parsedToken, err := jwt.Parse(tokenString, keyFunc)
	if err != nil {
		err := fmt.Errorf("error parsing token: %v", err)
		pkg.LogError("error parsing token", err)
		return nil, err
	}

	// Verifica se é válido
	if !parsedToken.Valid {
		err := fmt.Errorf("invalid token")
		pkg.LogError("invalid token", err)
		return nil, err
	}

	// valida se o issuer bate com o config
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		err := fmt.Errorf("invalid token claims")
		pkg.LogError("invalid token claims", err)
		return nil, err
	}
	//valida se o issuer bate com o config
	if claims["iss"] != plugin.Config.Issuer {
		err := fmt.Errorf("invalid issuer")
		pkg.LogError("invalid issuer", err)
		return nil, err
	}

	// valida se o token não está expirado
	if float64(time.Now().Unix()) > claims["exp"].(float64) {
		err := fmt.Errorf("token expired")
		pkg.LogError("token expired", err)
		return nil, err
	}

	pkg.LogInfo("all token validated")
	return claims, nil
}

func (plugin *SeverinoPlugin) checkPermission(claims jwt.MapClaims) bool {

	//validando se a claim name (roles) existe no token, se não, ja retorna false o que gerará um 401
	claimValue, ok := claims[plugin.Config.ClaimName]
	if !ok {
		pkg.LogError("claim not found", fmt.Errorf("claim not found"))
		return false
	}

	//validando se claimName é de fato um array, pois caso não, não é possivel associar diversos escopos.
	roles, ok := claimValue.([]interface{})
	if !ok {
		pkg.LogError("claim not an array", fmt.Errorf("claim not an array"))
		return false
	}

	//validação se todos os valores requiridos no array estão presentes
	required := make(map[string]bool)
	for _, val := range plugin.Config.RequiredValues {
		pkg.LogError("required value", fmt.Errorf("%s", val))
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
			pkg.LogError("required value not found", fmt.Errorf("required value not found"))
			return false
		}
	}

	pkg.LogInfo("all permissions validated")
	return true
}
