FROM alpine:latest

RUN apk update && \
    apk upgrade && \
    apk --no-cache add ca-certificates && \
    rm -rf /var/cache/apk/*
WORKDIR /code
USER 1001
COPY bin/linux/namespace-controller .
ENTRYPOINT ["/code/namespace-controller"]
