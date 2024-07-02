package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/go-resty/resty/v2"
)

type Req struct {
	LL string `json:"long_link"`
}

func main() {
	url := "http://google.com"
	c := make(chan string, 1000000)
	go func() {
		for {
			url_t := url + fmt.Sprint(rand.Int63n(9223372036854775807))
			r := Req{LL: url_t}
			b, _ := json.Marshal(r)
			resty.New().R().SetBody(b).Post("http://localhost/shortlink/get")
			c <- url_t
			time.Sleep(50 * time.Millisecond)
		}
	}()
	// go func() {
	// 	url_t := <-c
	// 	resty.New().R().Get(fmt.Sprintf("http://localhost/shortlink/redirect?url=%s", url_t))
	// 	c <- url_t
	// 	time.Sleep(5 * time.Millisecond)
	// }()
	time.Sleep(10 * time.Minute)
}
