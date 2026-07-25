package main

import (
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/vektah/gqlparser/v2/ast"

	"github.com/gianpaoloaranha/go-social-network/internal/adapters/in/graphql/generated"
	"github.com/gianpaoloaranha/go-social-network/internal/infra/config"
	"github.com/gianpaoloaranha/go-social-network/internal/infra/db"
	"github.com/gianpaoloaranha/go-social-network/internal/infra/graphql"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	postgresDB, closePostgres, err := db.ConnectToPostgres(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer closePostgres()

	if err := db.RunPostgresMigrations(postgresDB); err != nil {
		log.Fatal(err)
	}

	resolver := graphql.BuildResolvers(postgresDB)

	srv := handler.New(generated.NewExecutableSchema(generated.Config{Resolvers: resolver}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, nil))
}
