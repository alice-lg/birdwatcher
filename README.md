# birdwatcher

birdwatcher is a small HTTP server meant to provide an API defined by
Barry O'Donovan's
[birds-eye](https://github.com/inex/birds-eye-design/) to
[the BIRD internet routing daemon](http://bird.network.cz/).

## Why

The [INEX implementation](https://github.com/inex/birdseye) of
birdseye runs PHP, which is not always desirable (and performant)
in a route server setting. By using Go, we are able to work with
regular binaries, which means deployment and maintenance might be
more convenient.

Our version also has a few more capabilities, as you will
discover when looking at [the modules section](https://github.com/alice-lg/birdwatcher/blob/master/etc/birdwatcher/birdwatcher.conf)
of the config.

## Installation

You will need to have go installed to build the package.
Please make sure your go version is `>= 1.9`.

Running `go install github.com/alice-lg/birdwatcher@latest` will give you
a binary. You might need to cross-compile it for your
bird-running service (`GOARCH` and `GOOS` are your friends).

We provide a Makefile for more advanced compilation/configuration.
Running `make linux` will create a Linux executable (by default for
`amd64`, but that is configurable by providing the `ARCH` argument
to the Makefile).


#### 2.0 Breaking Change

The BIRD configuration setup (single/multi table, pipe/table prefixes) is no longer
configured in birdwatcher but directly in Alice-LG. Please have a look at the
[source section of the Alice-LG config example](https://github.com/alice-lg/alice-lg/blob/master/etc/alice-lg/alice.example.conf).


### BIRD configuration

Birdwatcher parses the output of `birdc[6]` and expects (for now)
the time format to be `iso long`. You need to configure

    timeformat base         iso long;
    timeformat log          iso long;
    timeformat protocol     iso long;
    timeformat route        iso long;

in your `/etc/bird[6].conf` for birdwatcher to work.

#### BIRD keep filtered routes
To also see filtered routes in configured BGP protocol instances, you need to make
sure that you have enabled the `import keep filtered on` option for affected bgp protocols.

    protocol bgp 'peerX' {
        ...
        import keep filtered on;
        ...
    }

Now you should be able to do a `show route filtered protocol peerX` in BIRD.

If you use a multi table setup you are also using the pipe protocol the connect the tables.
No special BIRD configuration is required to be able to query pipe filtered routes.

birdwatcher provides [various endpoints (see "available modules" section)](https://github.com/alice-lg/birdwatcher/blob/master/etc/birdwatcher/birdwatcher.conf)
to query routes filtered in bgp protocol as well as pipe protocol instances.

For use with [Alice-LG](https://github.com/alice-lg/alice-lg), make sure to set the appropriate BIRD config setup
in your [Alice-LG configuration](https://github.com/alice-lg/alice-lg/blob/master/etc/alice-lg/alice.example.conf).

#### BIRD tagging filtered routes
If you want to make use of the filtered route reasons in [Alice-LG](https://github.com/alice-lg/alice-lg), you need
to make sure that you are using BIRD 1.6.3 or up as you will need Large BGP Communities
(http://largebgpcommunities.net/) support.

You need to add a Large BGP Community just before you filter a route, for example:

    define yourASN = 12345
    define yourFilteredNumber = 65666
    define prefixTooLong = 1
    define pathTooLong = 2

    function importScrub() {
        ...
        if (net.len > 24) then {
            print "REJECTING: ",net.ip,"/",net.len," received from ",from,": Prefix is longer than 24: ",net.len,"!";
            bgp_large_community.add((YourASN,yourFilteredNumber,prefixTooLong));
            return false;
        }
        if (bgp_path.len > 64) then {
            print "REJECTING: ",net.ip,"/",net.len," received from ",from,": AS path length is ridiculously long: ",bgp_path.len,"!";
            bgp_large_community.add((yourASN,yourFilteredNumber,pathTooLong));
            return false;
        }
        ...
        return true;
    }

    function importFilter() {
        ...
        if !(importScrub()) then reject;
        ...
        accept;
    }

### Using Docker

You can run the birdwatcher for BIRD2 with docker:

    docker pull alicelg/birdwatcher:latest

    docker run -p 29184:29184 -v /var/run/bird.ctl:/usr/local/var/run/bird.ctl -it --rm birdwatcher:latest


Or build your own image:

    docker build . -t alicelg/birdwatcher:latest
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
[etc/birdwatcher/birdwatcher.conf](https://github.com/alice-lg/birdwatcher/blob/master/etc/birdwatcher/birdwatcher.conf).
You should be able to use it out of the box. If you need
to change it, it is well-commented and hopefully intuitive.
If you do not know how to configure it, please consider opening
[an issue](https://github.com/alice-lg/birdwatcher/issues/new).

## How

In the background `birdwatcher` runs the `birdc[6]` client, sends
commands and parses the result. It also leverages simple caching
techniques to help reduce the load on the BIRD service.

## Who

Initially developed by Daniel and MC from [Netnod](https://www.netnod.se/) in
two days at the RIPE 73 IXP Tools Hackathon in Madrid, Spain.

Running BIRD and parsing the results was added by [Veit Heller](https://github.com/hellerve/) on behalf of [ecix](http://ecix.net/).

With major contributions from: Patrick Seeburger and Benedikt Rudolph on behalf of [DE-CIX](https://de-cix.net).
