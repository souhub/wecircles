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

ENV WECIRCLES_S3_IMAGE_BUCKET wecircles-img-dev

ENV IMAGE_PATH https://wecircles-img-dev.s3-ap-northeast-1.amazonaws.com/web

EXPOSE 80

CMD ["go","run","/go/src/github.com/souhub/wecircles"]
