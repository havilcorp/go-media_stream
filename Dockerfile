FROM golang:1.20 AS builder
WORKDIR /build
COPY . .
RUN go test -cover ./...
RUN go build -o main cmd/app/main.go

FROM golang:1.20
RUN apt-get update && apt-get install ffmpeg -y
WORKDIR /build
COPY static static
COPY templates templates
COPY --from=builder /build/main /build/main
EXPOSE 8080
CMD ["./main"]