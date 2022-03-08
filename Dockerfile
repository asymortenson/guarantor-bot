FROM golang:1.15-alpine3.12 AS builder

COPY . /guarantorplace.com/
WORKDIR /guarantorplace.com/

RUN go mod download
RUN go build -o ./bin/bot ./cmd/bot/

FROM alpine:latest

WORKDIR /root/

COPY --from=0 /guarantorplace.com/bin/bot .
COPY --from=0 /guarantorplace.com/configs configs/

EXPOSE 8080

CMD ["./bot"]