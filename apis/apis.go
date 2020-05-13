package apis

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
)

type person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type handlers struct {
	logger *log.Logger
	client *redis.Client
}

func getCounter(logger *log.Logger, client *redis.Client) string {
	val, err := client.Get("counter").Result()
	if err != nil {
		logger.Printf("Error getCounter: %v", err)
	}
	return val
}

func setCounter(logger *log.Logger, client *redis.Client, counter string) string {
	counteri, _ := strconv.Atoi(counter)
	counteri++
	err := client.Set("counter", counteri, 0).Err()
	if err != nil {
		logger.Printf("Error setCounter: %v", err)
	}
	return strconv.Itoa(counteri)
}

func SetRedis(logger *log.Logger, client *redis.Client) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		logger.Printf("Processing redis SET request in %s\n", time.Now().Sub(startTime))

		h := &handlers{
			logger: logger,
			client: client,
		}

		counter := getCounter(h.logger, h.client)
		if counter == "" {
			counter = "0"
			counter = setCounter(h.logger, h.client, counter)
		}

		var user person
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			logger.Printf("Error SET decode: %v", err)
			w.WriteHeader(http.StatusBadRequest)
		}

		json, err := json.Marshal(person{Name: user.Name, Age: user.Age})
		if err != nil {
			logger.Printf("Error SET marshal: %v", err)
			w.WriteHeader(http.StatusBadRequest)
		}

		kids := fmt.Sprintf("id%s", counter)
		err = client.Set(kids, json, 0).Err()
		if err != nil {
			logger.Printf("Error SET: %v", err)
			w.WriteHeader(http.StatusBadRequest)
		} else {
			logger.Printf("New entry added with id (key): %s", kids)
			// increase counter
			setCounter(h.logger, h.client, counter)
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
			w.WriteHeader(http.StatusBadRequest)
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
