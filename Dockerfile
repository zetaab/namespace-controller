FROM alpine:latest

RUN apk add --no-cache ca-certificates
WORKDIR /code
USER 1001
COPY bin/linux/namespace-controller .
ENTRYPOINT ["/code/namespace-controller"]
