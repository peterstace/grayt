FROM golang:1

WORKDIR /go/src/github.com/peterstace/grayt
COPY . .

RUN go install ./...
