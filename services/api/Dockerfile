FROM golang:1.23 AS build

ENV GOPATH=/
WORKDIR /src/
COPY . .

RUN go mod download; CGO_ENABLED=0 go build -o /order-system-api ./services/api/main.go

FROM alpine:3.17

COPY --from=build /order-system-api /order-system-api
COPY ./services/api/migrations/ /migrations

EXPOSE 8080
CMD ["./order-system-api"]