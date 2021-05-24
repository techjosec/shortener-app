package redis

import (
	"fmt"
	"strconv"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	"github.com/techjosec/url-shortener-app/shortener"
)

type redisRepository struct {
	client *redis.Client
}

func newRedisClient(redisURL string) (*redis.Client, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opts)
	_, err = client.Ping().Result()
	if err != nil {
		return nil, err
	}

	return client, nil
}

func NewRedisRepository(redisURL string) (shortener.RedirectRepository, error) {

	redisRepo := &redisRepository{}
	client, err := newRedisClient(redisURL)
	if err != nil {
		return nil, errors.Wrap(err, "repository.redis.NewRedisRepository")
	}

	redisRepo.client = client

	return redisRepo, nil
}

func (rr *redisRepository) generateKey(code string) string {
	return fmt.Sprintf("redirect:%s", code)
}

func (rr *redisRepository) Find(code string) (*shortener.Redirect, error) {
	key := rr.generateKey(code)

	data, err := rr.client.HGetAll(key).Result()
	if err != nil {
		return nil, errors.Wrap(err, "repository.Redirect.redis.Find")
	}
	if len(data) == 0 {
		return nil, errors.Wrap(shortener.ErrRedirectNotFound, "repository.Redirect.redis.Find")
	}

	createdAt, err := strconv.ParseInt(data["created_at"], 10, 64)

	if err != nil {
		return nil, errors.Wrap(err, "repository.Redirect.redis.Find")
	}

	redirect := &shortener.Redirect{}
	redirect.Code = data["code"]
	redirect.URL = data["url"]
	redirect.CreatedAt = createdAt

	return redirect, nil
}

func (rr *redisRepository) Store(redirect *shortener.Redirect) error {
	key := rr.generateKey(redirect.Code)
	data := map[string]interface{}{
		"code":       redirect.Code,
		"url":        redirect.URL,
		"created_at": redirect.CreatedAt,
	}
	_, err := rr.client.HMSet(key, data).Result()
	if err != nil {
		return errors.Wrap(err, "repository.Redirect.redis.Store")
	}
	return nil
}
