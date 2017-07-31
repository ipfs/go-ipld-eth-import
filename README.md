# go-ipld-eth-import

`ETHEREUM CLIENT` ---> `IPFS`

Get your ethereum information to IPFS.

## Overview

(TODO)

## Usage

### Retriever of Block Data (by block hash)

(TODO)

### Loops

#### Loop of Last Block

(TODO)

#### Loop per Block Interval

(TODO)

## Other scripts

(TODO)

## Mustekala

(TODO)

### Dockerfile: IPFS + ETH Plugin

Quick and dirty workaround, so you can have hit the ground running with your
IPFS daemon.

#### Building it

Step up into the directory `mustekala/ipfs-plus-eth-plugin/`.

```
docker build -t ipfs-plus-eth-plugin -f Dockerfile --squash .
```

Don't forget the `--squash` option, if you don't want to end up with a 1.9GB
image!

#### Using it

Ideally, just use it inside the _mustekala_ combo. You can however, use it as
an `ipfs` command in its own right.

```
# Init your IPFS repo
docker run --rm -ti -v ~/.ipfs:/root/.ipfs ipfs-plus-eth-plugin ipfs init

# Start the IPFS daemon
docker run --rm -ti -v ~/.ipfs:/root/.ipfs ipfs-plus-eth-plugin ipfs daemon

# Insert some ETH data
cat eth-block-body-json-997522 | docker exec -i ipfs ipfs dag put --input-enc json --format eth-block

# (Will give you the cid z43AaGEzuAXhWf9pWAm63QCERtFpqcc6gQX3QBBNaG1syxGGhg6)

# Retrieve your ETH data
docker exec -i ipfs ipfs dag get z43AaGEzuAXhWf9pWAm63QCERtFpqcc6gQX3QBBNaG1syxGGhg6
```

## TODO

(TODO: Fill this list up)

## Credits

* Ideas taken from this @whyrsuleeping [repo](https://raw.githubusercontent.com/whyrusleeping/ipfs-eth-import/master/main.go).
