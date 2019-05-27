FROM golang
WORKDIR /go/src/github.com/jschavesr/mulan
COPY . .

RUN go get -d ./...
RUN go build -o goapp
CMD [ "./goapp" ]
