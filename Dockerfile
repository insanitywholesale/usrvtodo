# build stage
FROM golang:1.14 as build
WORKDIR /go/src/usrvtodo
COPY . .
ENV CGO_ENABLED 0
RUN go get -d -v
RUN go install -v


# run stage
FROM busybox:glibc
RUN adduser -D -u 5000 app
USER app:app
WORKDIR /go/bin/
COPY --from=build /go/bin/usrvtodo /go/bin/usrvtodo
EXPOSE 8080
ENV GIN_MODE release
CMD ["/go/bin/usrvtodo"]
