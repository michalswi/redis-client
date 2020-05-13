package rclient

import (
	"crypto/tls"
	"strconv"

	"github.com/go-redis/redis"
)

func NewClient(rhost string, rport string, rpass string, rtls string) *redis.Client {
	// TLSConfig is required to communicate with Azure Cache for Redis
	if rtls != "" {
		rtlsb, _ := strconv.ParseBool(rtls)
		client := redis.NewClient(&redis.Options{
			Addr:      rhost + ":" + rport,
			Password:  rpass,
			DB:        0,
			TLSConfig: &tls.Config{InsecureSkipVerify: rtlsb},
		})
		return client
	}
	client := redis.NewClient(&redis.Options{
		Addr:     rhost + ":" + rport,
		Password: rpass,
		DB:       0,
	})
	return client
}
