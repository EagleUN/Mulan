FROM golang
WORKDIR /go/src/github.com/jschavesr/mulan
COPY . .

RUN go get -d ./...
RUN go build -o goapp
EXPOSE 8080
CMD [ "./goapp" ]
