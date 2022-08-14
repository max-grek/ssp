package main

import (
	"fmt"
	"log"
	"test-assignment-cookie-sync/config"
	"test-assignment-cookie-sync/connector"
	"test-assignment-cookie-sync/service"
)

func main() {
	cfg, err := config.ParseFlags()
	if err != nil {
		die(err)
	}
	die(newApp(cfg))
}

func die(err error) {
	if err == nil {
		return
	}
	log.Fatalln(err)
}

func newApp(cfg *config.Config) error {
	db, err := connector.Initialize(&cfg.StorageConfig)
	if err != nil {
		return fmt.Errorf("initialize connector: %v", err)
	}
	http := newHTTPService(&cfg.HTTPConfig)
	ss := service.NewCookie(db)
	http.registerRoutes(ss)
	return http.run()
}
