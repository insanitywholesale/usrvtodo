# build stage
FROM golang:1.15 as build
ENV GOPROXY http://oldboi.hell,direct
ENV GO111MODULE off
WORKDIR /go/src/usrvtodo
COPY *go* ./
COPY *tmpl ./
COPY *db ./
RUN go get -d -v
# this doesn't actually work fully
# workaround in run stage is for this
RUN go build -v -ldflags '-linkmode external -extldflags "-static"'
RUN go install -v


# run stage
FROM busybox:glibc
# unable to link fully statically so this is needed
COPY --from=build /lib/x86*/libdl* /lib/
# set path for database
ENV DB_PATH /data/todo.db
# put gin in production mode
ENV GIN_MODE release
# add and use non-root user
RUN adduser -D -u 5000 -g 5000 app
USER app:app
WORKDIR /go/bin
# copy files from build stage that are required at runtime
COPY --from=build /go/bin/usrvtodo /go/bin/
COPY --from=build /go/src/usrvtodo/index.tmpl /go/bin/
COPY --from=build /go/src/usrvtodo/error.tmpl /go/bin/
COPY --from=build /go/src/usrvtodo/todo.db /data/
EXPOSE 8080
VOLUME /data
CMD ["/go/bin/usrvtodo"]
