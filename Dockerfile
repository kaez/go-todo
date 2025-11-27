FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache gcc musl-dev sqlite-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o go-todo ./cmd/api

FROM alpine:latest

RUN apk --no-cache add ca-certificates sqlite-libs

WORKDIR /root/

COPY --from=builder /app/go-todo .

ENV PORT=8002
ENV DB_PATH=/data/todos.db

EXPOSE 8082

RUN mkdir -p /data

CMD ["./go-todo"]
