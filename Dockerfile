FROM golang:alpine AS builder
RUN mkdir /biscotti
ADD . /biscotti/
WORKDIR /biscotti
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/biscotti ./cmd/biscotti/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /biscotti/bin/biscotti .
EXPOSE 3000
CMD ["./biscotti"]
