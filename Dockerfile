# This is a multi-stage Dockerfile and requires >= Docker 17.05
# https://docs.docker.com/engine/userguide/eng-image/multistage-build/
FROM golang:1.13-alpine as builder
RUN apk add build-base


RUN mkdir -p $GOPATH/src/github.com/shaybix/weather-monster
WORKDIR $GOPATH/src/github.com/shaybix/weather-monster

ADD . .
RUN go install github.com/shaybix/weather-monster

## Service container
FROM alpine
RUN apk add --no-cache bash
RUN apk add --no-cache ca-certificates

WORKDIR /bin/

COPY --from=builder /go/bin/weather-monster .

EXPOSE 3000

CMD exec /bin/weather-monster