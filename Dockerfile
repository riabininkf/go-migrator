FROM golang:1.14

WORKDIR /app
COPY . .

RUN go build -o ./.build/gomigrator main.go

RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.27.0
