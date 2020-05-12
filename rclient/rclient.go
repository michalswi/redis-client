package rclient

import "github.com/go-redis/redis"

func NewClient(rhost string, rport string) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: rhost + ":" + rport,
		// password not set
		Password: "",
		// use default db
		DB: 0,
	})
	return client
}
