FROM --platform=linux/amd64 golang:alpine AS builder

ARG CGO_ENABLED=0
WORKDIR /app

ENV TZ=Europe/Moscow

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN GOARCH=amd64 CGO_ENABLED=0 GOOS=linux go build -o backend ./cmd

FROM alpine

RUN apk add --no-cache tzdata

WORKDIR /app

COPY --from=builder /app/backend /backend
RUN chmod +x /backend

COPY --from=builder /app/database/migrations /app/database/migrations

CMD ["/backend"]