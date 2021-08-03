FROM golang:1.16-alpine3.14 as bulder

COPY . /pocketerClient
WORKDIR /pocketerClient

RUN go mod download
RUN go build -o ./runtime/.bin/bot cmd/bot/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=0 /pocketerClient/runtime/.bin/bot .
COPY --from=0 /pocketerClient/configs configs/

EXPOSE 80

CMD ["./bot"]