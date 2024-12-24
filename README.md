```
. quiz-api/
├── main.go                      # Arquivo principal para inicializar a aplicação
├── internal/
│   ├── domain/                  # Camada de domínio
│   │   ├── question/            # Subdomínio: Questões
│   │   │   ├── entity.go        # Entidade Question
│   │   │   ├── service.go       # Regras de negócio para questões
│   │   │   ├── repository.go    # Interface para repositórios
│   │   │   ├── service_test.go  # Testes para o service de questões
│   │   ├── category/            # Subdomínio: Categorias
│   │       ├── entity.go        # Entidade Category
│   │       ├── service.go       # Regras de negócio para categorias
│   │       ├── repository.go    # Interface para repositórios
│   │       ├── service_test.go  # Testes para o service de categorias
│   ├── application/             # Casos de uso (orquestração)
│   │   ├── question/            # Casos de uso para questões
│   │   │   ├── dto.go           # DTOs para Questões
│   │   │   ├── service.go       # Orquestração de serviços relacionados a Questões
│   │   ├── category/            # Casos de uso para categorias
│   │       ├── dto.go           # DTOs para Categorias
│   │       ├── service.go       # Orquestração de serviços relacionados a Categorias
│   ├── infrastructure/          # Implementação de infraestrutura
│   │   ├── http/                # Camada de transporte HTTP
│   │   │   ├── router.go        # Configuração das rotas
│   │   │   ├── question_handler.go
│   │   │   ├── category_handler.go
│   │   ├── database/            # Implementação dos repositórios
│   │       ├── mongo_client.go  # Inicialização do cliente MongoDB
│   │       ├── question_repo.go
│   │       ├── category_repo.go
│   ├── config/                  # Configurações gerais do sistema
│       ├── config.go            # Configuração de variáveis de ambiente
└── pkg/                         # Pacotes reutilizáveis
    ├── logger/                  # Pacote de logging
    ├── validation/              # Validações customizadas
```