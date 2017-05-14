FROM golang:1.8.0-alpine AS build
COPY server.go src/
RUN go build src/server.go

FROM alpine:3.5
COPY --from=build /go/server /server
COPY static static
EXPOSE 9000
CMD ["/server"]
