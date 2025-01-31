# severino-plugin
Um plugin simples de autorização e autenticação JWT para Kong API Gateway


## Ciclo de vida de um plugin para o Kong Api Gateway
* O Kong intercepta as requisições antes de chegar ao seu upstream (serviço de destino).
* Em cada fase (por exemplo, access, header_filter, response, etc.), o Kong chama o método correspondente do seu plugin para executar alguma lógica.
* A validação JWT, será realizada na fase de Access, pois será necessário inspecionar a request para verificar a existência de um header específco.


## Testando a requisição autenticada
Para realizar o teste no serviço cadastrado no kong, faça as seguintes requests:

### Testando unauthorizer (401)
Faça uma requisição com um valor de token inválido.
```bash
curl -i -H "Authorization: Bearer xpto" http://localhost:8000/api/todos/1
```

Saída esperada:
```bash
HTTP/1.1 401 Unauthorized
```
### Testando Forbiden (403)
Faça uma requisição com um token JWT válido, porém, com o escopo não permitido.

* utilize a automação [generate-token](./generate-token/main.go)para gerar um novo token JWT, mas altere o escopo de roles para:
```diff
- "roles": []string{"admin", "superadmin"},
+ "roles": []string{"admin", "batata"},
```

Isso irá alterar o escopo de permissão do claim do token JWT. 

```bash
curl -i -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJteS1pc3N1ZXIiLCJyb2xlcyI6WyJhZG1pbiIsImJhdGF0YSJdfQ.QznRmiKR0SNd5m2r_fB1guVdpZj_HA2FX1b4C3rLFhw" http://localhost:8000/api/todos/1
```

Saída esperada:
```bash
HTTP/1.1 403 Forbidden
```
### Testando Ok (200)
Faça uma requisição com um token JWT válido

* utilize a automação [generate-token](./generate-token/main.go)para gerar um novo token JWT.
* Certifique-se que o escopo do claim é: ```{"admin", "superadmin"}```

Saída esperada:
```bash
HTTP/1.1 200 OK
```

