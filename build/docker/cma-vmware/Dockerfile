# Initial preparations
FROM golang:1.11-alpine
RUN apk update && apk add git make curl

# Let's add the code in
ADD . /go/src/gitlab.com/mvenezia/cma-vmware

# Build package in container
WORKDIR /go/src/gitlab.com/mvenezia/cma-vmware
RUN make -f build/Makefile container-build

# Now that we have build a static binary, we can stuff it in a bare-bones alpine image
FROM alpine
COPY --from=0 /go/src/gitlab.com/mvenezia/cma-vmware/cma-vmware /cma-vmware

ENTRYPOINT ["/cma-vmware"]
