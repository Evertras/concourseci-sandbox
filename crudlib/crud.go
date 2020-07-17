package crudlib

import (
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
)

const hashKeyAge = "age"
const hashKeyColor = "color"

type Client struct {
	pool *redis.Pool
}

func NewClient(pool *redis.Pool) *Client {
	return &Client{
		pool,
	}
}

func TwoPlusTwo() int {
	return 4
}

func dogKey(name string) string {
	return fmt.Sprintf("dog:%s", name)
}

func (c *Client) CreateDog(name string, color string, age time.Duration) error {
	dbClient := c.pool.Get()
	key := dogKey(name)

	dbClient.Send("MULTI")
	dbClient.Send("HSET", key, hashKeyAge, age.String())
	dbClient.Send("HSET", key, hashKeyColor, color)

	if _, err := dbClient.Do("EXEC"); err != nil {
		return fmt.Errorf("dbClient.Do(): %w", err)
	}

	return nil
}

func (c *Client) GetDogAge(name string) (time.Duration, error) {
	dbClient := c.pool.Get()
	key := dogKey(name)

	ageString, err := redis.String(dbClient.Do("HGET", key, hashKeyAge))

	if err != nil {
		return 0, fmt.Errorf("dbClient.Do(): %w", err)
	}

	age, err := time.ParseDuration(ageString)

	if err != nil {
		return 0, fmt.Errorf("time.ParseDuration: %w", err)
	}

	return age, nil
}
