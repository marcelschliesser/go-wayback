FROM golang:1.21.0-alpine AS builder
WORKDIR /app
COPY ./src/ ./
RUN go mod download
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o webserver .

FROM alpine
WORKDIR /app
COPY --from=builder /app/ /app/
CMD ["./webserver"]