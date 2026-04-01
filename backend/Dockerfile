FROM --platform=linux/amd64 golang:alpine AS builder

ARG CGO_ENABLED=0
WORKDIR /app

ENV TZ=Europe/Moscow

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN GOARCH=amd64 CGO_ENABLED=0 GOOS=linux go build -o backend ./cmd

FROM alpine

RUN apk add --no-cache tzdata ffmpeg python3 py3-pip \
    && pip3 install --break-system-packages yt-dlp

WORKDIR /app

COPY --from=builder /app/backend /backend
RUN chmod +x /backend

COPY --from=builder /app/database/migrations /app/database/migrations

CMD ["/backend"]