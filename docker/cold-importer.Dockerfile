FROM ubuntu:16.04

MAINTAINER Herman Junge "chpdg42@gmail.com"

ADD ./build/bin/docker-cold-importer /usr/bin/cold-importer

ENTRYPOINT ["/usr/bin/cold-importer"]
