package models

type Auth struct {
	TokenHash string `redis:"token_hash"`
	UserID    string `redis:"user_id"`
}
