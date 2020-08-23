FROM golang:1-alpine AS build

RUN apk update && apk add make git gcc musl-dev

ARG SERVICE

ADD . /tmp/${SERVICE}

WORKDIR /tmp/${SERVICE}

RUN make clean install
RUN make ${SERVICE}

RUN mv ${SERVICE} /${SERVICE}

FROM alpine:3.9

ARG SERVICE

ENV APP=${SERVICE}

RUN apk add --no-cache ca-certificates && mkdir /app
COPY --from=build /${SERVICE} /app/${SERVICE}

ENTRYPOINT /app/${APP}
