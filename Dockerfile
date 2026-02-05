FROM golang:1.24-alpine AS builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o kypidbot ./cmd/bot

FROM alpine:3.21
WORKDIR /app
COPY --from=builder /build/kypidbot .
COPY --from=builder /build/messages.yaml .
CMD ["./kypidbot"]
