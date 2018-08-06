FROM golang:alpine
LABEL maintainer="Nicola Baldi (@${OWNER}) <nic.baldi@gmail.com>"

ARG PRODUCT_NAME=UNSPECIFIED
ARG VERSION=UNSPECIFIED
ARG OWNER=UNSPECIFIED

RUN apk add --update git bash openssh zip

EXPOSE 8080

WORKDIR $GOPATH/src/github.com/${OWNER}/${PRODUCT_NAME}
COPY . .
RUN /bin/bash scripts/build.sh --release
RUN cp "./pkg/$(go env GOOS)_$(go env GOARCH)/${PRODUCT_NAME}" $GOPATH/bin/imposter

LABEL "com.${OWNER}.${PRODUCT_NAME}.version"="${VERSION}"

WORKDIR $GOPATH
ENTRYPOINT ["imposter"]
