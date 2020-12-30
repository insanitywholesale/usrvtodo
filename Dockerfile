# build stage
FROM golang:alpine as build
RUN apk add --no-cache git build-base alpine-sdk
ENV GOPROXY http://oldboi.hell,direct
ENV GO111MODULE off
WORKDIR /go/src/usrvtodo
COPY *go* ./
COPY *tmpl ./
COPY *db ./
RUN go get -d -v
# this doesn't actually work fully
# workaround in run stage is for this
# after adding the tags, one of the previous errors went away
RUN go build -v -tags 'osusergo netgo static static_build' -ldflags '-linkmode external -extldflags "-static"'
RUN go install -v


# run stage
# changed from busybox to alpine
# busybox was giving me:
# standard_init_linux.go:211: exec user process caused "no such file or directory"
FROM alpine
# unable to link fully statically so this is needed
#COPY --from=build /lib/x86*/libdl* /lib/
# set path for database
ENV DB_PATH /data/todo.db
# put gin in production mode
ENV GIN_MODE release
# add and use non-root user
RUN adduser -D -u 5000 -g 5000 app
USER app:app
WORKDIR /data
WORKDIR /go/bin
# copy files from build stage that are required at runtime
COPY --from=build /go/bin/usrvtodo /go/bin/
COPY --from=build /go/src/usrvtodo/index.tmpl /go/bin/
COPY --from=build /go/src/usrvtodo/error.tmpl /go/bin/
COPY --from=build /go/src/usrvtodo/todo.db /data/
EXPOSE 8080
VOLUME /data
CMD ["/go/bin/usrvtodo"]
