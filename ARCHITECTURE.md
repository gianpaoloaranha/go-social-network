# Project Architecture - Go Social Network

This project follows **Hexagonal Architecture (Ports & Adapters)**. The
application core lives under `internal/app` and should remain independent from
GraphQL, PostgreSQL, Redis, HTTP, JWT libraries, and framework details.

The implemented input adapter is GraphQL. The implemented output adapters are
PostgreSQL repositories and Redis Pub/Sub.

## General Principle

```text
Outside world -> Adapter (in) -> Port -> Usecase -> Port -> Adapter (out) -> Outside world
```

The dependency direction points inward:

- adapters depend on ports/usecases
- usecases depend on domain and ports
- domain depends only on itself and the standard library where possible

`cmd/server/main.go` is the composition root. It loads configuration, opens
infrastructure resources, creates adapters/usecases/resolvers, and connects the
HTTP server.

## Current Structure

```text
.
├── cmd/server
│   └── main.go
├── internal
│   ├── adapters
│   │   ├── in/graphql
│   │   │   ├── generated
│   │   │   ├── resolver
│   │   │   └── schema
│   │   └── out
│   │       ├── db/postgres
│   │       │   └── model
│   │       └── pubsub/redis
│   ├── app
│   │   ├── domain
│   │   ├── ports
│   │   │   ├── auth
│   │   │   ├── comment
│   │   │   ├── post
│   │   │   ├── pubsub
│   │   │   └── user
│   │   └── usecase
│   └── infra
│       ├── authentication
│       ├── config
│       ├── db
│       ├── graphql
│       └── redis
├── gqlgen.yml
├── docker-compose.yml
└── README.md
```

## Layer Responsibilities

| Layer | Responsibility |
|---|---|
| `internal/app/domain` | Domain entities, domain errors, password hashing helpers, auth session payload |
| `internal/app/usecase` | Application flows: users, auth, posts, comments |
| `internal/app/ports` | Interfaces exposed or required by the application core |
| `internal/adapters/in/graphql` | GraphQL schema, generated gqlgen code, resolvers, GraphQL mapping/payload helpers |
| `internal/adapters/out/db/postgres` | GORM models and repository implementations |
| `internal/adapters/out/pubsub/redis` | Redis Pub/Sub implementation for subscriptions |
| `internal/infra` | Configuration, resource creation, JWT helpers, and application wiring |
| `cmd/server` | Runtime composition and HTTP server setup |

## Application Core

The core is centered around four usecase groups:

- `user`: create/update/delete users, follow/unfollow, query users
- `auth`: login, refresh token, logout
- `post`: create/update/delete posts, query posts, publish post events
- `comment`: create/update/delete comments, query comments, publish comment events

The core talks to persistence and Pub/Sub only through ports. For example,
`post` usecases receive repository and publisher interfaces rather than a GORM
database or Redis client.

## GraphQL Adapter

GraphQL is schema-first:

- SDL lives in `internal/adapters/in/graphql/schema/schema.graphql`
- gqlgen generated files live in `internal/adapters/in/graphql/generated`
- resolver implementations live in `internal/adapters/in/graphql/resolver`

The schema includes:

- public user/auth mutations such as `createUser`, `login`, and `refreshToken`
- protected mutations marked with `@auth`
- public queries for users, posts, and comments
- subscriptions for post/comment events

After changing the schema, regenerate code with:

```bash
go run github.com/99designs/gqlgen generate
```

## Authentication Flow

Authentication uses:

- bcrypt for password hashing
- JWT access tokens for protected GraphQL operations
- JWT refresh tokens for issuing a new access/refresh token pair
- a GraphQL `@auth` directive for declaring protected operations
- an HTTP middleware that reads `Authorization: Bearer <token>`

The flow is:

1. `createUser` receives a plaintext password and stores only its bcrypt hash.
2. `login` validates email/password and returns `AuthPayload`.
3. `AuthPayload` contains `accessToken`, `refreshToken`, and `user`.
4. The HTTP middleware validates access tokens and stores the authenticated
   `userID` in `context.Context`.
5. The GraphQL `@auth` directive rejects protected operations without a valid
   authenticated user in context.
6. Protected resolvers use the `userID` from context instead of trusting actor
   IDs from GraphQL input.

Access and refresh tokens have distinct claims through `tokenType`:

```json
{
  "userId": "user-id",
  "tokenType": "access"
}
```

```json
{
  "userId": "user-id",
  "tokenType": "refresh"
}
```

## Authenticated Actor Rule

For operations performed by the current user, the actor comes from the access
token, not from client input.

Examples:

- `createPost(input: { description })` uses the token user as `authorId`.
- `createComment(input: { postId, message })` uses the token user as
  `authorId`.
- `followUser(input: { userToFollowId })` uses the token user as follower.
- `updateUser(input: { ... })` updates the token user.
- `deleteUser` deletes the token user.

This avoids accepting arbitrary `userId`/`authorId` values from clients for
actions that should belong to the authenticated user.

## Persistence And Pub/Sub

PostgreSQL is accessed through GORM models and repositories in
`internal/adapters/out/db/postgres`. The database connection and migrations are
created in `internal/infra/db`.

Redis is used as a Pub/Sub backend for GraphQL subscriptions. The Redis client
is created in `internal/infra/redis`, while the adapter that implements the
application Pub/Sub port lives in `internal/adapters/out/pubsub/redis`.

## Infrastructure vs Adapters

Infrastructure packages create or configure resources:

- `infra/config` loads environment variables
- `infra/db` opens PostgreSQL and runs migrations
- `infra/redis` opens the Redis client
- `infra/graphql` builds resolvers and injects usecases
- `infra/authentication` creates/validates tokens and provides auth context
  helpers

Adapters implement application-facing contracts:

- PostgreSQL repositories implement repository ports
- Redis Pub/Sub implements the Pub/Sub port
- GraphQL resolvers call usecase ports

## Golden Rule

No file inside `internal/app` should import from `internal/adapters`.

The current `auth` usecase imports `internal/infra/authentication` for JWT
creation/validation. That is acceptable for this learning stage, but a stricter
hexagonal version would introduce an auth/token port, implemented by an
infrastructure adapter, so the usecase depends only on an interface.

## Roadmap

- [x] GraphQL input adapter
- [x] PostgreSQL repository adapter
- [x] Redis Pub/Sub adapter for subscriptions
- [x] JWT login, refresh token, and GraphQL `@auth`
- [ ] Persist refresh sessions/tokens for real logout and token revocation
- [ ] Authorization ownership checks for update/delete post/comment
- [ ] Dataloaders to reduce N+1 queries in nested GraphQL fields
- [ ] HTTP/REST adapter
- [ ] gRPC adapter
- [ ] CLI adapter
