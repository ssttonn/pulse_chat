# Build Stage
FROM golang:1.22.0-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum into root folder
COPY go.mod go.sum* ./
RUN go mod download

# Copy all source code into root folder
COPY . .

# Build binary file with compressed mode 
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main ./src/edge-ws/cmd/api

# Run Stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 8081

CMD [ "./main" ]