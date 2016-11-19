FROM golang:1.7.3

ENV GOPATH /go
ENV USER root

RUN go get github.com/bndr/gojenkins \
	&& go get github.com/docker/docker/api/types \
	&& go get github.com/docker/docker/api/types/container \
	&& go get github.com/docker/docker/client \
	&& go get golang.org/x/net/context

COPY . /go/src/github.com/stevebargelt/Dockhand
RUN go get -d -v github.com/stevebargelt/Dockhand
RUN go install github.com/stevebargelt/Dockhand

WORKDIR /go/src/github.com/stevebargelt/Dockhand
RUN go get -d -v
RUN go build -o dockhand main.go
#RUN go test github.com/stevebargelt/Dockhand/...
CMD ["/go/src/github.com/stevebargelt/Dockhand/dockhand", "run"]