ARG SRC_PATH_GLOBAL=/goradius

#
# BUILDER
#
FROM golang:1.20.5-alpine as builder

ARG SRC_PATH_GLOBAL
ENV SRC_PATH=${SRC_PATH_GLOBAL}

WORKDIR ${SRC_PATH}
COPY ./src ${SRC_PATH}/src

WORKDIR ${SRC_PATH}/src
RUN CGO_ENABLED=0 go build -v -o goradius ./


FROM alpine:latest as base

ARG SRC_PATH_GLOBAL
ENV SRC_PATH=${SRC_PATH_GLOBAL}

EXPOSE 2083

COPY certs ${SRC_PATH}/certs

#
# PRODUCTION IMAGE
#
FROM base as prod

ARG SRC_PATH_GLOBAL
ARG SRC_PATH=${SRC_PATH_GLOBAL}

COPY --from=builder ${SRC_PATH}/src/goradius ${SRC_PATH}
WORKDIR ${SRC_PATH}
ENTRYPOINT [ "./goradius" ]

FROM base as dev

COPY --from=golang:1.20.5-alpine /usr/local/go/ /usr/local/go/
ENV PATH="/usr/local/go/bin:${PATH}"

ARG SRC_PATH_GLOBAL
ARG SRC_PATH=${SRC_PATH_GLOBAL}
COPY . ${SRC_PATH}

WORKDIR ${SRC_PATH}/src
RUN go get -d ./...
ENTRYPOINT [ "go", "run", "main.go" ]