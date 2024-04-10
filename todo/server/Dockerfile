FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o myserver ./server/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/myserver .

COPY ./static /app/static

CMD ["/app/myserver"]