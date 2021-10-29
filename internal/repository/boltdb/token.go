package boltdb

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/boltdb/bolt"
	"github.com/slavasukhorukov/pocketer_bot/internal/repository"
)

type tokenRepo struct {
	db *bolt.DB
}

func NewTokenRepository(db *bolt.DB) repository.TokenRepository {
	return &tokenRepo{db: db}
}

func (repo *tokenRepo) Put(chatId int64, token string, bucket repository.Bucket) error {
	return repo.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		return b.Put(intToBytes(chatId), []byte(token))
	})
}

func (repo *tokenRepo) Get(chatId int64, bucket repository.Bucket) (string, error) {
	var token string
	if err := repo.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		token = string(b.Get(intToBytes(chatId)))
		if token == "" {
			return errors.New(fmt.Sprintf("Token with chatId = %d not found", chatId))
		}
		return nil
	}); err != nil {
		return "", err
	}

	return token, nil
}

func intToBytes(val int64) []byte {
	return []byte(strconv.FormatInt(val, 10))
}
