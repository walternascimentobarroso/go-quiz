# Usa a imagem base do Go
FROM golang:1.20-alpine

# Define o diretório de trabalho
WORKDIR /app

# Copia o arquivo go.mod
COPY go.mod ./

# Faz o tidy para garantir que as dependências estejam corretas
RUN go mod tidy

# Copia o código-fonte
COPY . .

# Baixa as dependências necessárias
RUN go get -d ./...

# Faz o build da aplicação
RUN go build -o /quiz-go

# Expõe a porta da aplicação
EXPOSE 8000

# Comando para rodar a aplicação
CMD ["/quiz-go"]
