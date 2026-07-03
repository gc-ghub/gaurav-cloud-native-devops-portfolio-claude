FROM golang:1.25-alpine AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/backend .

FROM alpine:3.20
RUN apk add --no-cache ca-certificates && \
    adduser -D -H appuser
WORKDIR /app
COPY --from=builder /out/backend ./backend
USER appuser
EXPOSE 8080
ENTRYPOINT ["./backend"]
