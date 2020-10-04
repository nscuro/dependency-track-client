FROM golang:1.15-alpine as build
RUN apk --no-cache add build-base git
WORKDIR /opt/dtrack
COPY . .
RUN make install

FROM alpine:latest
COPY --from=build /go/bin/dtrack /bin/dtrack
ENTRYPOINT [ "/bin/dtrack" ]
