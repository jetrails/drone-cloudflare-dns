FROM golang:alpine as build-env
WORKDIR /go/src/app
ENV GOOS=linux GOARCH=amd64 CGO_ENABLED=0
RUN apk -Uuq add git dep ca-certificates
COPY Gopkg.toml Gopkg.lock src/main.go /go/src/app/
RUN dep ensure
RUN go build -o main

FROM alpine
LABEL maintainer development@jetrails.com
COPY --from=build-env /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build-env /go/src/app/main /usr/local/bin/main
RUN chmod +x /usr/local/bin/main
ENTRYPOINT [ "main" ]
