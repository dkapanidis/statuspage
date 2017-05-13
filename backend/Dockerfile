FROM golang:1.8.0-alpine
COPY server.go src/
RUN go build src/server.go
COPY static static
EXPOSE 9000
CMD ["./server"]
