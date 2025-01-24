package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

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
	OptionText  string `json:"option_text"`
	IsCorrect   bool   `json:"is_correct"`
	Explanation string `json:"explanation"`
}
