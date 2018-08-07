FROM golang:alpine
WORKDIR $GOPATH/src/github.com/naighes/imposter
RUN apk add --update bash git zip openssh
COPY . .
RUN /bin/bash scripts/build.sh --release
RUN cd $GOPATH/src/github.com/naighes/imposter/pkg/$(go env GOOS)_$(go env GOARCH) && cp ./imposter $GOPATH/bin/imposter

FROM golang:alpine
LABEL maintainer="Nicola Baldi (@naighes) <nic.baldi@gmail.com>"
LABEL "com.naighes.imposter.version"="${VERSION}"
RUN apk add --update openssh
COPY --from=0 /go/bin/imposter /go/bin/imposter
EXPOSE 8080
WORKDIR $GOPATH
ENTRYPOINT ["imposter"]
