ARG GOVERSION=1.16
FROM golang:${GOVERSION} as build-env

ARG GOARCH
ENV GOARCH=${GOARCH}

WORKDIR /go/src/pod-deployer
COPY . /go/src/pod-deployer

RUN make build-local

# Copy into base image
FROM alpine:3.7
COPY --from=build-env /go/bin/pod-deployer /
COPY --from=build-env /go/src/pod-deployer/views /views

RUN adduser -s /bin/false -h /home/deployer -D deployer

USER deployer

CMD ["/pod-deployer"]

EXPOSE 8080
