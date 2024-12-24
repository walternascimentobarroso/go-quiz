package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	// Carregar as variáveis do arquivo .env
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Erro ao carregar o arquivo .env: %v", err)
	}
}

// GetMongoURI retorna a URI de conexão com o MongoDB. Pode ser configurada via variável de ambiente.
func GetMongoURI() string {
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		log.Fatal("MONGO_URI não configurado!")
	}
	return uri
}
