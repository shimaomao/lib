package models

import (
	"fmt"
	"github.com/speedland/wcg"
	"time"
)

type ApiToken struct {
	Token       string    `json:"token"`
	Description string    `json:"desc"`
	CreatedAt   time.Time `json:"created_at"`
}

func NewApiToken() *ApiToken {
	token, err := wcg.UUID()
	if err != nil {
		panic(err)
	}
	return &ApiToken{
		Token:       token,
		Description: "",
		CreatedAt:   time.Now(),
	}
}

func (token *ApiToken) Key() string {
	return token.Token
}

func (token *ApiToken) String() string {
	return fmt.Sprintf("<ApiToken %s>", token.Token)
}
