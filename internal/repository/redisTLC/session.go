package redisTLC

import (
	"errors"
	"fmt"

	"mainService/pkg/serverErrors"

	"github.com/gomodule/redigo/redis"
)

type IAuthRepository interface {
	AddSession(sessionID string, userID string) error
	DeleteSession(sessionID string) error
	ValidateSession(sessionID string) error
	GetUserIdBySession(sessionID string) (string, error)
}

type redisAuthRepository struct {
	sessionStorage *redis.Pool
}

func NewRedisAuthRepository(conn *redis.Pool) IAuthRepository {
	return &redisAuthRepository{
		sessionStorage: conn,
	}
}

func (p *redisAuthRepository) AddSession(sessionID string, userID string) error {
	connection := p.sessionStorage.Get()
	defer connection.Close()

	sessionKey := "sessions:" + sessionID

	result, err := redis.String(connection.Do("SET", sessionKey, userID))
	if err != nil {
		return serverErrors.INTERNAL_SERVER_ERROR
	} else if result != "OK" {
		return fmt.Errorf(result)
	}

	return nil
}

func (p *redisAuthRepository) DeleteSession(sessionID string) error {
	connection := p.sessionStorage.Get()
	defer connection.Close()

	sessionKey := "sessions:" + sessionID

	_, err := redis.Int(connection.Do("DEL", sessionKey))
	if err != nil {
		return serverErrors.INTERNAL_SERVER_ERROR
	}

	return nil
}

func (p *redisAuthRepository) ValidateSession(sessionID string) error {
	connection := p.sessionStorage.Get()
	defer connection.Close()

	sessionKey := "sessions:" + sessionID

	result, err := redis.Int(connection.Do("EXISTS", sessionKey))
	if result == 0 {
		return SESSION_NOT_FOUND
	} else if err != nil {
		return serverErrors.INTERNAL_SERVER_ERROR
	}

	return nil
}

func (p *redisAuthRepository) GetUserIdBySession(sessionID string) (string, error) {
	connection := p.sessionStorage.Get()
	defer connection.Close()

	sessionKey := "sessions:" + sessionID

	userID, err := redis.String(connection.Do("GET", sessionKey))
	if errors.Is(err, redis.ErrNil) {
		return "", SESSION_NOT_FOUND
	}
	if err != nil {
		return "", serverErrors.INTERNAL_SERVER_ERROR
	}

	return userID, nil
}
