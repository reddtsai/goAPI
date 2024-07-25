package main

import (
	"context"
	"log"

	"github.com/reddtsai/goAPI/cmd/http"
	"github.com/reddtsai/goAPI/pkg/blockaction/api"
	"github.com/reddtsai/goAPI/pkg/blockaction/storage"
)

func main() {
	ctx := context.Background()

	// TODO : config
	db, err := storage.NewBlockActionDB(ctx, storage.BlockActionDBCfg{
		UserName: "root",
		Password: "!QAZ2wsx",
		Address:  "127.0.0.1:3306",
		DBName:   "blockaction",
	})
	if err != nil {
		log.Fatalf("init db : %v\n", err)
	}
	httpHandler, err := api.NewBlockActionApi(api.SetStorage(db))
	if err != nil {
		log.Fatalf("init api : %v\n", err)
	}
	err = http.Execute(ctx, httpHandler)
	if err != nil {
		log.Fatalf("execute cmd : %v\n", err)
	}

}
