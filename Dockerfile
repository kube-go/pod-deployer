ARG GOVERSION=1.16.4
FROM golang:${GOVERSION} as build

ARG GOARCH
ENV GOARCH=${GOARCH}

WORKDIR /go/src/pod-deployer
COPY . /go/src/pod-deployer

RUN make build-local

# Copy into base image
FROM gcr.io/distroless/static:latest
COPY --from=build /go/bin/pod-deployer /

USER pod-deployer

CMD ["/pod-deployer"]

EXPOSE 8080
