FROM golang:1.14-alpine as build
RUN apk --no-cache add git
WORKDIR /opt/dtrack
COPY . .
RUN GO111MODULE=on go build -v ./cmd/dtrack

FROM alpine:latest
COPY --from=build /opt/dtrack/dtrack /bin/dtrack
ENTRYPOINT [ "/bin/dtrack" ]
