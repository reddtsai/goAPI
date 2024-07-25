package main

import (
	"context"
	"log"

	"github.com/spf13/viper"

	"github.com/reddtsai/goAPI/cmd/http"
	"github.com/reddtsai/goAPI/pkg/blockaction/api"
	"github.com/reddtsai/goAPI/pkg/blockaction/storage"
)

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("conf.d")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("init config : %v\n", err)
	}
}

func main() {
	ctx := context.Background()

	db, err := storage.NewBlockActionDB(ctx, storage.BlockActionDBCfg{
		UserName: viper.GetString("mysql-options.username"),
		Password: viper.GetString("mysql-options.password"),
		Address:  viper.GetString("mysql-options.addr"),
		DBName:   viper.GetString("mysql-options.db"),
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
