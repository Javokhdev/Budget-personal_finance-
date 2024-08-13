FROM golang:1.22.5 AS builder
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o myapp ./server/server.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/myapp .

COPY .env .
EXPOSE 8070
CMD ["./myapp"]