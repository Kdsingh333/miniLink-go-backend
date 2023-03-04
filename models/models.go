package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ShortenBody struct {
	LongUrl string `json:"longUrl"`
}

type CustomBody struct {
	LongUrl    string `json:"longUrl"`
	CustomCode string `json:"customCode"`
}

type UrlDoc struct {
	ID         primitive.ObjectID `bson:"_id"`
	UrlCode    string             `bson:"urlCode"`
	LongUrl    string             `bson:"longUrl"`
	ShortUrl   string             `bson:"shortUrl"`
	CreatedAt  time.Time          `bson:"createdAt"`
	ExperiesAt time.Time          `bson:"expiresAt"`
}
