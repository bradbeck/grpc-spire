FROM golang:alpine as builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o add-api -ldflags="-w -s" api/main.go

FROM scratch

WORKDIR /app

COPY --from=builder /app/add-api /usr/bin/

ENTRYPOINT [ "add-api" ]
