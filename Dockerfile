FROM golang:latest 
#as builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .

#RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .
#FROM alpine:latest
#RUN apk --no-cache add ca-certificates
#WORKDIR /root/
#COPY --from=builder /app/main .

EXPOSE 12345
ENTRYPOINT ["go", "run", "node.go", "consensus.go", "election.go", "main.go"]
