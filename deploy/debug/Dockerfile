FROM golang:1.22-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/store-service ./cmd/store-service

FROM alpine:3.19

RUN adduser -D -H -s /sbin/nologin app
WORKDIR /home/app

COPY --from=build /bin/store-service /usr/local/bin/store-service

EXPOSE 8080
USER app

ENTRYPOINT ["/usr/local/bin/store-service"]

