package app

import (
	"log"

	"github.com/go-redis/redis"
)

type Cache struct {
	Client *redis.Client
}

func (c *Cache) Setup(addr string) error {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := client.Ping().Result()
	log.Printf("Cache method=setup %v\n", pong)
	log.Println("Connected to redis")
	c.Client = client
	return err
}

func NewCache(redisHost string) *Cache {
	cache := &Cache{}
	err := cache.Setup(redisHost)
	if err != nil {
		panic(err)
	}
	return cache
}
