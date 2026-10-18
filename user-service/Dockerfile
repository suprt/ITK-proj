FROM golang:1.27.1-alpine AS builder

WORKDIR /app

COPY proto/go.mod proto/go.sum ./proto/
COPY shared/go.mod shared/go.sum ./shared/
COPY user-service/go.mod user-service/go.sum ./user-service/

ENV GOWORK=off
RUN cd user-service &&  go mod download

COPY proto/ ./proto/
COPY shared/ ./shared/
COPY user-service/ ./user-service/

RUN sh -c 'cd user-service && CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/user-service ./cmd'

FROM alpine:latest

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /out/user-service /app/user-service
COPY user-service/config.yaml /app/config.yaml

EXPOSE 50051

ENTRYPOINT ["/app/user-service"]