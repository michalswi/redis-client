package apis

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
)

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var kid = 1

func SetRedis(logger *log.Logger, client *redis.Client) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		logger.Printf("Processing redis SET request in %s\n", time.Now().Sub(startTime))

		var user User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			logger.Printf("Error: %v", err)
			w.WriteHeader(http.StatusBadRequest)
		}

		json, err := json.Marshal(User{Name: user.Name, Age: user.Age})
		if err != nil {
			logger.Printf("Error: %v", err)
			w.WriteHeader(http.StatusBadRequest)
		}

		kids := fmt.Sprintf("id%d", kid)
		err = client.Set(kids, json, 0).Err()
		if err != nil {
			logger.Printf("Error SET: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			logger.Printf("New entry added with id (key): %s", kids)
			kid++
		}
	}
}

func GetRedis(logger *log.Logger, client *redis.Client) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		logger.Printf("Processing redis GET request in %s\n", time.Now().Sub(startTime))
		uid := mux.Vars(r)
		kids := fmt.Sprintf("id%s", uid["uid"])
		val, err := client.Get(kids).Result()
		if err != nil {
			logger.Printf("Error GET: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Write([]byte(val))
	}
}

func PingRedis(logger *log.Logger, client *redis.Client) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		logger.Printf("Processing request in %s\n", time.Now().Sub(startTime))
		pong, err := client.Ping().Result()
		if err != nil {
			logger.Printf("Error: %v", err)
			w.WriteHeader(http.StatusServiceUnavailable)
		}
		w.Write([]byte(pong))
	}
}

func Home(logger *log.Logger, version string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		logger.Printf("Processing request in %s\n", time.Now().Sub(startTime))
		message := "redis-client"
		hostname, err := os.Hostname()
		if err != nil {
			logger.Printf("Get 'host name' failed: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		var rawHTML = `
		<html>
		<h1>%s</h1>
		<p><b>Hostname</b>: %s; <b>Version</b>: %s</p>
		</html>
		`
		fmt.Fprintf(w, rawHTML, message, hostname, version)
	}
}
