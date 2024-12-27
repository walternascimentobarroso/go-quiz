package utils

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"quiz-go/src/domain"
	"quiz-go/src/infrastructure/database/mongodb"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func HandleError(w http.ResponseWriter, err error, message string, statusCode int) {
	log.Printf("%s: %v", message, err)
	http.Error(w, message, statusCode)
}

func SendJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Erro ao codificar resposta em JSON: %v", err)
		http.Error(w, "Erro ao retornar resposta", http.StatusInternalServerError)
	}
}

func ConvertID(id string) (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(id)
}

func GetQuestionByID(id string) (domain.Question, error) {
	var question domain.Question
	objectID, err := ConvertID(id)
	if err != nil {
		return question, err
	}

	filter := bson.M{"_id": objectID}
	err = mongodb.QuestionCollection.FindOne(context.TODO(), filter).Decode(&question)
	return question, err
}
