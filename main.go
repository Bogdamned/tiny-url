package main

import (
	"log"

	"github.com/spf13/viper"
)

var id uint64

func main() {
	db, err := NewRedisDB(Config{
		Addr:     viper.GetString("db.addr"),
		Password: "",
		DB:       viper.GetInt("db.db"),
	})
	if err != nil {
		log.Fatalf("failed to initialize db: %s", err.Error())
	}

	urlDb := NewUrlRedis(db)
	urlDb.initCounter()

	handlers := NewHandler(NewLocalCache(), urlDb)
	run(handlers)
}

func initConfig() error {
	viper.AddConfigPath("config")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
