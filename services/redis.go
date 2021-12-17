package services

import (
	"fmt"
	"time"

	r "github.com/go-redis/redis"
)

const (
	prefix = "goServer"
	suffix = ":"
)

type Redis struct {
	Client *r.Client
}

func parsePrefix(key string) string {
	return fmt.Sprintf("%s%s%s", prefix, suffix, key)
}

func RConnect() *Redis {
	client := r.NewClient(&r.Options{
		Addr: "localhost:6379",
	})
	return &Redis{Client: client}
}

func (r *Redis) Ping() {
	pong, err := r.Client.Ping().Result()
	fmt.Println(pong, err)
}

func (r *Redis) Get(key string) (string, error) {
	key = parsePrefix(key)
	val, err := r.Client.Get(key).Result()
	return val, err
}

func (r *Redis) HGet(key string, field string) (string, error) {
	key = parsePrefix(key)
	val, err := r.Client.HGet( key, field).Result()
	return val, err
}

func (r *Redis) Set(key string, val interface{}, exp time.Duration) (string, error) {
	key = parsePrefix(key)
	d, err := r.Client.Set(key, val, exp).Result()
	return d, err
}

func (r *Redis) HSet(key string, field string, val interface{}) (bool, error) {
	key = parsePrefix(key)
	b, err := r.Client.HSet(key, field, val).Result()
	return b, err
}

func (r *Redis) HGetAll(key string) (map[string]string, error) {
	key = parsePrefix(key)
	m, err := r.Client.HGetAll(key).Result()
	return m, err
}

