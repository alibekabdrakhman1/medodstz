FROM golang:1.21

RUN go version

WORKDIR /app

COPY .. .
EXPOSE 8080

RUN go mod download
RUN go build ./cmd/app

CMD ["go", "run", "./cmd/app"]