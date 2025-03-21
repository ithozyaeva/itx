FROM --platform=linux/amd64 golang:alpine AS builder

ARG CGO_ENABLED=0
WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN GOARCH=amd64 CGO_ENABLED=0 GOOS=linux go build -o server ./cmd

FROM alpine

WORKDIR /app

COPY --from=builder /app/server /server
RUN chmod +x /server

CMD ["/server"]