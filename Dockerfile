FROM golang:1.24

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o chat-server

EXPOSE 8888

CMD ["./chat-server"]
