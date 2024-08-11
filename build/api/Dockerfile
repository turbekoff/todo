FROM golang:1.20-alpine as builder

WORKDIR /build

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /api ./cmd/api/main.go


FROM alpine

RUN apk update --no-cache && \
    apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /api /api

EXPOSE 80

CMD ["/api"]
