FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY backend/go.mod backend/go.sum ./backend/
WORKDIR /app/backend
RUN go mod download

WORKDIR /app
COPY . .
WORKDIR /app/backend
RUN go build -o /out/site-grupo-mano .

FROM alpine:3.21

WORKDIR /app/backend
RUN apk add --no-cache ca-certificates

COPY --from=builder /out/site-grupo-mano /app/backend/site-grupo-mano
COPY --from=builder /app/frontend /app/frontend

EXPOSE 8080

ENV PORT=8080 \
    GIN_MODE=release \
    FRONTEND_DIR=/app/frontend

CMD ["/app/backend/site-grupo-mano"]