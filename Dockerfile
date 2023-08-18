FROM golang:1.20

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o ./generate_coding_challenge_server_go

EXPOSE 8080

COPY configuration configuration
ENV APP_ENVIRONMENT production

CMD ["./generate_coding_challenge_server"]