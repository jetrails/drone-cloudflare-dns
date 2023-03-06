FROM golang:1.17.6-alpine as build-env
WORKDIR /usr/src/app
COPY . .
RUN go mod download && go mod verify
RUN go build -o ./bin/main ./cmd/main/main.go

FROM alpine
LABEL maintainer development@jetrails.com
COPY --from=build-env /usr/src/app/bin/main /usr/local/bin/main
RUN chmod +x /usr/local/bin/main
ENTRYPOINT [ "/usr/local/bin/main" ]
