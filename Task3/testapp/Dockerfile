FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY main.go .
RUN go build -o testapp main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/testapp .
EXPOSE 8080
CMD ["./testapp"]
