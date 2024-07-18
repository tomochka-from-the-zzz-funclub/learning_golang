package database

import (
	"fmt"
	"log"
	my_errors "shorten_links/internal/errors"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis"
)

type Redis struct {
	rbd *redis.Client
}

func NewRedis() *Redis {
	client := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	_, err := client.Ping().Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
		return nil
	}

	return &Redis{
		rbd: client,
	}
}

func (r *Redis) Set(shortlink string, data InfoLLink) error {
	err_l := r.rbd.Set(fmt.Sprintf("%s_longlink", shortlink), data.GetLongLink(), data.GetDeath().Sub(time.Now())).Err()
	if err_l != nil {
		return err_l
	}
	err_r := r.rbd.Set(fmt.Sprintf("%s_redirect", shortlink), strconv.Itoa(data.GetStatRedirect()), data.GetDeath().Sub(time.Now())).Err()
	if err_r != nil {
		return err_r
	}
	err_d := r.rbd.Set(fmt.Sprintf("%s_death", shortlink), data.GetDeath().Format(time.RFC3339), data.GetDeath().Sub(time.Now())).Err()
	if err_d != nil {
		return err_d
	}
	return nil
}

func (r *Redis) GetLongL(shortlink string) (string, error) {
	longlink, err := r.rbd.Get(fmt.Sprintf("%s_longlink", shortlink)).Result()
	if err != nil {
		return "", my_errors.ErrNoLlink
	}
	return longlink, nil
}

func (r *Redis) GetRedirect(shortlink string) (int, error) {
	redirect_str, err := r.rbd.Get(fmt.Sprintf("%s_redirect", shortlink)).Result()
	if err != nil {
		return 0, my_errors.ErrNoRedirect
	}
	redirect, err := strconv.Atoi(redirect_str)
	if err != nil {
		return 0, my_errors.ErrNoRedirect
	}
	return redirect, nil
}

func (r *Redis) GetDataDeath(shortlink string) (time.Time, error) {
	data_str, err := r.rbd.Get(fmt.Sprintf("%s_death", shortlink)).Result()
	if err != nil {
		return time.Time{}, my_errors.ErrNoDataDeath
	}
	datadeath, err := time.Parse(time.RFC3339, data_str)
	if err != nil {
		return datadeath, my_errors.ErrNoDataDeath
	}
	return datadeath, nil
}

func (r *Redis) GetAllData(shortlink string) (InfoLLink, error) {
	var (
		data DataLLink
		err  error
	)
	data.LongLink, err = r.rbd.Get(fmt.Sprintf("%s_longlink", shortlink)).Result()
	if err != nil {
		return data, my_errors.ErrNoLlink
	}
	redirect_str, err := r.rbd.Get(fmt.Sprintf("%s_redirect", shortlink)).Result()
	if err != nil {
		return data, my_errors.ErrNoRedirect
	}
	data.StatRedirect, err = strconv.Atoi(redirect_str)
	if err != nil {
		return data, my_errors.ErrNoRedirect
	}
	data_str, err := r.rbd.Get(fmt.Sprintf("%s_death", shortlink)).Result()
	if err != nil {
		return data, my_errors.ErrNoDataDeath
	}
	data.Death, err = time.Parse(time.RFC3339, data_str)
	if err != nil {
		return data, my_errors.ErrNoDataDeath
	}
	return data, nil
}

func (r *Redis) GetAll() (map[string]InfoLLink, error) {
	// Получаем все ключи, хранящиеся в Redis
	keys, err := r.rbd.Keys("*").Result()
	if err != nil {
		return nil, err
	}
	// Извлекаем значения по всем ключам
	values, err := r.rbd.MGet(keys...).Result()
	if err != nil {
		return nil, err
	}

	result := make(map[string]InfoLLink)
	for i, key := range keys {
		var data InfoLLink
		data.SetLongLink("")
		data.SetStatRedirect(0)
		data.SetDeath(time.Date(2020, time.April, 17, 12, 34, 56, 0, time.UTC))

		if strings.HasSuffix(key, "_longlink") {
			key = strings.TrimSuffix(key, "_longlink")
			longlink := values[i].(string)
			data.SetDeath(result[key].GetDeath())
			data.SetLongLink(longlink)
			data.SetStatRedirect(result[key].GetStatRedirect())
			result[key] = data
		} else if strings.HasSuffix(key, "_redirect") {
			key = strings.TrimSuffix(key, "_redirect")
			redirect, err := strconv.Atoi(values[i].(string))
			if err != nil {
				return result, err
			}
			data.SetDeath(result[key].GetDeath())
			data.SetLongLink(result[key].GetLongLink())
			data.SetStatRedirect(redirect)
			result[key] = data
		} else if strings.HasSuffix(key, "_death") {
			key = strings.TrimSuffix(key, "_death")
			datadeath, err := time.Parse(time.RFC3339, values[i].(string))
			if err != nil {
				return result, err
			}
			data.SetDeath(datadeath)
			data.SetLongLink(result[key].GetLongLink())
			data.SetStatRedirect(result[key].GetStatRedirect())
			result[key] = data
		}
	}

	return result, nil
}

func (r *Redis) Increment(key string) error {
	return r.rbd.IncrBy(key+"_redirect", 1).Err()
}
