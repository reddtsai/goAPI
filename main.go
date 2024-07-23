package main

import (
	"context"
	"log"

	"github.com/reddtsai/goAPI/cmd/api"
)

func main() {
	ctx := context.Background()
	err := api.Execute(ctx)
	if err != nil {
		log.Fatalln(err)
	}

}
