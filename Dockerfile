FROM golang:1

WORKDIR /go/src/github.com/peterstace/grayt
COPY . .

RUN go get github.com/satori/go.uuid
RUN go install -v ./...
