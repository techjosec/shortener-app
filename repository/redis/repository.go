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

func fillRedirect(redisData map[string]string) (*shortener.Redirect, error) {

	createdAt, err := strconv.ParseInt(redisData["created_at"], 10, 64)

	if err != nil {
		return nil, err
	}

	redirect := shortener.Redirect{}
	redirect.Code = redisData["code"]
	redirect.URL = redisData["url"]
	redirect.CreatedAt = createdAt

	return &redirect, nil
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

	redirect, err := fillRedirect(data)
	if err != nil {
		return nil, errors.Wrap(err, "repository.Redirect.redis.fillRedirect")
	}

	return redirect, nil
}

func (rr *redisRepository) ListAll() (*[]shortener.Redirect, error) {

	var redirects []shortener.Redirect

	index := 0
	iterator := rr.client.Scan(0, "redirect:*", 0).Iterator()

	for iterator.Next() {

		data, err := rr.client.HGetAll(iterator.Val()).Result()
		if err != nil {
			return nil, errors.Wrap(err, "repository.Redirect.redis.FindAll")
		}
		if len(data) == 0 {
			return nil, errors.Wrap(shortener.ErrRedirectNotFound, "repository.Redirect.redis.FindAll")
		}

		redirect, err := fillRedirect(data)
		if err != nil {
			return nil, errors.Wrap(err, "repository.Redirect.redis.fillRedirect")
		}

		redirects = append(redirects, *redirect)
		index = index + 1
	}

	if err := iterator.Err(); err != nil {
		return nil, errors.Wrap(err, "repository.Redirect.redis.FindAll")
	}

	if len(redirects) == 0 {
		return nil, errors.Wrap(shortener.ErrRedirectNotFound, "repository.Redirect.redis.FindAll")
	}

	return &redirects, nil
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
