FROM golang:1.23 AS build

ENV GOPATH=/
WORKDIR /src/
COPY . .

RUN go mod download; CGO_ENABLED=0 go build -o /order-system-bill ./services/billing/main.go

FROM alpine:3.17

COPY --from=build /order-system-bill /order-system-bill

CMD ["./order-system-bill"]
