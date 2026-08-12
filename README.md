# Go Social Network

A study project built to explore **Hexagonal Architecture (Ports & Adapters)**
in Go, using a social-network domain with users, posts, comments, follows,
authentication, and subscriptions.

The core idea: the application core (`domain`, `usecase`, and `ports`) should
not depend on delivery or persistence details. GraphQL, gRPC, PostgreSQL,
Redis, JWT, and other tools are connected through adapters and infrastructure
wiring.

## Goals

- Practice Hexagonal Architecture / Ports & Adapters in Go
- Build a schema-first GraphQL API with gqlgen
- Add a gRPC input adapter alongside GraphQL, as a step toward exploring
  microservice-style splits
- Keep business logic separated from delivery and persistence mechanisms
- Use PostgreSQL through repository adapters
- Use Redis Pub/Sub for GraphQL subscriptions
- Add JWT authentication with protected GraphQL operations

## Tech stack

| Concern | Choice |
|---|---|
| Language | Go |
| API | GraphQL with gqlgen, gRPC with protoc |
| HTTP server | `net/http` |
| Database | PostgreSQL with GORM |
| Pub/Sub | Redis |
| Authentication | JWT access/refresh tokens + bcrypt password hashing |
| Local runtime | Docker Compose |

## Architecture

See [ARCHITECTURE.md](./ARCHITECTURE.md) for the full breakdown.

Quick summary:

```text
Outside world -> Adapter (in) -> Port -> Usecase -> Port -> Adapter (out) -> Outside world
```

Current important directories:

```text
cmd/server                         # Application entry point, HTTP + gRPC wiring
cmd/grpc-client                    # Demo gRPC client (calls CreateUser, CreatePost)
proto                              # .proto contracts shared by server and client
proto/gen                          # Generated gRPC/protobuf Go code (do not hand-edit)
internal/app/domain                # Domain entities and domain errors
internal/app/usecase                # User, auth, post, and comment use cases
internal/app/ports                 # Usecase and repository contracts
internal/adapters/in/graphql       # GraphQL schema, generated code, resolvers
internal/adapters/in/grpc          # gRPC service implementations (UserService, PostService)
internal/adapters/out/db/postgres  # PostgreSQL repositories and GORM models
internal/adapters/out/pubsub/redis # Redis Pub/Sub adapter
internal/infra                     # Config, DB, Redis, GraphQL/gRPC wiring, JWT helpers
```

## Getting started

### Prerequisites

- Go 1.25+
- Docker and Docker Compose

### Environment

Copy the example environment file:

```bash
cp .env.example .env
```

Important variables:

```env
PORT=8080
GRPC_PORT=50051
POSTGRES_HOST=db
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=password
POSTGRES_DB_NAME=social_network
REDIS_ADDR=redis:6379
JWT_SECRET_KEY=secretkey
```

For local execution outside Docker, set `POSTGRES_HOST=localhost` and
`REDIS_ADDR=localhost:6379`.

### Running with Docker Compose

```bash
docker compose up --build
```

The GraphQL Playground will be available at:

```text
http://localhost:8080/
```

The gRPC server listens on `localhost:50051` (mapped from the container, see
`GRPC_PORT`).

### Running locally

Start PostgreSQL and Redis first, then run:

```bash
go run ./cmd/server
```

This starts both the GraphQL/HTTP server (`PORT`) and the gRPC server
(`GRPC_PORT`) in the same process, each on its own listener.

## Authentication

Passwords are hashed with bcrypt before being stored. Login returns an
`accessToken`, a `refreshToken`, and the authenticated user.

Create a user:

```graphql
mutation {
  createUser(input: {
    name: "Ada Lovelace"
    email: "ada@example.com"
    password: "secret123"
  }) {
    user {
      id
      name
      email
    }
    errors {
      field
      message
    }
  }
}
```

Login:

```graphql
mutation {
  login(email: "ada@example.com", password: "secret123") {
    accessToken
    refreshToken
    user {
      id
      name
      email
    }
    errors {
      field
      message
    }
  }
}
```

For protected operations, add this in the Playground **HTTP Headers** panel:

```json
{
  "Authorization": "Bearer <accessToken>"
}
```

Refresh an access token:

```graphql
mutation {
  refreshToken(refreshToken: "<refreshToken>") {
    accessToken
    refreshToken
    user {
      id
      name
    }
  }
}
```

Protected mutations use the authenticated user from the token. For example,
`createPost` does not receive `authorId`:

```graphql
mutation {
  createPost(input: { description: "Hello from an authenticated user" }) {
    post {
      id
      description
      author {
        id
        name
      }
    }
    errors {
      field
      message
    }
  }
}
```

The GraphQL schema uses the `@auth` directive to mark operations that require a
valid access token.

## gRPC

Alongside GraphQL, the server exposes a gRPC API defined in
`proto/socialnetwork.proto` (`UserService.CreateUser`,
`PostService.CreatePost`). It currently plays the role of an internal,
trusted-caller RPC layer rather than a public-facing entry point like
GraphQL: `CreatePostRequest` takes `author_id` explicitly instead of deriving
it from an authenticated session, since the caller is assumed to already be a
trusted party (e.g. another service), not an end-user client. See
[ARCHITECTURE.md](./ARCHITECTURE.md) for the reasoning behind this
distinction.

A demo client that calls both RPCs in sequence (create a user, then create a
post authored by that user) lives in `cmd/grpc-client`:

```bash
go run ./cmd/grpc-client
```

It connects to `localhost:<GRPC_PORT>` (default `50051`), so it works both
against `go run ./cmd/server` and against the Docker Compose `app` service.

## Development

Regenerate GraphQL code after changing the schema:

```bash
go run github.com/99designs/gqlgen generate
```

Regenerate gRPC/protobuf code after changing `proto/socialnetwork.proto`
(requires `protoc` with the `protoc-gen-go` and `protoc-gen-go-grpc`
plugins):

```bash
protoc -I proto \
       --go_out=proto/gen --go_opt=paths=source_relative \
       --go-grpc_out=proto/gen --go-grpc_opt=paths=source_relative \
       socialnetwork.proto
```

Run tests:

```bash
go test ./...
```

In sandboxed environments, you may need a writable Go build cache:

```bash
GOCACHE=/private/tmp/go-social-network-gocache go test ./...
```

## License

Not defined yet.
