package main

import (
	"log"
	"strconv"

	"github.com/go-redis/redis"
)

type UrlRepository interface {
	Insert(id int, url string) error
	Get(id int) (string, error)
	GetID() (int, error)
	SetID(id int) error
	DeleteLast() error
}
type Config struct {
	Addr     string
	Password string
	DB       int
}

func NewRedisDB(cfg Config) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password, // no password set
		DB:       cfg.DB,       // use default DB
	})

	return rdb, nil
}

type UrlRedis struct {
	db *redis.Client
}

func NewUrlRedis(db *redis.Client) *UrlRedis {
	urlDB := &UrlRedis{db: db}
	return urlDB
}

func (u *UrlRedis) initCounter() error {
	_, err := u.db.Get("id").Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			u.db.Set("id", 1, 0)
			return nil
		}

		return err
	}

	return nil
}

func (u *UrlRedis) Insert(id int, url string) error {
	idStr := strconv.Itoa(id)
	return u.db.Set(idStr, url, 0).Err()
}

func (u *UrlRedis) Get(id int) (string, error) {
	idStr := strconv.Itoa(id)

	url, err := u.db.Get(idStr).Result()
	if err != nil {
		log.Println(err)
		return "", err
	}

	return url, nil
}

func (u *UrlRedis) GetID() (int, error) {
	id, err := u.db.Get("id").Result()
	if err != nil {
		log.Println(err)
		return 0, err
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	return idInt, nil
}

func (u *UrlRedis) SetID(id int) error {
	return u.db.Set("id", strconv.Itoa(id), 0).Err()
}

//DeleteLast isrequired for testing
func (u *UrlRedis) DeleteLast() error {
	id, err := u.GetID()
	if err != nil {
		return err
	}

	err = u.db.Del(strconv.Itoa(id)).Err()
	if err != nil {
		return err
	}

	id--
	return u.SetID(id)
}
