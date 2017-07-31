FROM golang:1.8.3

# This is a quick and dirty un-optimized Dockerfile
# to make available ipfs with the ipld-eth plugin.
# You may want to squash the resulting image, as
# it weights around 1.91GB!
#
# docker build -t ipfs-plus-eth-plugin -f Dockerfile --squash .

# Get dependencies
RUN go get -u -v github.com/whyrusleeping/gx \
 && go get -u -v github.com/whyrusleeping/gx-go \
 && go get -v github.com/ethereum/go-ethereum \
 && go get -v github.com/ipfs/go-ipfs

# Get ipfs gx dependencies
RUN cd /go/src/github.com/ipfs/go-ipfs \
 && gx --verbose install

# Get the plugin itself
# Clone it from github and change its branch
RUN git clone https://github.com/hermanjunge/go-ipld-eth.git /go/src/github.com/ipfs/go-ipld-eth \
 && cd /go/src/github.com/ipfs/go-ipld-eth \
 && git checkout refactor/block-header

# Make changes to the plugin and build ipfs with the eth plugin
RUN echo "ipldeth github.com/ipfs/go-ipld-eth/plugin 0" >> /go/src/github.com/ipfs/go-ipfs/plugin/loader/preload_list \
 && sed -i 's/package main/package plugin/' /go/src/github.com/ipfs/go-ipld-eth/plugin/eth.go \
 && cd /go/src/github.com/ipfs/go-ipld-eth \
 && gx-go hook pre-test \
 && cd /go/src/github.com/ipfs/go-ipfs/ \
 && make build \
 && mv /go/src/github.com/ipfs/go-ipfs/cmd/ipfs/ipfs /usr/bin/ipfs

# Remove usused elements
ENV PATH /usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin

RUN rm -rf /usr/local/go \
 && rm -rf /go \
 && apt-get purge -y g++ gcc libc6-dev make pkg-config \
 && apt-get update \
 && apt-get upgrade -y \
 && apt-get clean -y \
 && apt-get autoclean -y \
 && apt-get autoremove -y \
 && rm -rf /usr/share/locale/* \
 && rm -rf /var/cache/debconf/*-old \
 && rm -rf /var/lib/apt/lists/* \
 && rm -rf /usr/share/doc/*
