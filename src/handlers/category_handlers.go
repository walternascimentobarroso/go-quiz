package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"quiz-go/src/domain"
	"quiz-go/src/infrastructure/database/mongodb"
	"quiz-go/src/utils"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateCategory(w http.ResponseWriter, r *http.Request) {
	var category domain.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		utils.HandleError(w, err, "Erro ao decodificar JSON", http.StatusBadRequest)
		return
	}

	category.ID = primitive.NewObjectID()
	_, err := mongodb.CategoryCollection.InsertOne(context.TODO(), category)
	if err != nil {
		utils.HandleError(w, err, "Erro ao inserir categoria", http.StatusInternalServerError)
		return
	}

	utils.SendJSONResponse(w, category, http.StatusCreated)
}

func GetCategories(w http.ResponseWriter, r *http.Request) {
	var categories []domain.Category
	cursor, err := mongodb.CategoryCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		utils.HandleError(w, err, "Erro ao buscar categorias", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var category domain.Category
		if err := cursor.Decode(&category); err != nil {
			utils.HandleError(w, err, "Erro ao decodificar categoria", http.StatusInternalServerError)
			return
		}
		categories = append(categories, category)
	}

	utils.SendJSONResponse(w, categories, http.StatusOK)
}

func UpdateCategory(w http.ResponseWriter, r *http.Request) {
	categoryID := mux.Vars(r)["id"]

	var updatedCategory domain.Category
	if err := json.NewDecoder(r.Body).Decode(&updatedCategory); err != nil {
		utils.HandleError(w, err, "Erro ao decodificar JSON", http.StatusBadRequest)
		return
	}

	objectID, err := utils.ConvertID(categoryID)
	if err != nil {
		utils.HandleError(w, err, "ID inválido", http.StatusBadRequest)
		return
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": bson.M{"name": updatedCategory.Name}}

	result := mongodb.CategoryCollection.FindOneAndUpdate(context.TODO(), filter, update, options.FindOneAndUpdate().SetReturnDocument(options.After))
	if result.Err() != nil {
		utils.HandleError(w, result.Err(), "Erro ao atualizar categoria", http.StatusInternalServerError)
		return
	}

	var updatedCategoryResponse domain.Category
	if err := result.Decode(&updatedCategoryResponse); err != nil {
		utils.HandleError(w, err, "Erro ao recuperar categoria atualizada", http.StatusInternalServerError)
		return
	}

	utils.SendJSONResponse(w, updatedCategoryResponse, http.StatusOK)
}

func DeleteCategory(w http.ResponseWriter, r *http.Request) {
	categoryID := mux.Vars(r)["id"]

	objectID, err := utils.ConvertID(categoryID)
	if err != nil {
		utils.HandleError(w, err, "ID inválido", http.StatusBadRequest)
		return
	}

	filter := bson.M{"_id": objectID}
	deleteResult, err := mongodb.CategoryCollection.DeleteOne(context.TODO(), filter)
	if err != nil || deleteResult.DeletedCount == 0 {
		utils.HandleError(w, err, "Erro ao remover categoria", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
