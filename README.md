# birdwatcher

birdwatcher is a small HTTP server meant to provide an API defined by
Barry O'Donovan's
[birds-eye](https://github.com/inex/birds-eye-design/) to
[the BIRD routing daemon](http://bird.network.cz/).

## Installation

You will need to have go installed to build the package.
Running `go get github.com/ecix/birdwatcher` will give you
a binary. You might need to cross-compile it for your
bird-running servive (`GOARCH` and `GOOS` are your friends).

## Why

The [INEX implementation](https://github.com/inex/birdseye) of
birdseye runs PHP, which is not always desirable (and performant)
in a routeserver setting. By using Go, we are able to work with
regular binaries, which means deployment and maintenance might be
more convenient.

## How

In the background `birdwatcher` runs the `birdc` client, sends
commands and parses the result. It also leverages simple caching
techniques to help reduce the load on the bird service.

## Who

Initially developed by Daniel and MC from [Netnod](https://www.netnod.se/) in
two days at the RIPE 73 IXP Tools Hackathon in Madrid, Spain.

Running bird and parsing the results was added by [Veit Heller](https://github.com/hellerve/) on behalf of [ecix](http://ecix.net/).
