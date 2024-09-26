FROM golang:1.22.5-alpine AS builder

WORKDIR /app

COPY . .

RUN go build -o application .


FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/application .

COPY --from=builder /app/.env .

EXPOSE 8080

CMD ["./application"]
