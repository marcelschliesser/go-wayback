FROM golang:1.20.6-alpine3.18

WORKDIR /app

COPY src/ .

RUN go build -o main .

EXPOSE 8080

CMD ["./main"]
