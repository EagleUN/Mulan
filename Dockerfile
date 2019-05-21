FROM golang:alpine AS build-env
WORKDIR /go/src/github.com/jschavesr/mulan
COPY . .
RUN apk add --no-cache git

RUN go get -d ./...
RUN go build -o goapp
CMD [ "./goapp" ]