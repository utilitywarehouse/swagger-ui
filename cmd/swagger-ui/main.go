package main

import (
	"log"
	"net/http"

	"github.com/utilitywarehouse/swagger-ui/swaggerui"
)

func main() {
	s := swaggerui.Handler()

	log.Panic(http.ListenAndServe(":8080", s))
}
