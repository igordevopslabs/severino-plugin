.PHONY: run build install-tools go-checks

GOBIN := $(shell go env GOPATH)/bin

run:
	@echo "==> Iniciando a aplicação..."
	docker-compose up -d
	@echo "==> waiting app up..."

build:
	@echo "==> Building app"
	docker-compose up --build

install-tools: 
	@echo "==> Installing gotest"
	@go install github.com/rakyll/gotest@latest
	@echo "==> Installing staticcheck"
	@go install honnef.co/go/tools/cmd/staticcheck@latest
	@echo "==> Installing govulncheck"
	@go install golang.org/x/vuln/cmd/govulncheck@latest
	@echo "==> Tools installed in $(GOBIN)"

go-checks:
	@echo "Rodando validações de segurança no codigo"
	@echo "==> Running staticcheck"
	@cd plugins/severino && $(GOBIN)/staticcheck ./...
	@cd generate-token && $(GOBIN)/staticcheck ./...
	@echo "==> Running govulncheck"
	@cd plugins/severino && $(GOBIN)/govulncheck ./...
	@cd generate-token && $(GOBIN)/govulncheck ./...
	@echo "==> Running unit tests"
	@cd plugins/severino && go test -v -cover ./...

