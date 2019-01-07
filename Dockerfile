FROM golang:alpine as BUILD

RUN apk add --no-cache make git

ENV CGO_ENABLED=0

WORKDIR /opt/app

COPY . .

RUN make build

FROM alpine as RUN

RUN apk add ca-certificates

COPY --from=BUILD /opt/app/bin/* /usr/bin/

ENTRYPOINT ["weather-reporter"]
