package main

import (
	"context"
	"log"
	"net/http"
	"time"

	gqlgen "github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/vektah/gqlparser/v2/ast"

	"github.com/gianpaoloaranha/go-social-network/internal/adapters/in/graphql/generated"
	redispubsub "github.com/gianpaoloaranha/go-social-network/internal/adapters/out/pubsub/redis"
	"github.com/gianpaoloaranha/go-social-network/internal/app/domain"
	"github.com/gianpaoloaranha/go-social-network/internal/infra/authentication"
	"github.com/gianpaoloaranha/go-social-network/internal/infra/config"
	"github.com/gianpaoloaranha/go-social-network/internal/infra/db"
	appgraphql "github.com/gianpaoloaranha/go-social-network/internal/infra/graphql"
	redisinfra "github.com/gianpaoloaranha/go-social-network/internal/infra/redis"
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

	redisClient, closeRedis, err := redisinfra.Connect(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer closeRedis()

	broker := redispubsub.NewBroker(redisClient)
	resolver := appgraphql.BuildResolvers(postgresDB, broker, broker)

	graphqlConfig := generated.Config{Resolvers: resolver}
	graphqlConfig.Directives.Auth = func(ctx context.Context, obj any, next gqlgen.Resolver) (any, error) {
		if _, ok := authentication.UserIDFromContext(ctx); !ok {
			return nil, domain.NewError(domain.ErrUnauthorized, "Unauthorized")
		}

		return next(ctx)
	}

	srv := handler.New(generated.NewExecutableSchema(graphqlConfig))

	srv.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
	})
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", authentication.AuthMiddleware(srv))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, nil))
}
