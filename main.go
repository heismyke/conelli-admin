package main

import (
	"fmt"
	"log"

	"github.com/conelli/admin-backend/api"
	"github.com/conelli/admin-backend/config"
)

func main() {
	if err := config.Envs.Validate(); err != nil {
		log.Fatalf("configuration error: %v", err)
	}

	server, err := api.NewApi(fmt.Sprintf(":%s", config.Envs.PORT))
	if err != nil {
		log.Fatalf("failed to initialize application: %v", err)
	}

	log.Fatal(server.Run())
}
