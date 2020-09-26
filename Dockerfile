FROM golang:1.15-alpine

RUN apk add --no-cache git

WORKDIR /go/src/github.com/zherdev/go-panic-beautifier
COPY . .

ENV GOPATH /go
ENV PATH $HOME/goApps/bin:$PATH

RUN go get -d -v ./...
RUN go build -v -o go-panic-beautifier cmd/goPanicBeautifier/main.go

EXPOSE 80

CMD ["./go-panic-beautifier", "conf.yaml"]
