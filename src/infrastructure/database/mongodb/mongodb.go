package mongodb

import (
	"context"
	"log"

	"quiz-go/src/infrastructure/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Client é a variável global do cliente MongoDB
var Client *mongo.Client

// QuestionCollection é a coleção de perguntas no MongoDB
var QuestionCollection *mongo.Collection

// CategoryCollection é a coleção de categorias no MongoDB
var CategoryCollection *mongo.Collection

// Connect conecta ao MongoDB e configura as coleções
func Connect() {
	clientOptions := options.Client().ApplyURI(config.GetMongoURI())
	var err error
	Client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatalf("Erro ao conectar ao MongoDB: %v", err)
	}
	log.Println("Conectado ao MongoDB com sucesso")

	// Configuração das coleções
	QuestionCollection = Client.Database("quizdb").Collection("questions")
	CategoryCollection = Client.Database("quizdb").Collection("categories")
}
