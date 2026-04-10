package main

import (
	"context"
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/abhinavmedikonda/OMS/observability"
	"github.com/kelseyhightower/envconfig"
)

type AppConfig struct {
	AccountURL string `envconfig:"ACCOUNT_SERVICE_URL"`
	CatalogURL string `envconfig:"CATALOG_SERVICE_URL"`
	OrderURL   string `envconfig:"ORDER_SERVICE_URL"`
}

func main() {
	ctx := context.Background()
	provider, err := observability.Setup(ctx, "graphql")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := provider.Shutdown(ctx); err != nil {
			log.Printf("observability shutdown: %v", err)
		}
	}()

	var cfg AppConfig
	err = envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	s, err := NewGraphQLServer(cfg.AccountURL, cfg.CatalogURL, cfg.OrderURL)
	if err != nil {
		log.Fatal(err)
	}

	srv := handler.NewDefaultServer(s.ToExecutableSchema())
	http.Handle("/graphql", observability.HTTPHandler(srv, "graphql"))
	http.Handle("/playground", observability.HTTPHandler(playground.Handler("abhi", "/graphql"), "playground"))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
