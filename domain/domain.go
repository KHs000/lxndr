package domain

import (
	"context"

	"github.com/mongodb/mongo-go-driver/mongo"
)

type (
	// Error ..
	Error struct {
		Code  int
		Error string
	}

	// Response ..
	Response struct {
		Message string   `json:"message"`
		Data    []string `json:"data"`
	}

	// User ..
	User struct {
		Email string // email
		Hash  string // hash email
		Token string // password
	}

	// Connection ..
	Connection struct {
		Ctx    context.Context
		Client *mongo.Client
	}

	// Collection ..
	Collection struct {
		Database string
		CollName string
	}
)
