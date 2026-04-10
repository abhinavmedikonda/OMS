package main

import (
	"context"
	"log"
	"time"

	"github.com/abhinavmedikonda/OMS/catalog"
	"github.com/abhinavmedikonda/OMS/observability"
	"github.com/kelseyhightower/envconfig"
	"github.com/tinrab/retry"
)

type Config struct {
	DatabaseURL string `envconfig:"DATABASE_URL"`
}

func main() {
	ctx := context.Background()
	provider, err := observability.Setup(ctx, "catalog")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := provider.Shutdown(ctx); err != nil {
			log.Printf("observability shutdown: %v", err)
		}
	}()

	var cfg Config
	err = envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	var r catalog.Repository
	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		r, err = catalog.NewElasticRepository(cfg.DatabaseURL)
		if err != nil {
			log.Println(err)
		}
		return
	})
	defer r.Close()

	log.Println("Listening on port 8080...")
	s := catalog.NewService(r)
	log.Fatal(catalog.ListenGRPC(s, 8080))
}
