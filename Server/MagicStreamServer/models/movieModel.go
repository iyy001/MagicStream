package models

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Genre struct {
	GenreID   int    `bson:"genre_id" json:"genre_id" validate:"required"`
	GenreName string `bson:"genre_name" json:"genre_name" validate:"required,min=2,max=100"`
}
type Ranking struct {
	RankingValue int    `bson:"ranking_value" json:"ranking_value" validate:"required,min=1,max=5"`
	RankingName  string `bson:"ranking_name" json:"ranking_name" validate:"oneof=Excellent Good Okay Bad Terrible"`
}

type Movie struct {
	ID          bson.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ImdbID      string        `json:"imdb_id" bson:"imdb_id" validate:"required"`
	Title       string        `json:"title" bson:"title" validate:"required, min=2,max=500"`
	PosterPath  string        `json:"poster_path" bson:"poster_path"`
	YoutubeID   string        `json:"youtube_id" bson:"youtube_id"`
	Genre       []Genre       `json:"genre" bson:"genre"`
	AdminReview string        `json:"admin_review" bson:"admin_review"`
	Ranking     Ranking       `json:"ranking" bson:"ranking"`
}
