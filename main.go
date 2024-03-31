package main

import (
	"checkrr-notification-service/Helpers"
	"checkrr-notification-service/Router"
	"checkrr-notification-service/Services"
	"github.com/spf13/viper"
	"log"
)

func main() {

	viper.SetConfigName("")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()

	if err != nil {
		Helpers.Log(err, "Reading Env")
	}

	ch := make(chan []byte, 1e5)

	go Services.StartConsumer(ch)

	app := Router.New(ch)

	log.Fatal(app.Listen(":8081"))

}
