APP_NAME=quiz-go
DOCKER_COMPOSE=docker-compose

# Builda a aplicação
build:
	@docker build -t $(APP_NAME) .

# Sobe os serviços com Docker Compose
up:
	@$(DOCKER_COMPOSE) up -d

# Derruba os serviços
down:
	@$(DOCKER_COMPOSE) down

# Limpa os containers e imagens
clean:
	@$(DOCKER_COMPOSE) down --rmi all

restart: clean build up

# Exibe os logs
logs:
	@$(DOCKER_COMPOSE) logs -f

# Roda os testes (adicionar testes depois)
test:
	@go test ./...
