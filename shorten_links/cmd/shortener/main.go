package main

import (
	"fmt"
	"net/http"
	my_errors "shorten_links/internal/errors"
	"shorten_links/internal/transport"
)

func main() {
	fmt.Println(my_errors.ErrEqualJSON.Error())
	transport.HandleCreate()

	fmt.Println(http.ListenAndServe(":80", nil))
}
