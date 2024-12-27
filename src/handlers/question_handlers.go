package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"quiz-go/src/domain"
	"quiz-go/src/infrastructure/database/mongodb"
	"quiz-go/src/utils"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetQuestionByID(id string) (domain.Question, error) {
	var question domain.Question
	objectID, err := utils.ConvertID(id)
	if err != nil {
		return question, err
	}

	filter := bson.M{"_id": objectID}
	err = mongodb.QuestionCollection.FindOne(context.TODO(), filter).Decode(&question)
	return question, err
}

func GetQuestion(w http.ResponseWriter, r *http.Request) {
	questionID := mux.Vars(r)["id"]

	question, err := GetQuestionByID(questionID)
	if err != nil {
		utils.HandleError(w, err, "Erro ao buscar a questão", http.StatusNotFound)
		return
	}

	log.Printf("Get Questão com ID %s", questionID)
	utils.SendJSONResponse(w, question, http.StatusOK)
}

func DeleteQuestion(w http.ResponseWriter, r *http.Request) {
	questionID := mux.Vars(r)["id"]

	question, err := GetQuestionByID(questionID)
	if err != nil {
		utils.HandleError(w, err, "Erro ao buscar a questão", http.StatusNotFound)
		return
	}

	// Remove the question
	filter := bson.M{"_id": question.ID}
	deleteResult, err := mongodb.QuestionCollection.DeleteOne(context.TODO(), filter)
	if err != nil || deleteResult.DeletedCount == 0 {
		utils.HandleError(w, err, "Erro ao remover a questão", http.StatusInternalServerError)
		return
	}

	log.Printf("Get Questão com ID %s", questionID)
	w.WriteHeader(http.StatusNoContent)
}

func UpdateQuestion(w http.ResponseWriter, r *http.Request) {
	questionID := mux.Vars(r)["id"]

	var updatedDetails domain.QuestionDetails
	if err := json.NewDecoder(r.Body).Decode(&updatedDetails); err != nil {
		utils.HandleError(w, err, "Erro no corpo da requisição", http.StatusBadRequest)
		return
	}

	objectID, err := utils.ConvertID(questionID)
	if err != nil {
		utils.HandleError(w, err, "ID inválido", http.StatusBadRequest)
		return
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": bson.M{"question": updatedDetails}}

	result := mongodb.QuestionCollection.FindOneAndUpdate(context.TODO(), filter, update, options.FindOneAndUpdate().SetReturnDocument(options.After))
	if result.Err() != nil {
		utils.HandleError(w, result.Err(), "Erro ao atualizar questão", http.StatusInternalServerError)
		return
	}

	var updatedQuestion domain.Question
	if err := result.Decode(&updatedQuestion); err != nil {
		utils.HandleError(w, err, "Erro ao recuperar questão atualizada", http.StatusInternalServerError)
		return
	}

	utils.SendJSONResponse(w, updatedQuestion, http.StatusOK)
}

func CreateQuestion(w http.ResponseWriter, r *http.Request) {
	var questionDetails domain.QuestionDetails
	if err := json.NewDecoder(r.Body).Decode(&questionDetails); err != nil {
		utils.HandleError(w, err, "Erro ao decodificar JSON", http.StatusBadRequest)
		return
	}

	newQuestion := domain.Question{
		ID:       primitive.NewObjectID(),
		Question: questionDetails,
	}

	_, err := mongodb.QuestionCollection.InsertOne(context.TODO(), newQuestion)
	if err != nil {
		utils.HandleError(w, err, "Erro ao inserir questão no MongoDB", http.StatusInternalServerError)
		return
	}

	utils.SendJSONResponse(w, newQuestion, http.StatusCreated)
}

func GetQuestions(w http.ResponseWriter, r *http.Request) {
	// Obter parâmetros de query
	categoryFilter := r.URL.Query().Get("category")
	randomize := r.URL.Query().Get("random") == "true"
	limitParam := r.URL.Query().Get("limit")

	// Converter limit para um número, se possível
	var limit int64 = 0 // Se limit for 0, significa sem limite (todos os itens)
	if limitParam != "" {
		var err error
		limit, err = strconv.ParseInt(limitParam, 10, 64)
		if err != nil || limit <= 0 {
			// Se o parâmetro limit não for válido, podemos retornar erro ou usar o padrão
			utils.HandleError(w, err, "Valor inválido para 'limit'", http.StatusBadRequest)
			return
		}
	}

	// Construir filtro para consulta
	filter := bson.M{}
	if categoryFilter != "" {
		filter = bson.M{"question.categories": categoryFilter}
	}

	var pipeline []bson.M
	pipeline = append(pipeline, bson.M{"$match": filter})
	if randomize {
		pipeline = append(pipeline, bson.M{"$sample": bson.M{"size": limit}})
	}

	if !randomize && limit > 0 {
		pipeline = append(pipeline, bson.M{"$limit": limit})
	}

	// Agregar perguntas no MongoDB
	cursor, err := mongodb.QuestionCollection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		utils.HandleError(w, err, "Erro ao buscar perguntas", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.TODO())

	// Decodificar resultados
	var questions []domain.Question
	for cursor.Next(context.TODO()) {
		var question domain.Question
		if err := cursor.Decode(&question); err != nil {
			utils.HandleError(w, err, "Erro ao decodificar pergunta", http.StatusInternalServerError)
			return
		}
		questions = append(questions, question)
	}

	// Verificar se há resultados
	if len(questions) == 0 {
		utils.HandleError(w, nil, "Nenhuma pergunta encontrada", http.StatusNotFound)
		return
	}

	// Enviar resposta JSON
	utils.SendJSONResponse(w, questions, http.StatusOK)
}
