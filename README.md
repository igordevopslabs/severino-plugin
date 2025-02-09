# severino-plugin
![Golang Scans](https://github.com/igordevopslabs/severino-plugin/actions/workflows/go-checks.yml/badge.svg)

O Severino-Plugin é um plugin customizado para o Kong API Gateway, desenvolvido em Golang, que valida tokens JWT (JSON Web Token) e controla o acesso às rotas com base nas permissões definidas.

Ele garante que:

1. Somente requisições autenticadas com um JWT válido sejam aceitas.
2. A assinatura do JWT seja validada com base na chave secreta configurada.
3. As permissões do usuário (roles) sejam verificadas antes de permitir o acesso.

## Como funciona a autenticação e autorização?
* O plugin extrai o token JWT do cabeçalho Authorization: `Bearer <TOKEN>.`
* Ele valida a assinatura do token (chave secreta para HS256).
* Ele verifica se o usuário tem as permissões necessárias (baseado na claim "roles" do JWT).

Se todas as condições forem atendidas, a requisição é permitida. Caso contrário, retorna 401 (Unauthorized) ou 403 (Forbidden).

## Setup do Ambiente
Para podermos testar o plugin, será necessário subir uma stack do Kong, com controlplane e dataplane, além da interface gráfica do Kong Manager para que possamos criar os services e definir os plugins.

Você será capaz de subir toda essa stack apenas executando o comando:
```bash
docker-compose up --build
```

Isso irá subir toda a stack do kong, além da sua database.

**OBS:** Pode ser que seja necessário realizar alguns ajustes no Dockerfile devido a particularidades do SO em que esta sendo executado, por exemplo, se você for testar a partir de MacOs, ajustes nos sockets de rede podem ser necessários.

Uma vez que você tenha inicializado todos os serviços usando o Docker Compose, várias interfaces estarão disponíveis para interação e administração do Kong. A seguir estão as mais relevantes:

* `Kong Admin API`: Acessível via `http://localhost:8001`, esta é a API administrativa que você usará para todas as operações programáticas no Kong. Seja para adicionar novos serviços, rotas ou para ativar plugins, tudo pode ser feito aqui através de chamadas RESTful.

* `Kong Admin GUI`: Esta é a interface gráfica do usuário e pode ser acessada através de `http://localhost:8002`. A GUI oferece uma maneira mais visual e interativa para gerenciar os aspectos do seu gateway API. Ele é especialmente útil para visualizar o fluxo de tráfego, modificar configurações existentes ou adicionar novas funcionalidades de forma mais intuitiva.

## Troubleshooting
Se você encontrar problemas ao usar este ambiente Kong, a primeira coisa a fazer é conferir os logs dos contêineres em questão. Isso pode ser feito usando `docker logs` <container_name>. Além disso, as interfaces de administração podem fornecer informações úteis sobre a configuração atual e o estado dos plugins e serviços. Se você suspeitar de problemas com plugins específicos, verifique se eles estão corretamente listados e ativados tanto na Admin API quanto na GUI. Certifique-se também de que a estrutura de pastas e os arquivos YAML estão corretamente formatados e localizados nos lugares corretos.

## Setando um service e rota via console do Kong Manager:

#### Definindo o Service Gateway
1. No menu lateral esquerdo, escolha a opção **Gateway Services**
2. **New Gateway Service**
3. Defina um nome para seu service, por exemplo, ```tasks```
4. No menu abaixo, em Service Endpoint, escolha a opção **Protocol, Host, Port and Path** e informe os valores de acordo. Aqui vai a sugestão para usar uma [fake API](https://jsonplaceholder.typicode.com/) já publicada na internet, isso facilitará os testes e elimina a necessidade de subir uma API para testar o plugin.

#### Definindo a Rota
Após a definição do Service, vamos à definição da Rota.

1. Clique no serviço existente que deseja criar uma rota.
2. Na guia de **Routes** clique em **+New Routes**.
3. De um nome a rota, de preferência algo que tenha alguma ligação com o serviço.
4. Selecione os protocolos adequados, para a maioria dos exemplos **HTTP/HTTPS** irá servir.
5. Defina um path (ou rota) para o seu serviço, por exemplo: ```/api```.
6. **OPCIONAL**: Você pode personalizar os métodos HTTP permitidos para a rota, por exemplo, `GET`, `POST`, `PATCH`, dentre outros.

## Habilitando o Plugin
Com o service e route definidos, agora é necessário habilitar o plugin.

`SUGESTÃO`: Habilitar o plugin individualmente por rota, isso evita que outros serviços/rotas não sejam afetados por algum mau funcionamento ou comportamento indesejado do plugin.

**Configuração do Plugin no Kong**
- O plugin pode ser ativado em uma rota específica, um serviço ou globalmente no Kong.
- A configuração é feita via Kong Manager (interface gráfica) ou Admin API.

**Configuração via Kong Manager**:
1. Acesse o Kong Manager (http://localhost:8002).
2. Vá até **Routes**, selecione a rota desejada.
3. No menu Plugins, clique em Add Plugin e escolha `severino-plugin`.
4. Configure os campos necessários:
    * `issuer` → "my-issuer" (deve ser o mesmo no JWT).
    * `algorithm` → "HS256".
    * `secret_key` → Para validar a assinatura do JWT.
    * `claim_name` → "roles" (onde as permissões estão armazenadas no JWT).
    * `required_values` → Exemplo: ["admin", "superadmin"] (define quem pode acessar a API).
5. Clique em Save para ativar o plugin.

## Testes e Validações

#### Testando a requisição autenticada
Para realizar o teste no serviço cadastrado no kong, faça as seguintes requests:


#### Testando requisição sem token (401)
Faça uma requisição sem token no header.
```bash
curl -i http://localhost:8000/api/todos/1
```

Saída esperada:
```bash
HTTP/1.1 401 Unauthorized
token token header empty or not valid
```

#### Testando unauthorizer (401)
Faça uma requisição com um valor de token inválido.
```bash
curl -i -H "Authorization: Bearer xpto" http://localhost:8000/api/todos/1
```

Saída esperada:
```bash
HTTP/1.1 401 Unauthorized
token header is not valid or expired
```
#### Testando Forbiden (403)
Faça uma requisição com um token JWT válido, porém, com o escopo não permitido.

* utilize a automação [generate-token](./generate-token/README.md) para gerar um novo token JWT, mas altere o escopo de roles para:
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
token has no permission
```
#### Testando Ok (200)
Faça uma requisição com um token JWT válido

* utilize a automação [generate-token](./generate-token/README.md) para gerar um novo token JWT.
* Certifique-se que o escopo do claim é: ```{"admin", "superadmin"}```

Saída esperada:
```bash
HTTP/1.1 200 OK
X-Kong-Upstream-Latency: 68
X-Kong-Proxy-Latency: 1
X-Kong-Request-Id: 4e95ae53b62681644db4596e640a0f71

{
  "userId": 1,
  "id": 1,
  "title": "delectus aut autem",
  "completed": false
}
```

## Contribuições
Contribuições são bem-vindas! Abra uma issue ou um Pull Request.

## Licença
Este projeto está sob a licença MIT. Consulte o arquivo [LICENSE.md](LICENSE.md) para mais detalhes.

