FROM golang:1.25.0-trixie AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o go-social-network ./cmd/server/main.go

FROM golang:trixie

WORKDIR /app

COPY --from=builder /app/go-social-network .
COPY --from=builder /app/internal/adapters/in/graphql/schema ./internal/adapters/in/graphql/schema

EXPOSE 8080

CMD ["./go-social-network"]
