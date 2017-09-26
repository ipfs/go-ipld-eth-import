## Docker

### Overview

Guide to build your programs (including IPFS + eth plugin) in docker containers.

### Requirements

* Docker. This tutorial was prepared using `Docker version 17.06.2-ce, build cec0b72`.

### Compiler machine

In order to have light images, we are preparing an image containing all elements
needed to build both cold and hot importers, and IPFS with its Ethereum plugin.
Your mileage may vary, so, you can really look this guide as a help to get you
started and optimize your images as you see fit.

#### Step 1: Build the image

To create the compiler image, fire the following command

```
docker build -t golang-ipfs-compiler -f docker/golang-ipfs-compiler.Dockerfile .
```

Sit down, start netflix and wait for everything to be nice and installed.

You will find your image with

```
[ hj: go-ipld-eth-import ]$ docker images
REPOSITORY             TAG                 IMAGE ID            CREATED             SIZE
golang-ipfs-compiler   latest              21f21206dba1        8 seconds ago       948MB
```

#### Step 2: Obtain all needed dependencies

This is a more lenghty process, we will get you the current dependencies of both
IPFS and Ethereum in a directory of your choose, that way, we don't spoil actual
versions of your programs in the local machines, should be the case you are not
running Ubuntu in there (Hello OSX Users).

Replace `$HOME/.golang-linux-ipfs-compiler` for the directory you want, run this
command, and keep enjoying that `Rick and Morty` chapter.

```
docker run \
	-ti --rm \
	--name golang-ipfs-compiler \
	-v $HOME/.golang-linux-ipfs-compiler:/go \
	-v $PWD/docker/compiler-get-all-deps.sh:/workspace/compiler-get-all-deps.sh \
	-w /workspace \
	golang-ipfs-compiler /workspace/compiler-get-all-deps.sh
```

When this ends, We should be good to go and ready to make those images.

### Build the cold importer image

Just run this command

```
./docker/build-cold-importer-image.sh
```

Or

```
make docker-cold
```

It will

* Build (with the image above) the file ./build/bin/docker-cold-importer
* Create a docker image holding that executable.

From there you can just push the image to your docker repository,
making it available to your servers.

**NOTE**: We are aware that the image weights +/- `130MB`. Optimization
is advisable, Pull Requests are welcome!

### Run the cold importer container

* TODO
* Need to move forward in the connection to the ipfs shell.
