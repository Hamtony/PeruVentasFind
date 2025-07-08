# Dockerfile
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copiar archivos de módulos primero
COPY go.mod go.sum ./
RUN go mod download

# Luego el resto
COPY . .

RUN go build -o main_node main_node.go
RUN go build -o worker_node worker_node.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main_node .
COPY --from=builder /app/worker_node .
COPY ./data ./data

RUN apk add --no-cache mongodb-tools

CMD ["./main_node"]
