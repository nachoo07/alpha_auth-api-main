FROM golang:1.23.1

WORKDIR /usr/src/app

ENV SSL_MODE=require

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /usr/local/bin/app ./cmd/api/main.go

EXPOSE 8080

COPY entrypoint.sh /usr/local/bin/entrypoint.sh
RUN chmod +x /usr/local/bin/entrypoint.sh

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
CMD ["/usr/local/bin/app"]
