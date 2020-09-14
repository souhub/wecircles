FROM golang:1.15.0-alpine3.12

WORKDIR /go/src/github.com/souhub/wecircles

RUN apk add --no-cache \
    alpine-sdk \
    git

COPY . .

RUN go build

ENV DB_USER root

ENV DB_PROTOCOL tcp

ENV DB_NAME wecircles

EXPOSE 80

CMD ["go","run","/go/src/github.com/souhub/wecircles"]
