# birdwatcher

birdwatcher is a small HTTP server meant to provide an API defined by
Barry O'Donovan's
[birds-eye](https://github.com/inex/birds-eye-design/) to
[the BIRD routing daemon](http://bird.network.cz/).

## Why

The [INEX implementation](https://github.com/inex/birdseye) of
birdseye runs PHP, which is not always desirable (and performant)
in a routeserver setting. By using Go, we are able to work with
regular binaries, which means deployment and maintenance might be
more convenient.

Our version also has a few more capabilities, as you will
discover when looking at [the modules section](https://github.com/ecix/birdwatcher/blob/master/etc/ecix/birdwatcher.conf)
of the config.

## Installation

You will need to have go installed to build the package.
Running `go get github.com/ecix/birdwatcher` will give you
a binary. You might need to cross-compile it for your
bird-running servive (`GOARCH` and `GOOS` are your friends).

We provide a Makefile for more advanced compilation/configuration.
Running `make linux` will create a Linux executable (by default for
`amd64`, but that is configurable by providing the `ARCH` argument
to the Makefile).

### Building an RPM

Building RPMs is supported through [fpm](https://github.com/jordansissel/fpm).
If you have `fpm` installed locally, you can run `make rpm`
to create a RPM in the folder `RPMS`. If you have a remote
build server with `fpm` installed, you can build and fetch
an RPM with `make remote_rpm BUILD_SERVER=<buildserver_url>`
(requires SSH access).

### Deployment

If you want to deploy `birdwatcher` on a system that uses
RPMs, you should be able to install it after following the
instructions on [building an RPM](#building-an-rpm).

We do not currently support other deployment methods.

## Configuration

An example config with sane defaults is provided in
[etc/ecix/birdwatcher.conf](https://github.com/ecix/birdwatcher/blob/master/etc/ecix/birdwatcher.conf).
You should be able to use it out of the box. If you need
to change it, it is well-commented and hopefully intuitive.
If you do not know how to configure it, please consider opening
[an issue](https://github.com/ecix/birdwatcher/issues/new).

## How

In the background `birdwatcher` runs the `birdc` client, sends
commands and parses the result. It also leverages simple caching
techniques to help reduce the load on the bird service.

## Who

Initially developed by Daniel and MC from [Netnod](https://www.netnod.se/) in
two days at the RIPE 73 IXP Tools Hackathon in Madrid, Spain.

Running bird and parsing the results was added by [Veit Heller](https://github.com/hellerve/) on behalf of [ecix](http://ecix.net/).
