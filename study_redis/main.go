package main

import (
	//"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-redis/redis"
	"github.com/rs/zerolog"
	// "github.com/redis/go-redis/v9"
	// "gopkg.in/redis.v5"
)

type Redis struct {
	rbd *redis.Client
}

func NewRedis() *Redis {
	return &Redis{
		rbd: redis.NewClient(&redis.Options{
			Addr:     "localhost:7000",
			Password: "",
			DB:       0,
		}),
	}
}

func (r *Redis) CreateClient() {
	r.rbd = redis.NewClient(&redis.Options{
		Addr:     "localhost:7000",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

func (r *Redis) Set(key, value string, ttl time.Duration) error {
	return r.rbd.Set(key, value, ttl).Err()
}

func (r *Redis) Get(key string) (string, error) {
	val, err := r.rbd.Get(key).Result()
	return val, err
}

func ExampleClient() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:7000",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	err := rdb.Set("key", "value", 0).Err()
	if err != nil {

		panic(err)
	}

	val, err := rdb.Get("key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)

	val2, err := rdb.Get("key2").Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", val2)
	}
}

type HandlersBuilder struct {
	r  *Redis
	lg zerolog.Logger
}

func HandleCreate() {
	hb := HandlersBuilder{
		r:  NewRedis(),
		lg: zerolog.New(os.Stderr).With().Timestamp().Logger().Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.UnixDate}),
	}
	if err := hb.r.rbd.Ping().Err(); err != nil {
		log.Fatal(err)
	}
	hb.r.CreateClient()
	http.HandleFunc("/get", hb.Get())
	http.HandleFunc("/set", hb.Set())
}

func ParseSet(r *http.Request) (string, string, time.Duration, error) {
	var data struct {
		Key      string `json:"key"`
		Value    string `json:"value"`
		TimeLife string `json:"time"`
	}
	err := json.NewDecoder(r.Body).Decode(&data)
	var empty time.Duration
	if err != nil {
		return data.Key, data.Value, empty, err
	}
	d, _ := time.ParseDuration(data.TimeLife)
	return data.Key, data.Value, d, err
}

func ParseGet(r *http.Request) (string, error) {
	var data struct {
		Key string `json:"key"`
	}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		return data.Key, err
	}
	return data.Key, err
}

func (hb *HandlersBuilder) Set() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			key, val, ttl, err_pars := ParseSet(r)
			if err_pars != nil {
				fmt.Fprint(w, "Неверный запрос")
				w.WriteHeader(http.StatusBadRequest)
			} else {
				err := hb.r.rbd.Set(key, val, ttl)
				if err.Val() != "OK" {
					// 	hb.lg.Warn().
					// 		Msgf("message from func Set %v", err.Err().Error())
					w.WriteHeader(http.StatusBadRequest)
				} else {
					fmt.Fprint(w, "Элемент с ключом: ", key, " успешно добавлен")
				}
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}

func (hb *HandlersBuilder) Get() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			key, err_pars := ParseGet(r)
			if err_pars != nil {
				fmt.Fprint(w, "Неверный запрос")
				w.WriteHeader(http.StatusBadRequest)
			} else {
				val, err := hb.r.Get(key)
				if err != nil {
					w.WriteHeader(http.StatusNotFound)
					fmt.Fprint(w, "Элемент с ключом ", key, " не найден")
				} else {
					fmt.Fprint(w, "Элемент с ключом ", key, " найден: ", val)
				}
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}

func main() {
	HandleCreate()

	fmt.Println(http.ListenAndServe("8081:80", nil))
}
