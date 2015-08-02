FROM golang:1.4.2

### Install godep, dependency manager for golang ###
RUN go get github.com/tools/godep
RUN go get -u github.com/golang/lint/golint

### Set up working directory ###
RUN mkdir -p /go/src/github.com/brainly/pencil-go
WORKDIR /go/src/github.com/brainly/pencil-go
ADD . /go/src/github.com/brainly/pencil-go

### Build the project including dependencies ###
RUN godep go install

### Make sure all tests pass ###
RUN godep go test ./...
RUN golint ./...

ENTRYPOINT ["/go/bin/pencil-go"]
