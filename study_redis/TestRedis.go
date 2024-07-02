package main

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	rbd *redis.Client
}

var ctx = context.Background()

func (r *Redis) CreateClient() {
	r.rbd = redis.NewClient(&redis.Options{
		Addr:     "localhost:7000",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

func (r *Redis) Set(key, value string) {
	err := r.rbd.Set(ctx, key, value, 0).Err()
	if err != nil {
		panic(err)
	}
}

func (r *Redis) Get(key string) {
	val, err := r.rbd.Get(ctx, key).Result()
	if err == redis.Nil {
		fmt.Println(key, " does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println(key, val)
	}
}

func ExampleClient() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:7000",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	err := rdb.Set(ctx, "key", "value", 0).Err()
	if err != nil {

		panic(err)
	}

	val, err := rdb.Get(ctx, "key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)

	val2, err := rdb.Get(ctx, "key2").Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", val2)
	}
	// Output: key value
	// key2 does not exist
}

func main() {
	//ExampleClient()
	var ex Redis
	ex.CreateClient()
	ex.Set("apple", "яблоко")
	ex.Set("banana", "банан")
	ex.Set("grapes", "виноград")
	ex.Get("kiwi")
	ex.Get("apple")

}
