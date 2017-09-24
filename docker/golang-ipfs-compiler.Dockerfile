FROM ubuntu:16.04

MAINTAINER Herman Junge "chpdg42@gmail.com"

################################################################################
#
#  The purpose of this image is to prepare a compiling environment
#  for your go applications (cold and hot importers), as well as your
#  ipfs + plugin client.
#
#  Please refer to this directory README.md file for a step-by-step tutorial.
#
################################################################################

RUN apt-get update && \
    apt-get install -y \
        curl \
        git \
        make \
        vim-gnome \
        build-essential && \
    apt-get upgrade -y

WORKDIR /tmp

RUN curl -O https://storage.googleapis.com/golang/go1.8.3.linux-amd64.tar.gz && \
    tar -xvf go1.8.3.linux-amd64.tar.gz && \
    GOROOT=/tmp/go GOPATH=/tmp/gopath /tmp/go/bin/go get golang.org/x/tools/cmd/goimports && \
    mv /tmp/go /usr/local && \
    mv /tmp/gopath/bin/goimports /usr/local/bin/goimports

ENV GOROOT /usr/local/go
ENV GOPATH /go
ENV PATH "${PATH}:${GOROOT}/bin:${GOPATH}/bin"
