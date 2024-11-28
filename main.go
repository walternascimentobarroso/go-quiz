package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Question struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Question QuestionDetails    `json:"question"`
}

type QuestionDetails struct {
	Description   string   `json:"description"`
	Explanation   string   `json:"explanation"`
	Difficulty    string   `json:"difficulty"`
	Categories    []string `json:"categories"`
	AllowMultiple bool     `json:"allow_multiple"`
	Options       []Option `json:"options"`
}

type Option struct {
	OptionText string `json:"option_text"`
	IsCorrect  bool   `json:"is_correct"`
}

var client *mongo.Client
var questionCollection *mongo.Collection

func deleteQuestion(w http.ResponseWriter, r *http.Request) {
	log.Println("Rota DELETE chamada")
	// Define path variable for question ID
	vars := mux.Vars(r)
	questionID, ok := vars["id"]
	if !ok {
		log.Println("Erro: ID da questão não fornecido")
		http.Error(w, "É necessário fornecer o ID da questão para removê-la", http.StatusBadRequest)
		return
	}

	log.Printf("Received DELETE request for question ID: %s", questionID)

	// Converte o ID da string para ObjectID
	objectID, err := primitive.ObjectIDFromHex(questionID)
	if err != nil {
		log.Printf("Erro ao converter ID da string para ObjectID: %v", err)
		http.Error(w, "ID da questão inválido", http.StatusBadRequest)
		return
	}

	// Define o filtro para a questão a ser removida
	filter := bson.M{"_id": objectID}

	// Realiza a remoção da questão
	deleteResult, err := questionCollection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Printf("Erro ao remover questão do MongoDB: %v", err)
		http.Error(w, "Erro ao remover a questão", http.StatusInternalServerError)
		return
	}

	// Verifica o número de documentos removidos
	if deleteResult.DeletedCount == 0 {
		log.Printf("Questão com ID %s não encontrada", questionID)
		http.Error(w, "Questão não encontrada", http.StatusNotFound)
		return
	}

	log.Printf("Questão com ID %s removida com sucesso", questionID)

	// Define o status HTTP 204 No Content
	w.WriteHeader(http.StatusNoContent)
}

func main() {
	// Configurações do logger
	log.SetFlags(log.LstdFlags | log.Lshortfile) // Adiciona timestamp e nome do arquivo

	// Conexão com o MongoDB
	clientOptions := options.Client().ApplyURI("mongodb://mongo:27017")
	var err error
	client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatalf("Erro ao conectar ao MongoDB: %v", err)
	}
	log.Println("Conectado ao MongoDB com sucesso")

	questionCollection = client.Database("quizdb").Collection("questions")

	r := mux.NewRouter()
	log.Println("Configurando rotas")
	r.HandleFunc("/questions/{id}", deleteQuestion).Methods("DELETE")
	r.HandleFunc("/questions", createQuestion).Methods("POST")
	r.HandleFunc("/questions", getQuestions).Methods("GET")
	r.HandleFunc("/questions/{id}", updateQuestion).Methods("PUT")

	log.Println("Iniciando servidor na porta 8000")
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func updateQuestion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // Obtém os parâmetros da rota
	id := params["id"]

	objID, err := primitive.ObjectIDFromHex(id) // Converte o ID para ObjectID do MongoDB
	if err != nil {
		log.Printf("ID inválido: %v", err)
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var updatedData QuestionDetails
	if err := json.NewDecoder(r.Body).Decode(&updatedData); err != nil {
		log.Printf("Erro ao decodificar o corpo da requisição: %v", err)
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}

	// Cria o documento de atualização
	update := bson.M{"$set": bson.M{
		"question.description":    updatedData.Description,
		"question.explanation":    updatedData.Explanation,
		"question.difficulty":     updatedData.Difficulty,
		"question.categories":     updatedData.Categories,
		"question.allow_multiple": updatedData.AllowMultiple,
		"question.options":        updatedData.Options,
	}}

	// Atualiza o documento no MongoDB
	result, err := questionCollection.UpdateOne(context.TODO(), bson.M{"_id": objID}, update)
	if err != nil {
		log.Printf("Erro ao atualizar a pergunta: %v", err)
		http.Error(w, "Erro ao atualizar a pergunta", http.StatusInternalServerError)
		return
	}

	if result.MatchedCount == 0 {
		log.Printf("Nenhuma pergunta encontrada com o ID: %s", id)
		http.Error(w, "Nenhuma pergunta encontrada", http.StatusNotFound)
		return
	}

	log.Printf("Pergunta atualizada com sucesso: %s", id)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Pergunta atualizada com sucesso"})
}

func createQuestion(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		Question QuestionDetails `json:"question"`
	}

	// Decodifica o corpo da solicitação
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		log.Printf("Erro ao decodificar JSON: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Log dos dados recebidos
	log.Printf("Recebido: %+v", requestBody)

	// Gera um novo ObjectID para a questão
	newQuestion := Question{
		ID:       primitive.NewObjectID(),
		Question: requestBody.Question, // Preenche os detalhes da questão
	}

	// Inserir a questão no MongoDB
	_, err := questionCollection.InsertOne(context.TODO(), newQuestion)
	if err != nil {
		log.Printf("Erro ao inserir questão no MongoDB: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Questão inserida com sucesso")

	// Define o cabeçalho Content-Type como application/json
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	// Retorna a questão criada como resposta em formato JSON
	if err := json.NewEncoder(w).Encode(newQuestion); err != nil {
		log.Printf("Erro ao codificar resposta em JSON: %v", err)
	}
}

func getQuestions(w http.ResponseWriter, r *http.Request) {
	var questions []Question

	// Define o cabeçalho Content-Type como application/json
	w.Header().Set("Content-Type", "application/json")

	cursor, err := questionCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Printf("Erro ao buscar perguntas: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var question Question
		if err := cursor.Decode(&question); err != nil {
			log.Printf("Erro ao decodificar pergunta: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		questions = append(questions, question)
	}

	log.Printf("Total de perguntas recuperadas: %d", len(questions))

	// Define o status HTTP 200 OK
	w.WriteHeader(http.StatusOK)

	// Retorna a lista de perguntas como resposta em formato JSON
	if err := json.NewEncoder(w).Encode(questions); err != nil {
		log.Printf("Erro ao codificar resposta em JSON: %v", err)
	}
}
