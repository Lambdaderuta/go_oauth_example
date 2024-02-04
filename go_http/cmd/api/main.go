package main

import (
	"fmt"
	"go_http/internal/server"
	"go_http/internal/auth"
)

func main() {

	auth.NewAuth()
	server := server.NewServer()

	err := server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
