FROM golang:1.24
WORKDIR /app

RUN mkdir -p insider-project
COPY . /app/insider-project/

WORKDIR /app/insider-project
RUN go mod tidy
RUN go build -o app ./cmd/server

RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz && \
    mv migrate /usr/local/bin/migrate

EXPOSE 8080
CMD ["./app"]
