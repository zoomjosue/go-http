FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod ./
COPY . .
RUN go build -o servidor .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/servidor .
COPY --from=builder /app/data ./data
EXPOSE 24918
CMD ["./servidor"]