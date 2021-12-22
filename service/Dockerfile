FROM golang:alpine as builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o add-service -ldflags="-w -s" service/main.go

FROM scratch

WORKDIR /app

COPY --from=builder /app/add-service /usr/bin/

ENTRYPOINT [ "add-service" ]
