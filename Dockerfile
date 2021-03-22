FROM golang:latest

ENV GOPATH /usr/local/go/bin/go
RUN apt-get update
COPY . /tmp/app
WORKDIR /tmp/app

EXPOSE 12345
CMD ["/usr/local/go/bin/go", "run", "main.go", "election.go", "node.go"]
