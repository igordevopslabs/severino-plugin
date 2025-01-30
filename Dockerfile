FROM golang:1.22 AS builder

WORKDIR /go/src/severino-plugin

COPY ./plugins/severino .

RUN go build -o severino-plugin .

# Usar a imagem oficial do Kong como base
FROM kong:3

# Copiar o arquivo de plugins
COPY plugins.yaml /plugins.yaml


# Instalar plugins listados no arquivo YAML
RUN yq e '.plugins[] | "luarocks install \(.name) \(.version)"' /plugins.yaml | sh

# Configurar permissões para os diretórios do Kong
USER root
RUN mkdir -p /usr/local/kong/logs /usr/local/kong/sockets /usr/local/kong/pids \
    && chown -R 1000:1000 /usr/local/kong

COPY --from=builder /go/src/severino-plugin/severino-plugin /usr/local/bin/severino-plugin

# Copiar e configurar script de inicialização
COPY ./scripts/init-kong.sh /init-kong.sh
RUN chmod +x /init-kong.sh

# Voltar para o usuário padrão do Kong
USER 1000

# Copiar as configurações do Kong
COPY ./config /etc/kong/config

# Definir o script como ponto de entrada
ENTRYPOINT ["/init-kong.sh"]
