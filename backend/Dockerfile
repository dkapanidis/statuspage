FROM golang:1.8.0-alpine
COPY server.go src/
RUN go build src/server.go
EXPOSE 80
CMD ["./server"]
