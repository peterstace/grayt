FROM golang:1

WORKDIR /go/src/github.com/peterstace/grayt
COPY . .

# TODO: remove uuid once it's no longer a dependency
RUN go get github.com/satori/go.uuid
RUN go install -v ./...
