FROM golang:alpine

WORKDIR /go/src/app

COPY . .

RUN go get -d -v ./...

RUN go mod tidy

RUN go install -v ./...

EXPOSE 8080

CMD ["app"]
