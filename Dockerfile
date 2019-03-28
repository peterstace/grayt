FROM golang:1

# TODO: remove uuid once it's no longer a dependency
RUN go get github.com/satori/go.uuid

WORKDIR /go/src/github.com/peterstace/grayt
COPY . .

RUN go install ./...
