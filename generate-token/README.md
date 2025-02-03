## Gerador de Token JWT para o Severino-Plugin

Este projeto é um gerador simples de tokens JWT em Golang, projetado para testar a autenticação e autorização do Severino-Plugin no Kong API Gateway.
O token gerado pode ser usado para fazer requisições autenticadas no Kong e testar a validação do plugin Severino-Plugin.

## Como Funciona?
Este programa gera um JWT (JSON Web Token) assinado com o algoritmo HS256, contendo as seguintes informações:

* `Issuer` (iss): "my-issuer" (precisa corresponder à configuração do Kong).
* `Roles` (roles): ["admin", "superadmin"] (permite testar diferentes permissões no Severino-Plugin).
* `Expiração` (exp): Token válido por 24 horas.
* `Assinatura`: Usa a chave secreta "secretkey@123!" para assinar o token.

O JWT gerado pode ser usado para testar a autenticação e autorização no Kong API Gateway.

## Uso

* Clone o repositório (se necessário):
```bash
https://github.com/igordevopslabs/severino-plugin.git
cd generate-token
```

* Execute o projeto:
```bash
go run main.go
```

* Saída esperada:
```bash
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJteS1pc3N1ZXIiLCJyb2xlcyI6WyJhZG1pbiIsInN1cGVyYWRtaW4iXSwiZXhwIjoxNzQ1ODIwMDAwfQ.Di98sE-K3PS_0o8CByhAZhGjcHfntGyyNmSP5PQU5Lo
```

Com o token em mãos você já será capaz de interagir com o serviço cadastrado no Kong.