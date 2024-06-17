package main

import (
	"fmt"
	"hw/internal/transport"
	"net/http"
)

func main() {
	transport.HandleCreate()

	fmt.Println(http.ListenAndServe(":80", nil))
}
