package main

import (
	"fmt"
	"log"
	"meetup/internal/auth"
	"meetup/internal/pkg/config"
	"meetup/internal/pkg/repository/postgres"
	"meetup/internal/router"
)

func main() {
	// config
	cfg := config.GetConf()
	fmt.Println("user:", cfg.DBUsername)
	// databases
	postgresDB := postgres.New(cfg.DBUsername, cfg.DBPassword, cfg.Port, cfg.DBName)
	// authenticator
	authenticator := auth.New(postgresDB)

	// router
	r := router.New(authenticator, postgresDB)
	fmt.Println(cfg.Port)
	log.Fatalln(r.Init(fmt.Sprintf(":%s", "8080")))

}
