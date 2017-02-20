Release Status
==============

A small daemon to track the state of QA releases *because they take
long...*

The service tracks the state of releases, which so far is binary:

  - running
  - not running

You can change the state by sending a get request to either the start or
the stop endpoint and the server will respond with true or false:

| state/request | start | stop |
|---------------|:-----:|:----:|
| not running   |   1   |   0  |
| running       |   0   |   1  |

Basically it allows you to only start *one release at a time*.


Usage
-----

To get the server running locally you should build a binary which you
can run or distribute:

    go build
	RS_PORT=8080 ./releasestatus

The program will refuse to run unless a port is specified in the
environment.

It is also easily possible to run it in a Docker container from scratch.
([more info](#docker))

To mark a release as running call the start endpoint:

	curl localhost:8080/start

To stop mark it as done, call the stop endpoint:

	curl localhost:8080/stop

You can also add a `name` parameter to let other attempted starters know
who started the running release:

	curl localhost:8080/start?name=sion

The response will be `0` or `1`, based on the table above, and you can
use this in deploy scripts to determine whether to proceed with the
release or abort.


Hacking
-------

The program's structure is quite simple.

It has the callbacks for the start and stop endpoints in the main
function which sets up the web server and a helper function to get the
server port from the environment at start up.  The state of the running
release is stored in a `Release` struct holding information about who
started it and when.

If you'd like to contribute, please make sure you run `gofmt` on your
code first, to keep the style consistent.


Docker
------

There is a [docker/Dockerfile](docker/Dockerfile) for building
releasestatus as a Docker container from scratch.  It adds the
Linux-built releasestatus binary as the Docker command, as well as a
zoneinfo archive for setting the timezone inside the container.

It has a Makefile to reduce the amount of typing for the Docker build
because the docker commands are quite long.

The Makefile's default target will **build and run the Docker image**.
It has some default variables which can be overridden from the
environment:

 - `TZ`:   the timezone the container will run in (default: CET)
 - `PORT`: the port the service will run on (default: 8080)
 - `TAG`:  the image tag in Docker (default: releasestatus)

It also has separate make targets:

 - `container`: just build the image
 - `run`: run the container (assumes it is built)
 - `clean`: removes the binaries and archive used during the process

The container from scratch is very small, currently ~5mb, it contains
only the go binary and zone info.  For this to work correctly the binary
is built with netgo included (which has a native implementation of host
lookup etc), so that it can handle network requests.  The go build tool
uses special environment variables to build for Linux amd64, suitable
for running in Docker.  Because I also want to be able to set the time
zone for the container (I'd like the logging timestamps to be in Central
European Time) it needs a zoneinfo archive, which is borrowed from the
system's `/usr/share/zoneinfo` directory at build time.  These are the
only 2 files the container currently needs to run, and in fact, if
you're fine with running it in UCT, you don't even need the zoneinfo
file.

The run target will use your environment variables (if provided) to
override my default settings, and the container target can use the `TAG`
variable if you'd like to give the image a special tag for deployment.
The container will be removed when it stops running, to avoid clutter
among your containers.

This all-together makes building, running and deploying the dockerised
version of this service really convenient and easy.

---
Written by Siôn le Roux <sion.leroux@schibsted.hu>
