FROM golang:1.18 as build
WORKDIR /src
COPY go.mod go.sum /src/
RUN go mod download -x
COPY . .
RUN GOOS=linux GOARCH=amd64 make bin/registry

FROM scratch
LABEL org.opencontainers.image.authors="OBADA Foundation <techops@obada.io>"

ARG REGISTRY_UID="2017"
ARG REGISTRY_GID="2017"

COPY --from=build /src/registry /usr/sbin/registry
COPY --from=build "/etc/ssl/certs/ca-certificates.crt" "/etc/ssl/certs/ca-certificates.crt"

EXPOSE 2017
EXPOSE 2018

USER "${REGISTRY_UID}:${REGISTRY_GID}"
ENTRYPOINT ["/usr/sbin/registry", "server"]
